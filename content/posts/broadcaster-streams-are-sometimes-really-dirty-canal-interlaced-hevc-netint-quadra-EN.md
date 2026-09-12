---
title: 'Broadcaster streams are sometimes really dirty'
date: '2026-09-12'
lang: en
description: 'A Canal+ satellite channel in HEVC 1080i kept our encoders in software decoding, since neither NVDEC nor the NetInt Quadra accept interlaced HEVC. Looking closely at the stream, the pictures are progressive: only three flags say otherwise. A 450-line ffmpeg bitstream filter erases them, the card decodes, and salvages four times more pictures than the CPU on this 12% corrupted feed.'
ogDescription: 'HEVC 1080i the NetInt Quadra refused: the pictures are progressive, only three flags lie. An ffmpeg bitstream filter erases them.'
keywords: ffmpeg, HEVC, Canal+, bitstream filter, NetInt, Quadra, interlaced, field-coded, 1080i, MPEG-TS, DVB, satellite, broadcast, capture, Go
summary: 'How a field-coded HEVC 1080i Canal+ feed, refused by the NetInt Quadra decoder, made it to hardware decoding thanks to an ffmpeg bitstream filter that erases the interlace signalling. With the code, the ffmpeg command line, the results on a corrupted feed, and NetInt''s response.'
---

In the [capture platform](recording-broadcast-24-7-capture-platform-architecture-EN.html), there is one feed I ended up calling "the crappy Canal feed". It is a channel from the Canal+ satellite bouquet, "Barker à la une", in HEVC 1080i:

- an interlacing that none of our hardware decoders accept
- 12% corrupted packets
- wrong timestamps
- an `ffprobe` that cannot even probe it.

Broadcaster streams are sometimes really dirty, and this one is my worst example.

## No hardware decoding

Our encoders run on NVIDIA GPUs or NetInt TPUs. NVDEC cannot decode interlaced HEVC, and neither can NetInt, so these channels were decoded in software, on the CPU. And that eats CPU, to the point of losing multicast packets when many interlaced HEVC feeds are being processed.

So this had to be optimised a bit.

## This feed is interlaced without really being interlaced

Here is what `ffprobe` says about it:

```
Stream #0:0: Video: hevc (Main), yuv420p(tv, bt709, top first), 1440x540, 50 fps
```

540 lines high, 50 pictures per second, where you would expect 1080 lines at 25 fps. The feed sends each half-picture (each field) as a full picture of its own. This is called *field-coded*. For the codec, these are ordinary progressive pictures, just half as tall. Pairing them back up is the decoder's job afterwards.

So where does the "interlaced" that breaks everything come from? From three flags sitting next to the pictures, telling the decoder "careful, this is interlaced":

- two bits in the general stream header (the `profile_tier_level`, present in both the VPS and the SPS);
- two bits in the display information of the SPS (the VUI);
- a small optional message attached to each picture (the `pic_timing` SEI), saying whether it is the top or the bottom field.

These flags change nothing about the pictures. They only say how to display them. So the Quadra can decode these pictures. It refuses only because it has been told they are interlaced.

Hence the idea: erase the flags. The pictures stay identical, the decoder has no reason left to refuse. And to pair the fields back up afterwards, ffmpeg already has a filter for that, `weave`, unfortunately on the CPU.

So we end up with:

```NetInt Quadra decode → SW weave,pp=lb → hwupload → Quadra scale + encode```

instead of:

```CPU decode → SW weave,pp=lb → hwupload → Quadra scale + encode```

## A bitstream filter

The simplest way to modify the stream before the decoder is a bitstream filter (BSF). It sees the packets go by, modifies them, hands them back. You enable it with `-bsf:v`. It costs almost nothing, and since it only sits on the encoding process, the files recorded on disk stay exactly as they came off the satellite.

The filter does three things.

**1. The two bits in the general header.** They are always at the same place from the start of the block. No need to read anything, we write two bits blind:

```c
static void set_bit(uint8_t *buf, int pos, int val)
{
    if (val)
        buf[pos >> 3] |=   0x80 >> (pos & 7);
    else
        buf[pos >> 3] &= ~(0x80 >> (pos & 7));
}

static void patch_vps(uint8_t *buf, int size)
{
    /* PTL of the VPS starts at 16 (NAL header) + 32 fixed bits */
    if (size * 8 >= 48 + 42) {
        set_bit(buf, 48 + 40, 1);  /* general_progressive_source_flag */
        set_bit(buf, 48 + 41, 0);  /* general_interlaced_source_flag  */
    }
}
```

**2. The two bits in the VUI.** Everything before it has to be read, doing nothing with it, just to know where you have got to. The core of the function (cut down, the full version is 120 lines):

```c
    /* PTL of the SPS starts at a fixed offset: 16 (NAL header) + 8 bits */
    set_bit(buf, 24 + 40, 1);   /* general_progressive_source_flag */
    set_bit(buf, 24 + 41, 0);   /* general_interlaced_source_flag  */

    init_get_bits8(&gb, buf, size);
    skip_bits(&gb, 16);                    /* NAL header */
    skip_bits(&gb, 4);                     /* sps_video_parameter_set_id */
    max_sub = get_bits(&gb, 3);            /* sps_max_sub_layers_minus1 */
    skip_bits1(&gb);                       /* temporal_id_nesting */
    if (parse_ptl(&gb, max_sub, &ptl_pos) < 0)
        goto fail;

    get_ue_golomb_long(&gb);               /* sps_seq_parameter_set_id */
    chroma = get_ue_golomb_long(&gb);      /* chroma_format_idc */
    ...
    num_rps = get_ue_golomb_long(&gb);
    for (i = 0; i < num_rps; i++)
        if (parse_st_ref_pic_set(&gb, i, num_rps, num_delta_pocs) < 0)
            goto fail;
    ...
    if (!get_bits1(&gb))                   /* vui_parameters_present */
        return;

    /* VUI up to field_seq_flag */
    if (get_bits1(&gb)) {                  /* aspect_ratio_info_present */
        if (get_bits(&gb, 8) == 255)       /* aspect_ratio_idc */
            skip_bits_long(&gb, 32);
    }
    ...
    skip_bits1(&gb);                       /* neutral_chroma_indication */

    set_bit(buf, get_bits_count(&gb), 0);     /* field_seq_flag */
    set_bit(buf, get_bits_count(&gb) + 1, 0); /* frame_field_info_present */
    return;

fail:
    av_log(bsf, AV_LOG_WARNING,
           "SPS parse failed, only PTL flags rewritten\n");
```

If the SPS is damaged and the read fails, we land in `fail`. The two bits of the general header have already been written, the rest is left as is, and the stream carries on.

**3. The pic_timing messages.** Basically, "top field" or "bottom field" for each picture. These messages travel inside an SEI block, which is a plain list of `(type, size, content)` messages. We copy every message except those of type 1, and if the block ends up empty we drop it:

```c
static int strip_pic_timing_sei(uint8_t *buf, int size)
{
    int in = 2, out = 2, kept = 0;

    while (in < size) {
        int type = 0, len = 0, msg_start = in;

        if (buf[in] == 0x80 && in == size - 1)
            break;                          /* rbsp trailing bits */
        while (in < size && buf[in] == 0xFF)
            type += buf[in++] ? 255 : 0;
        type += buf[in++];
        while (in < size && buf[in] == 0xFF) {
            len += 255;
            in++;
        }
        len += buf[in++];
        if (in + len > size)
            return size;                    /* parse failure: keep NAL as is */
        in += len;

        if (type != SEI_TYPE_PIC_TIMING) {
            memmove(buf + out, buf + msg_start, in - msg_start);
            out += in - msg_start;
            kept++;
        }
    }
    if (!kept)
        return 0;                           /* drop the whole NAL */
    buf[out++] = 0x80;                      /* rbsp_stop_one_bit */
    return out;
}
```

A `03` byte is inserted into the data of an HEVC stream every time the sequence `00 00` is followed by a byte between 0 and 3, so that it does not look like the start of a block. So we have to remove these `03` bytes before counting bits, make the changes, then put them back. Otherwise every position is shifted as soon as a `03` sits somewhere before the VUI.

The rest of the file splits the packet into blocks, only deals with the three types we care about (VPS, SPS, SEI), copies everything else verbatim, and applies the same treatment to the initial header that the demuxer provides separately. Here is the filter declaration, with a guard so it compiles both on the ffmpeg 4.4 of our old encoders and on the 7.1 of the NetInt tree:

```c
#if LIBAVCODEC_VERSION_MAJOR >= 59
const FFBitStreamFilter ff_hevc_force_progressive_bsf = {
    .p.name         = "hevc_force_progressive",
    .p.codec_ids    = hfp_codec_ids,
    .priv_data_size = sizeof(HFPContext),
    .init           = hfp_init,
    .close          = hfp_close,
    .filter         = hfp_filter,
};
#else
const AVBitStreamFilter ff_hevc_force_progressive_bsf = {
    .name           = "hevc_force_progressive",
    ...
};
#endif
```

## Using it

The filter goes on the input. Then `weave` pairs the fields back up, and we deinterlace:

```
ffmpeg -bsf:v hevc_force_progressive -c:v h265_ni_quadra_dec -i buffer.m3u8 \
  -filter_complex '[0:v]weave,pp=lb,ni_quadra_hwupload,ni_quadra_scale=1920:1080[v]' \
  -map '[v]' -c:v h264_ni_quadra_enc ...
```

`weave` needs to know whether the first field is the top or the bottom one. That information was in the messages we just removed. So it has to be read before the filter. The Canal+ feed starts with the top field, which is the default for `weave`, so here it works without any option.

Deinterlacing stays on the CPU, with `pp=lb`, because the card's deinterlacer does not accept this format either. We archive and analyse these pictures, we do not rebroadcast them. Linear blend is three times cheaper than `yadif` and the difference is invisible for this use.

## Outcome

With 12.6% of packets flagged as errored by the tuner, the CPU salvages 91 pictures. The Quadra salvages 377. The hardware decoder copes with corruption four times better than libavcodec. I would not have bet on that, and in the end, less CPU load on our encoders, about a 70% gain.

I sent NetInt a ticket with three samples: the real feed, a clean synthetic one, and the real one passed through the filter. Since the firmware decodes these pictures perfectly once the flags are erased, there is no reason to keep this limitation.

I have to say NetInt has been very cooperative. The ticket was read by someone who understood the problem, the samples were put to use, and they offered to implement this on their side as a simple pp. That is not the answer you get from every hardware vendor. And in the meantime, my patch has been running in production since late August.
