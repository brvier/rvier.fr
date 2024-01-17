Title: J'ai hacké ma sonette !
Date: 2023-05-05
Lang: FR

En tant que développeur, la concentration et la focalisation sont cruciales pour accomplir les tâches de manière efficace et efficiente. Malheureusement, travailler dans un environnement bruyant ou distrayant peut représenter un défi important pour atteindre une concentration optimale. C'est là que les casques antibruit sont utiles. Ils bloquent les sons externes et réduisent les distractions sonores.
Mais je travaille également chez moi et je reçois souvent des livraisons, donc je n'entends pas la sonnette de la porte.

Du coup, j ai hacké ma sonette !

## Objectifs

- Garder la solution sans fil
- Être toujours averti par le son sur le recepteur
- Être informé par une notification sur le téléphone et l'ordinateur
- Résoudre le problème de décharge rapide de la batterie du tranmetter
- Réaliser un projet avec un Raspberry Pi Pico et MicroPython :)

## Solutions

- Pirater le récepteur RF :
  - Pirater le récepteur RF pourrait être la solution la plus facile, sauf que l'actuel émetteur RF utilise une pile CRC2032.
  Cet émetteur ne fonctionne pas bien si la pile est inférieure à 70% car la portée de transmission est trop faible et le récepteur ne la reçoit pas.

- Écoute de la RF 433 Mhz :
  - Une autre solution pourrait consister à écouter la RF 433 Mhz, mais là encore, je dois modifier la puissance de l'émetteur RF.

- Pirater l'émetteur RF :
  - Le bouton est pull down to gnd lorsqu'il est enfoncé.
  - L'émetteur fonctionne bien avec une alimentation de 3,3 V.

  J'ai donc décidé de partir sur cette solution, ayant a ma disposition un mini paneau solaire, son chargeur, et un pico.


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
         │   │G
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

On va pas se mentir, monitorer un port GPIO, ce n'est pas la partie difficile. Le soucis est la consomation du RPI Pico. 21mA sans "deepsleep".
La methode deepsleep du port de micropython ne permettant pas d'etre reveillé par un Interupt GPIO, j ai trouvé cette lib qui le fait pour nous : https://github.com/tomjorquera/pico-micropython-lowpower-workaround .

Resultats :

- 21mA sans le mode "dormant"
- 0.06mA avec

Pour les notification j utilise ma propre instance auto hebergé de gotify.

Source : https://github.com/brvier/sonette

## Picture

![Sonette](images/sonette.png)

