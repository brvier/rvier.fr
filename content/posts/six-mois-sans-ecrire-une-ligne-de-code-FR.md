---
title: 'Six mois sans écrire une ligne de code'
date: '2026-11-23'
lang: fr
description: 'Depuis six mois, je n''écris plus de code moi-même, sauf pour le plaisir d''un casse-tête ou d''une session de débogage en assembleur. Comment j''en suis arrivé là, de la complétion de ligne à l''agent à qui je donne l''architecture, les habitudes que les LLM m''ont fait prendre et qui améliorent le code, et ce qu''ils font maintenant de mes tickets, de mes mails et de Slack.'
ogDescription: 'De la complétion de ligne à l''agent qui code sous mon architecture : six mois sans taper de code, et un code meilleur qu''avant.'
keywords: LLM, IA, agents de code, complétion de code, qualité logicielle, tests, changelog, refactoring, architecture logicielle, support, productivité, développeur
summary: 'Six mois sans taper une ligne de code, hors casse-têtes et assembleur. Le chemin en quatre étapes, de la complétion de ligne à l''agent qui code sous mon architecture, les habitudes qui ont amélioré mon code en trois mois, et le même principe appliqué aux tickets de support, aux mails et à Slack.'
---

Cela fait six mois que je n'ai pas écrit une ligne de code. Enfin, presque : j'en écris encore pour le plaisir, quand j'ai envie d'un casse-tête, ou pour débugger de l'assembleur. Tout le reste, au travail comme sur mes projets personnels, c'est un LLM qui le tape.

Je ne m'y suis pas mis par conviction, ni du jour au lendemain. C'est devenu, étape après étape, beaucoup plus rapide et plus efficace que de taper moi-même.

## Étape 1 : finir la ligne

Au début, le LLM complétait la ligne en cours. Je commençais, il devinait la fin. C'était de l'autocomplétion améliorée, mais déjà efficace : il prédisait le plus souvent exactement ce que j'allais écrire.

## Étape 2 : écrire les lignes suivantes

Puis il s'est mis à proposer les lignes d'après, et c'est là que j'ai commencé à gagner du temps pour de bon.

J'écris une requête SQL : il propose la structure Go qui reçoit le résultat, avec les bons champs et les bons types, d'après la requête juste au-dessus. Je colle la réponse JSON d'une API : il en tire la structure correspondante, tags compris. Rien de difficile dans ce code. Seulement des dizaines de lignes fastidieuses, où l'on se trompe par inattention.

## Étape 3 : décrire une petite fonction

Ensuite, j'ai arrêté de commencer les lignes. Je décrivais une petite fonction, exactement : ce qu'elle reçoit, ce qu'elle rend, les cas limites. Le LLM l'écrivait, je la relisais.

Le mot important est « exactement ». À ce stade, je restais l'auteur : je savais quel code je voulais, je le faisais seulement taper par un autre.

## Étape 4 : donner l'architecture, laisser coder

Aujourd'hui, je ne décris plus les fonctions. Je donne l'architecture du code et l'architecture logicielle complète, et le LLM écrit le reste : les fonctions, les tests, les migrations, l'interface.

Je donne l'architecture parce que c'est là que les LLM font encore beaucoup d'erreurs. Laissés seuls, ils empilent : une abstraction de trop, une dépendance inutile, deux façons de faire la même chose à trois fichiers d'écart. Mes projets commencent donc par un document : le contexte, des principes numérotés, ce qui est interdit. « Le code décide, le modèle rédige. » « Rien ne part sans un geste de l'utilisateur. » J'en parlais déjà à propos des [conventions maison transformées en skills](encoding-house-conventions-into-ai-agent-skills-EN) : ce qui n'est pas écrit n'existe pas pour un agent.

Cela dit, il y a de moins en moins à écrire. Au début, j'utilisais énormément de skills pour aider le LLM à prendre ses décisions. Aujourd'hui, je n'en garde que deux ou trois, et ce sont des préférences métier : la manière de rédiger les changelogs et les messages de commit, plus deux ou trois choix d'architecture de déploiement, Ansible par exemple. Pour le reste, il n'a plus besoin qu'on lui tienne la main.

Quant aux erreurs d'architecture, elles ne dureront pas. Dans deux ans, ce ne sera sûrement plus le cas. Peut-être même dans six mois.

## Trois mois d'habitudes qui améliorent le code

J'ai déjà écrit que dans mon équipe [la qualité n'avait pas baissé](from-twelve-devs-to-four-EN). Je vais plus loin : depuis trois mois, mon code est meilleur qu'avant. Pas parce que le modèle écrit mieux que moi. Parce que tout ce que je savais devoir faire, et que je ne faisais pas faute de temps, ne coûte plus rien.

**Un bug corrigé, un test qui l'empêche de revenir.** Exemple récent : une frise horaire dont tous les repères s'empilaient au bord gauche de la page. La politique de sécurité du serveur bloquait les attributs `style` écrits dans le HTML, sans rien dire. La correction fait dix lignes. Le test livré avec parcourt toutes les pages et tous les scripts, et refuse le moindre style en ligne : ce bug ne reviendra pas, ses cousins non plus. Sur mon dernier projet, il y a près d'une centaine de fichiers de test, un pour deux fichiers de code. Je n'ai jamais eu ce ratio de ma vie.

**Le commit ne passe pas si les tests ne passent pas.** Un hook de pre-commit avec `gofmt`, `go vet` et les tests : tout le monde sait qu'il faut le faire, personne ne prend l'heure pour le régler. Là, c'est une phrase. Et il sert autant à l'agent qu'à moi : un LLM qui peut lancer les tests corrige ses erreurs seul, un LLM qui ne le peut pas « vérifie » en relisant son propre diff.

**Un changelog qui dit ce qui a changé.** Mes messages de commit ressemblaient souvent à « fix », « wip », « fix again ». Aujourd'hui, ils décrivent le comportement : « une lecture abandonnée par le client rend 499 au lieu d'une 500 ». Et le changelog est mis à jour dans le même commit. Je ne les écris pas, mais je les relis, et c'est devenu ma première revue de code : si la ligne ne décrit pas ce que j'ai demandé, quelque chose a dérapé.

**Les refactorings qu'on remettait à plus tard.** Sur [OhayōDojo](go-pocketbase-templ-htmx-EN), la liste des « un jour, il faudra » a été vidée en quelques jours : code mort supprimé, un seul formateur de montants pour les pages, les messages, les PDF et les CSV, une plomberie commune aux trois importateurs, le contrôle d'accès centralisé. Rien de difficile, seulement long et ennuyeux, donc jamais prioritaire. Quand un refactoring coûte vingt minutes de relecture au lieu de deux jours de frappe, il se fait.

**Des outils jetables pour vérifier.** Un faux serveur LLM et un faux serveur IMAP pour tester hors réseau, une commande qui prend des captures d'écran de l'interface, un `--dry-run` sur chaque tâche planifiée. Je ne construisais presque jamais ce genre d'outillage, trop cher pour ce qu'il rapporte. Il vient maintenant avec le projet.

**Jeter sans regret.** Début octobre, j'ai commencé un assistant personnel : 34 000 lignes de Go en un mois. À l'usage, l'idée de départ était mauvaise, et fin octobre j'ai tout repris de zéro. Avec 34 000 lignes écrites à la main, j'aurais défendu mon architecture par simple amour-propre. Quand le code ne coûte presque rien, on arrête de s'y attacher, et on prend de meilleures décisions.

## Le gros du travail : reprendre du legacy

Tout cela ne s'est pas passé sur des projets neufs et propres. Le gros de mon travail, ces derniers mois, a été de reprendre du code legacy dans un état assez lamentable.

Les raisons sont classiques. Plusieurs développeurs juniors sont passés dessus, puis plusieurs prestataires. Les spécifications étaient imprécises, voire inexistantes. Et l'architecture posait problème partout : une multitude de micro-logiciels, chacun chargé d'une tâche, chacun avec son code spaghetti. Plusieurs d'entre eux lisaient et modifiaient les mêmes bases de données, les mêmes tables, et tous dépendaient les uns des autres.

Bref, une architecture choisie parce qu'elle était à la mode, par des développeurs qui ne savaient pas vraiment ce qu'ils faisaient.

C'est ce chantier-là que j'ai mené sans taper de code.

## Et pas seulement le code

Le même glissement s'est produit en dehors de l'éditeur. J'utilise maintenant un LLM pour les tickets de support : il trie, ventile 80 % des réponses et les pré-rédige. Il me reste à valider. Même chose pour les mails, les questions sur Slack et les demandes de la hiérarchie.

Il rédige aussi tout ce qui entoure le code. Il résume les modifications dans les messages de commit et les changelogs. Il prépare ma revue quotidienne pour le daily meeting, à partir de ce qui a réellement été fait la veille et non de ce dont je me souviens. Et il produit une documentation et des notes de version pour les utilisateurs qui ne sont pas techniques : le même changement, raconté sans nom de fonction ni numéro de migration. Cette traduction-là, je ne la faisais presque jamais, et ce sont pourtant eux qui utilisent le logiciel.

C'est le schéma de l'étape 4, appliqué ailleurs : la machine prépare, je décide. Rien ne part sans que je l'aie lu.

## Je relis tout, je corrige de moins en moins

À l'heure actuelle, je relis encore tout le code généré, et les tests générés avec. Mais j'ai de moins en moins besoin de demander des modifications. La plupart du temps, je lis et je valide, sans rien changer.

Conséquence logique : j'utilise désormais Claude Code ou OpenCode en mode automatique. L'agent enchaîne ses actions sans me demander la permission à chaque étape, et je relis le résultat à la fin.

Une exception : quand je lui confie une grande architecture, je passe en mode plan. Je décris un plan bien précis, et il ne commence à coder qu'une fois ce plan arrêté.

Et je ne parle pas forcément des modèles les plus poussés, comme Opus 5.5 : c'est vrai aussi avec des modèles plus modestes. Pour trier mes mails, par exemple, c'est un Qwen 3.8 de 27 milliards de paramètres, servi par vLLM sur ma propre machine, qui fait le travail.

## Ce qui reste mon travail

Je ne code plus, mais je n'ai jamais autant fait de développement. Le travail s'est simplement déplacé.

Aujourd'hui, je passe surtout mon temps à mettre de l'ordre dans les demandes des utilisateurs, qui arrivent brouillonnes et en vrac, d'autant qu'il n'y a pas de chef de projet. Mon rôle ressemble de plus en plus à celui d'un product owner, ou plutôt à un mélange : chef de projet, architecte logiciel, un peu de gestion d'équipe, et toujours du lead dev, pour guider les équipes et empêcher les conneries architecturales que leur LLM génère parfois.

Car c'est assez étonnant à voir : un LLM est un outil puissant, et même les plus puissants, comme Fable, peuvent être très, très mal utilisés. Entre certaines mains, ils produisent des résultats aberrants : des architectures dégueulasses, un code bien crade, instable, impossible à faire évoluer et mal documenté.

Alors je décide ce qui se construit et ce qui ne se construit pas. Je tiens l'architecture. Je lis les diffs, les changelogs et les tests. Je dis non, souvent : non à la dépendance de trop, non à l'abstraction prématurée, non à la fonctionnalité que personne n'a demandée.

Aucune des habitudes de cette liste n'est nouvelle. Ce sont celles que tous les livres recommandent depuis vingt ans et que presque personne ne tient, parce que le temps manque toujours. Le LLM n'a pas inventé la rigueur, il l'a rendue gratuite.

## D'ailleurs, ce billet

Ce billet a été écrit par un LLM. Le plan, je l'ai fait à la main. Les faits, je les ai dictés. Mais les phrases sont de lui.

J'ai tout dicté en marchant sur mon tapis de marche, un doigt sur la touche espace de mon portable pour lui parler. De tout le clavier, c'est la seule touche qui a servi.

Et quand le reste du clavier me manque, il reste l'assembleur.
