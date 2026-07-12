---
title: 3 Erreurs Fréquentes d'utilisation du Module Logging de Python
date: '2024-01-19'
lang: fr
description: 'Les 3 erreurs les plus fréquentes avec le module logging de Python : mauvaise instanciation du Logger, exceptions masquées avec error() au lieu de exception(), et mauvais niveaux de log.'
ogDescription: 'Mauvaise instanciation du Logger, exceptions masquées, mauvais niveaux de log : les 3 erreurs les plus fréquentes avec le module logging de Python.'
keywords: Python, logging, Dev
summary: Mauvaise instanciation du Logger, exceptions masquées avec error() au lieu de exception(), et mauvais niveaux de log.
---

Le module python **logging** est un excellent module, mais je le vois souvent très mal utilisé. Voici les 3 erreurs les plus fréquentes, et les plus grosses que j'ai croisées.

## Instanciation

Les objets Logger ne doivent JAMAIS être instanciés directement, mais toujours à travers la fonction du module `logging.getLogger(__name__)`.

Il est même recommandé d'utiliser `__name__` comme nom dans la documentation officielle.

**Pas bien :**

```
import logging
logger = logging.Logger(__name__)
```

**Bien :**

```
import logging
logger = logging.getLogger(__name__)
```

## Cacher les informations des exceptions

Quoi de plus pénible à débugger qu'une erreur sans les informations de l'exception ? Au lieu d'utiliser la méthode **error**, utilisez la méthode **exception** : celle-ci a le même niveau d'erreur (ERROR), mais va passer en plus la trace de l'exception.

**Pas bien :**

```
import logging
logger = logging.getLogger(__name__)

try:
    truc()
except TrucError as err:
    logger.error("Truc failed : %s", err)
```

**Bien :**

```
import logging
logger = logging.getLogger(__name__)

try:
    truc()
except TrucError:
    logger.exception("Truc failed")
```

## Ne pas utiliser le bon niveau d'erreur

J'ai vu des centaines de gigaoctets de logs inutiles dans de nombreux projets, car presque tous les appels à logging étaient au niveau WARNING, voire pire, ERROR. Voire même des logs d'exécution de chaque requête SQLAlchemy en INFO, cela n'a rien à faire en dehors de DEBUG.

Voici une description des niveaux plutôt claire de sametmax.fr :

```
CRITICAL    50  Le programme complet est en train de partir en couille.
ERROR       40  Une opération a foirée.
WARNING     30  Pour avertir que quelque chose mérite l'attention : enclenchement d'un mode particulier, detection d'une situation rare, un lib optionelle peut être installée.
INFO        20  Pour informer de la marche du programme. Par exemple : "Starting CSV parsing"
DEBUG       10  Pour dumper des informations quand vous débuggez. Par exemple savoir ce qu'il y a dans ce putain de dictionnaire.
```
