Title: 3 Erreurs Fréquentes d'utilisation du Module Logging de Python
Date: 2024-01-19
Tags: Dev, Python
Lang: FR

Le module python **logging** est un excellent module, mais je le vois souvent
très mal utilisé. Voici les 3 erreurs les plus fréquentes,
et les plus grosses que j'ai croisées.

## Instanciation

Les objets Logger ne doivent JAMAIS être instanciés directement,
mais toujours à travers la fonction du module logging.getLogger(__name__).

Il est même recommandé d'utiliser \_\_name\_\_ comme nom dans la documentation officielle

**Pas bien:**

    #!python
    import logging
    logger = logging.Logger(__name__)

**Bien :**

    #!python
    import logging
    logger = logging.getLogger(__name__)

## Cacher les informations des exceptions

Quoi de plus pénible à débugger qu'une erreur sans les informations de l'exception.
Au lieu d'utiliser la méthode **error**, utilisez la méthode **exception**,
celle-ci a le même niveau d'erreur (ERROR). Mais va passer en plus la trace de
l'exception.

**Pas bien:**

    #!python
    import logging
    logger = logging.Logger(__name__)

    try:
        truc()
    except TrucError as err:
        logger.error("Truc failed : %s", err)

**Bien :**

    #!python
    import logging
    logger = logging.Logger(__name__)

    try:
        truc()
    except TrucError:
        logger.exception("Truc failed")

## Ne pas utiliser le bon niveau d'erreur

J'ai vu des centaines de gigaoctets de logs inutiles dans de nombreux projets,
car le niveau de presque tous les appels à logging était en niveau WARNING,
voir pire ERROR. Voir même des logs d'exécutions de chaque requêtes SQLALCHEMY en INFO,
cela n'a rien a faire en dehors de DEBUG.


Voici une description des niveaux plutôt claire de sametmax.fr :

    CRITICAL    50 	Le programme complet est en train de partir en couille.
    ERROR       40 	Une opération a foirée.
    WARNING     30 	Pour avertir que quelque chose mérite l’attention : enclenchement d’un mode particulier, detection d’une situation rare, un lib optionelle peut être installée.
    INFO        20 	Pour informer de la marche du programme. Par exemple : “Starting CSV parsing”
    DEBUG       10 	Pour dumper des informations quand vous débuggez. Par exemple savoir ce qu’il y a dans ce putain de dictionnaire.



