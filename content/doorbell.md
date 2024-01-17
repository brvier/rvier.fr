Title: I hacked my doorbell !
Date: 2023-05-05
Lang: EN

As a developer, concentration and focus are crucial to accomplish tasks effectively and efficiently. Unfortunately, working in a noisy or distracting environment can present a significant challenge to achieving optimal concentration. This is where noise-canceling headphones come in handy. They block out external sounds and reduce noise distractions.
But I also work from home and often get deliveries, so I don't hear the doorbell.

So I hacked my doorbell!

## Objectives

- Keep the wireless solution
- Always be notified by the sound on the receiver
- Be informed by a notification on the phone and the computer
- Solve the problem of fast discharge of the battery of the transceiver
- Realize a project with a Raspberry Pi Pico and MicroPython :)

## Solutions

- Hack the RF receiver:
  - Hacking the RF receiver might be the easiest solution, except that the current RF transmitter uses a CRC2032 battery.
  This transmitter does not work well if the battery is below 70% because the transmission range is too low and the receiver does not receive it.

- Listening to 433 Mhz RF:
  - Another solution could be to listen to the 433 Mhz RF, but again, I have to change the power of the RF transmitter.

- Hacking the RF transmitter:
  - The button is pull down to gnd when pressed.
  - The transmitter works well with a 3.3 V power supply.

  So I decided to go with this solution, having at my disposal a mini solar panel, its charger, and a pico.

````
 ┌────────────────────────────────────────────────────────────────────────────────┐
 │                                                                                │
 │   ┌────────────────────────────────────────────────────────────────────────┐   │
 │   │                                                                        │   │
 │   │        Button                                         CRC2032          │   │
 │   │     ┌──────────┐                                   ┌───────────┐       │   │
 │   │     │          │                                   │     -     ├───────┼───┼─────────────────┐
 │   │     │          │                                   │           │       │   │                 │
 │   │     │          │                                   │           │       │   │                 │
 │   │     │          │                                   │     +     ├───────┼───┼─────────────┐   │
 │   │     └────┬─────┘                                   └───────────┘       │   │             │   │
 │   │          │                                                             │   │             │   │
 │   │          │                                                             │   │             │   │
 │   └──────────┼─────────────────────────────────────────────────────────────┘   │             │   │
 │              │                                                                 │             │   │
 └──────────────┼─────────────────────────────────────────────────────────────────┘             │   │
                │                                                                               │   │
                │                                                                               │   │
                │                                                                               │   │
                │       220 Ohms                                                                │   │
                └───────/\/\/\/───────────────────────────────────────────────────────────┐     │   │
                                                                                          │     │   │
                                                                                          │     │   │
                                                                                          │     │   │
                                                                                          │     │   │
                                                                                          │     │   │
                                                                                          │     │   │
                                                                                          │     │   │
                                                                                          │     │   │
                                                                                          │     │   │
                                                                                          │     │   │
                                                                                          │     │   │
                                                                                          │     │   │
                                                                                          │     │   │
   ┌────────────────────────────────────────────────────────────────────────────────┐     │     │   │
   │                                                                                │     │     │   │
   │                                                                                │     │     │   │
   │                                                                                │     │     │   │
   │                                                                                │     │     │   │
   │                                                                                │     │     │   │
   │                                     Raspberry Pico W                           │     │     │   │
   │                                                                                │     │     │   │
   │                                                                                │     │     │   │
   │                                                                                │     │     │   │
   │                                                                                │     │     │   │
   │   VSYS GND      3.3V                                              PIN 17       │     │     │   │
   │                                                                                │     │     │   │
   └─────┬───┬────────┬───────────────────────────────────────────────────┬─────────┘     │     │   │
         │   │        │                                                   │               │     │   │
         │   │        │                                                   └───────────────┘     │   │
         │   │        │                                                                         │   │
         │   │        │                                                                         │   │
         │   │        │                                                                         │   │
         │   │        │                                                                         │   │
         │   │        │                                                                         │   │
         │   │        └─────────────────────────────────────────────────────────────────────────┘   │
         │   │                                                                                      │
         │   │                                                                                      │
         │   ├──────────────────────────────────────────────────────────────────────────────────────┘
         │   │
         │   │
         │   │
         │   │
         │   │
         │   │
         │   │
         │   │
         │   │
         │   │
         │   │
         │   │
         │   │
┌────────┴───┴───────────────────┐               ┌──────────────────────────────┐
│         Solaar Charger         │               │         Solaar Panel         │
│                                ├───────────────┤                              │
│                                │               │                              │
│                                │               │                              │
│                                ├───────────────┤                              │
│                                │               │                              │
│                                │               │                              │
│                                │               │                              │
│                                │               │                              │
│                                │               │                              │
│                                │               │                              │
└────────────────────────┬────┬──┘               └──────────────────────────────┘
                         │    │
                         │    │                  ┌──────────────────────────────┐
                         │    │                  │        Battery 18650         │
                         │    │                  │                              │
                         │    └──────────────────┤                              │
                         │                       │                              │
                         │                       │                              │
                         │                       │                              │
                         │                       │                              │
                         └───────────────────────┤                              │
                                                 │                              │
                                                 │                              │
                                                 │                              │
                                                 └──────────────────────────────┘

````

## Software
Let's not lie, monitoring a GPIO port is not the hard part. The problem is the consumption of the Pico RPI. 21mA without deepsleep.
The deepsleep method of the micropython port does not allow to be woken up by a GPIO Interupt, I found this lib that does it for us : https://github.com/tomjorquera/pico-micropython-lowpower-workaround .

Results :

- 21mA without the "sleep" mode
- 0.06mA with

For notification I use my own self-hosted instance of gotify.

Source : https://github.com/brvier/ihackedmydoorbell

## Picture

![Sonette](images/sonette.png)

