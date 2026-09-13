---
title: 'Ma vie en texte brut, quatre ans après : Planova et un fichier Markdown par jour'
date: '2026-09-13'
lang: fr
description: 'Quatre ans après ma première organisation en texte brut, j''ai remplacé todo.txt et mes formats maison par un seul fichier Markdown par jour : Planova sur Android et Linux, une CLI en Go, vim et Syncthing.'
ogDescription: 'Remplacer todo.txt et les formats maison par un fichier Markdown par jour : Planova sur Android et Linux, une CLI en Go, vim et Syncthing.'
image: https://rvier.fr/images/planova_screenshot_main.png
keywords: texte brut, plain text, Markdown, PIM, todo, agenda, journal, notes, Planova, Flutter, Go, Syncthing
summary: 'Version française du billet sur Planova : pourquoi j''ai abandonné todo.txt et mes formats maison pour un fichier Markdown par jour, le refill du matin, l''application Flutter, la CLI en Go et la synchronisation par Syncthing.'
---

En 2022 j'écrivais [My life in plain text](my-life-in-plain-text-EN.html) : les tâches dans un `todo.txt`, les événements dans un format maison `agendatxt`, les dépenses dans un autre format maison, les notes en Markdown, et une application Android de démonstration, MOrg, pour tenir tout ça ensemble sur mobile. Presque quatre ans plus tard, tout est toujours en texte brut, mais presque tout le reste a changé. Le changement principal tient en une phrase : **j'ai arrêté de découper ma vie par type de données, et je la découpe par jour.**

## Ce qui n'allait pas dans la première version

La structure de 2022 avait l'air propre : un fichier par sujet (`todo.txt`, `agenda.txt`, `expenses.txt`, un dossier journal). À l'usage, trois problèmes revenaient sans arrêt :

- **Un format maison demande un parseur maison.** `agendatxt` et `expensetxt` étaient triviaux à définir, mais chaque outil qui y touchait, l'application Android, les alias shell, les scripts, devait réimplémenter le même parsing et les mêmes cas limites. Un format que je suis seul à utiliser n'est pérenne qu'en théorie.
- **Une journée était éparpillée dans plusieurs fichiers.** « Qu'est-ce qui s'est passé mardi ? » voulait dire grepper une date dans l'agenda, dans la liste des tâches faites et dans le journal, puis recoller les morceaux de tête.
- **Un `todo.txt` global ne fait que grossir.** Sans date attachée à quoi que ce soit, le fichier est devenu une liste de culpabilité : des centaines de lignes, aucune idée de quand une tâche avait été ajoutée ou abandonnée en silence.

MOrg, l'application mobile, n'a franchement jamais dépassé le stade de la démonstration. Elle marchait, je m'en servais tous les jours, mais chaque format maison la rendait plus lourde à maintenir.

## Un fichier Markdown par jour

Le remplaçant s'appelle [Planova](https://planova.rvier.fr/), et son idée centrale est le fichier quotidien : `dailies/20260712.md`, du Markdown ordinaire, quatre sections :

```
## Events

- @09:30 Standup
- @14:00 Review Planova release

## Tasks

- [ ] Write the blog post about Planova
- [x] Fix the widget refresh on Android

## Journal

Spent the morning chasing a cache invalidation bug...

## Notes

Links and thoughts that belong to today.
```

Il n'y a aucun format à apprendre au-delà de deux conventions que les utilisateurs de Markdown connaissent déjà : les cases à cocher pour les tâches, et un préfixe `@HH:MM` pour transformer un élément de liste en événement. Tout ce qui concerne une journée, ce qui était prévu, ce qui a été fait, ce que j'en ai pensé, vit dans un seul fichier que n'importe quel éditeur, `grep` ou `git diff` sait lire. Le reste de l'arborescence est tout aussi banal, et c'est voulu :

```
Org/
├── dailies/
│   ├── 20260711.md
│   └── 20260712.md
├── archives/
└── notes/
    ├── Work/
    │   └── project_1.md
    └── Personal/
        └── project_2.md
```

## Le refill, à la place de la liste de tâches globale

Si les tâches vivent dans des fichiers quotidiens, que deviennent celles qu'on n'a pas finies ? C'était mon principal doute sur le modèle par jour, et la réponse est devenue ma fonctionnalité préférée. Chaque matin, Planova propose un *refill* : les tâches non faites de tous les jours précédents, les plus récentes d'abord, et on choisit celles qu'on reporte à aujourd'hui.

C'est un petit rituel, trente secondes, mais il inverse le mode de défaillance de `todo.txt`. Une tâche doit maintenant *mériter* sa place dans le fichier du jour, chaque jour. Ce que je refuse de reporter plusieurs fois de suite, c'est ce que j'ai visiblement décidé de ne pas faire, et les vieux fichiers quotidiens gardent une trace honnête du moment où j'ai laissé tomber.

<img src="../images/planova_screenshot_main.png" alt="Écran principal de Planova : calendrier du mois au-dessus du fichier Markdown du jour" loading="lazy" width="360">

## Sur mobile : l'application Planova

Planova est une application Flutter, qui tourne sur Android et sur Linux. Comme MOrg elle est minimaliste et pleine de partis pris, mais cette fois elle a dépassé le stade de la démonstration : la version 1.3 est sortie et c'est mon outil quotidien depuis plus d'un an :

- un calendrier mensuel avec le Markdown du jour sélectionné en dessous, modifiable sur place, avec un double appui pour cocher ou décocher une tâche ;
- des notifications générées directement depuis les lignes `@HH:MM`, sans base de données de calendrier à côté ;
- un widget d'écran d'accueil avec les prochains événements et les tâches du jour ;
- l'import ICS, pour qu'une invitation partagée depuis un client mail atterrisse dans le bon fichier quotidien ;
- une vue notes sur le dossier `notes/`, avec recherche, sous-dossiers et renommage.

## Sur le bureau : vim, et une CLI en Go

Sous Linux je vis toujours dans un terminal, donc les fichiers quotidiens sont à un `vim` de distance. Pour les opérations structurées, j'ai écrit une petite CLI compagnon en Go :

```
planova calendar                      # vue du mois, avec marqueurs événements/tâches
planova todo --add "Buy groceries"    # ajoute à la section ## Tasks du jour
planova event --add "@14:30 Meeting"  # ajoute à la section ## Events du jour
planova daily --edit                  # ouvre le fichier du jour dans $EDITOR
```

Comme l'application et la CLI ne partagent rien d'autre que les fichiers eux-mêmes, elles ne peuvent pas être en désaccord. Le fichier *est* l'état. N'importe quel script capable d'ajouter une ligne à un fichier Markdown est un client Planova valide.

## Toujours pas de synchronisation intégrée, exprès

Planova n'a ni compte, ni serveur, ni synchronisation à lui, et c'est une fonctionnalité. Mon dossier `Org/` est synchronisé entre le téléphone, le portable et la machine de bureau avec Syncthing ; Nextcloud, Dropbox ou un dépôt git feraient tout aussi bien l'affaire. Le travail de l'application est de lire et d'écrire des fichiers, pas de les déplacer.

## Conclusion

Quatre ans plus tard, le pari du texte brut a bien vieilli : les fichiers de 2022 sont toujours lisibles, ce que peu d'applications de l'époque peuvent dire. Le pari des formats maison, lui, non : moins j'ai inventé de conventions, plus le système a duré. Un fichier Markdown par jour, un dossier de notes, et des outils qui traitent les fichiers comme la seule source de vérité, c'est finalement toute la structure dont ma vie a besoin.

Planova est sous licence MIT, comme tout le reste de cette organisation. Si l'approche vous parle :

- [planova.rvier.fr](https://planova.rvier.fr/) : le site de l'application ;
- [Planova sur GitHub](https://github.com/brvier/Planova) : l'application Flutter (Android et Linux) ;
- la CLI en Go, pour le côté terminal des mêmes fichiers.
