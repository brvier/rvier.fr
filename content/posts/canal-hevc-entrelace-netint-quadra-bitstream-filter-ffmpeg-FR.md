---
title: 'Les flux des diffuseurs sont parfois vraiment très crades'
date: '2026-09-12'
lang: fr
description: 'Une chaîne Canal+ reçue par satellite en HEVC 1080i faisait tourner nos encodeurs en décodage logiciel, ni NVDEC ni NetInt Quadra n''acceptant le HEVC entrelacé. En regardant le flux de près, les images sont progressives : seuls trois drapeaux disent le contraire. Un bitstream filter ffmpeg de 450 lignes les efface, la carte décode, et sauve quatre fois plus d''images que le CPU sur ce flux corrompu à 12 %.'
ogDescription: 'Du HEVC 1080i que le NetInt Quadra refusait : les images sont progressives, seuls trois drapeaux mentent. Un bitstream filter ffmpeg les efface.'
keywords: ffmpeg, HEVC, Canal+, bitstream filter, NetInt, Quadra, entrelacé, field-coded, 1080i, MPEG-TS, DVB, satellite, broadcast, captation, Go
summary: 'Comment un flux Canal+ en HEVC 1080i field-coded, refusé par le décodeur NetInt Quadra, est passé en décodage matériel grâce à un bitstream filter ffmpeg qui efface la signalisation d''entrelacement. Avec le code, la ligne ffmpeg, les résultats sur un flux corrompu, et la réponse de NetInt.'
---

Dans la [plateforme de captation](recording-broadcast-24-7-capture-platform-architecture-EN.html), il y a un flux que j'ai fini par appeler « le flux merdique de Canal ». C'est une chaîne du bouquet Canal+ reçue par satellite, « Barker à la une », en HEVC 1080i :

- un entrelacement qu'aucun de nos décodeurs matériels n'accepte
- 12 % de paquets corrompus
- des timestamps faux
- un `ffprobe` qui n'arrive pas à la sonder. 

Les flux des diffuseurs sont parfois vraiment très crades, et celui-là est mon pire exemple.

## Pas de décodage Hardware

Nos encodeurs tournent sur GPU NVIDIA ou des TPUs NetInt. NVDEC ne sait pas décoder du HEVC entrelacé, et NetInt non plus, donc ces chaînes étaient décodées en logiciel, sur le CPU. Et ca consomme, au point d'avoir des pertes de packets multicasts quand de nombreux flux HEVC entrelacés sont traités.

Fallait donc essayer d'optimiser un peu cela.

## Ce flux est entrelacé sans vraiment l'etre

Voilà ce que `ffprobe` en dit :

```
Stream #0:0: Video: hevc (Main), yuv420p(tv, bt709, top first), 1440x540, 50 fps
```

540 lignes de haut, 50 images par seconde, la ou l'on devrait s attendre a du 1080 25fps. Le flux envoie chaque demi-image (chaque trame, ou *field*) comme une image à part entière. On appelle ça du *field-coded*. Pour le codec, ce sont des images progressives ordinaires, juste moitié moins hautes. C'est au décodeur de recoller les paires ensuite.

Alors d'où vient le « entrelacé » qui fait tout planter ? De trois drapeaux posés à côté des images, qui disent au décodeur « attention, c'est de l'entrelacé » :

- deux bits dans l'en-tête général du flux (le `profile_tier_level`, présent dans le VPS et le SPS) ;
- deux bits dans les infos d'affichage du SPS (le VUI) ;
- un petit message optionnel attaché à chaque image (le SEI `pic_timing`), qui dit si c'est la trame du haut ou du bas.

Ces drapeaux ne changent rien aux images. Ils disent juste comment les afficher. Le Quadra sait donc décoder ces images. Il refuse uniquement parce qu'on lui a dit qu'elles étaient entrelacées.

D'où l'idée, effacer les drapeaux. Les images restent identiques, le décodeur n'a plus de raison de refuser. Et pour recoller les trames deux par deux après, ffmpeg a déjà un filtre qui fait ça, `weave`, malheureusement sur CPU.

Il nous restera donc : 

```NetInt Quadra decode→ SW weave,pp=lb → hwupload → Quadra scale + encode```

au lieu de :

```CPU decode  → SW weave,pp=lb → hwupload → Quadra scale + encode```

## Un bitstream filter

Le plus simple pour modifier le flux avant le decodeur, c est d utiliser un bitstream filter (BSF). Il voit passer les paquets, les modifie, les rend. On l'active avec `-bsf:v`. Ça ne coûte presque rien, et comme il n'est posé que sur le process d'encodage, les fichiers enregistrés sur disque restent tels qu'ils sont arrivés du satellite.

Le filtre fait trois choses.

**1. Les deux bits de l'en-tête général.** Ils sont toujours au même endroit depuis le début du bloc. Pas besoin de lire quoi que ce soit, on écrit deux bits en aveugle :

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

**2. Les deux bits du VUI.** Il faut lire tout ce qui précède, sans en faire rien, juste pour savoir où on est arrivé. Le cœur de la fonction (coupé, la version complète fait 120 lignes) :

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

Si le SPS est abîmé et que la lecture échoue, on tombe dans `fail`. Les deux bits de l'en-tête général ont déjà été écrits, le reste est laissé tel quel, et le flux continue.

**3. Les messages pic_timing.** En gros, « trame du haut » ou « trame du bas » pour chaque image. Ces messages voyagent dans un bloc SEI, qui est une simple liste de messages `(type, taille, contenu)`. On recopie tous les messages sauf ceux de type 1, et si le bloc devient vide on le supprime :

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

Un octet `03` est inséré dans les données d'un flux HEVC à chaque fois que la séquence `00 00` est suivie d'un octet de 0 à 3, pour que ça ne ressemble pas à un début de bloc. Il nous faut donc retirer ces `03` avant de compter les bits, faire les modifications, puis les remettre. Sinon les positions sont décalées dès qu'un `03` traîne avant le VUI. 

Le reste du fichier découpe le paquet en blocs, ne s'occupe que des trois types qui nous intéressent (VPS, SPS, SEI), recopie tout le reste tel quel, et applique le même traitement à l'en-tête initial que le démuxeur fournit à part. Voilà la déclaration du filtre, avec une garde pour compiler aussi bien sur le ffmpeg 4.4 de nos vieux encodeurs que sur le 7.1 de l'arbre NetInt :

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

## L'utiliser

Le filtre se pose sur l'entrée. Ensuite `weave` recolle les trames deux par deux, et on désentrelace :

```
ffmpeg -bsf:v hevc_force_progressive -c:v h265_ni_quadra_dec -i buffer.m3u8 \
  -filter_complex '[0:v]weave,pp=lb,ni_quadra_hwupload,ni_quadra_scale=1920:1080[v]' \
  -map '[v]' -c:v h264_ni_quadra_enc ...
```

`weave` doit savoir si la première trame est celle du haut ou du bas. Cette info était dans les messages qu'on vient de supprimer. Il faut donc la lire avant le filtre. Le flux Canal+ commence par la trame du haut, ce qui est le défaut de `weave`, donc ici ça passe sans option.

Le désentrelacement reste sur le CPU, avec `pp=lb`, parce que le désentrelaceur de la carte n'accepte pas non plus ce format. On archive et on analyse ces images, on ne les rediffuse pas. Le mélange linéaire est trois fois moins cher que `yadif` et la différence ne se voit pas pour cet usage.

## Bilans

Avec 12,6 % de paquets marqués en erreur par le tuner. Le CPU en tire 91 images. Le Quadra en tire 377. Le décodeur matériel encaisse quatre fois mieux la corruption que libavcodec. Je n'aurais pas parié là-dessus, et au final, moins de charge CPU sur nos encodeurs, environ 70% de gain.

J'ai envoyé un ticket à NetInt avec trois échantillons : le flux réel, un flux synthétique propre, et le réel passé par le filtre. Puisque le firmware décode parfaitement ces images une fois les drapeaux effacés, aucune raison de conserver cette limitation.

Je dois dire que NetInt a été très collaboratif. Le ticket a été lu par quelqu'un qui comprenait le problème, les échantillons ont servi, et ils ont proposé d'implémenter ça de leur côté sous forme d'un simple pp. Ce n'est pas la réponse qu'on obtient de tous les fabricants de matériel. Et en attendant, mon patch tourne en production depuis fin août.



