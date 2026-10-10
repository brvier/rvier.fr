---
title: 'Octet pour octet : reproduire en Go un vieux format binaire non documenté'
date: '2026-06-05'
lang: fr
translation: byte-for-byte-reproducing-a-legacy-binary-format-EN
featured: true
description: 'Réécrire en Go un pipeline média legacy imposait de reproduire octet pour octet un format binaire non documenté : rétro-ingénierie de floor(min(mean(|s16|)/256, 127)) par fenêtre de 10 ms, retenue fractionnaire, et tests unitaires de conformité.'
ogDescription: 'Rétro-ingénierie d''un format de profil audio non documenté, retenue fractionnaire reproduite bug pour bug, et conformité prouvée par des tests à l''octet près.'
keywords: Go, Golang, rétro-ingénierie, reverse engineering, legacy, code legacy, format binaire, migration, réécriture, ffmpeg, PCM, tests unitaires, golden files
summary: 'Retrouver un format de profil audio non documenté, retenue fractionnaire comprise, bug pour bug, et prouver la conformité par des tests à l''octet près.'
---

Toute réécriture d'un système legacy finit par buter sur le même mur : un format de fichier, un message réseau ou une somme de contrôle que personne n'a documenté, que l'auteur d'origine a emporté avec lui en quittant l'entreprise, et dont une douzaine de consommateurs en aval dépendent *exactement tel qu'il est*. Voici l'histoire de l'un d'eux : un petit fichier annexe (*sidecar*) de « profil » audio, que j'ai dû reproduire octet pour octet en réécrivant en Go un démon de conversion média vieux de vingt ans.

## Le format que personne n'a mis par écrit

L'ancien pipeline transcode des enregistrements broadcast, et à côté de chaque fichier média converti il écrit un petit fichier annexe `.ap` (« audio profile »). Les outils en aval s'en servent pour dessiner des formes d'onde et pour se positionner sur les passages audio intéressants sans décoder le média. Il n'y a ni en-tête, ni nombre magique, ni spécification, seulement un flux brut d'octets, un pour chaque tranche de 10 millisecondes d'audio source.

Comme les consommateurs comparent ces fichiers et les mettent en cache, « à peu près pareil » ne suffisait pas : la réécriture en Go devait produire les mêmes octets que l'ancien binaire C, sur toutes les entrées, pour toujours. Le code source de l'ancien binaire était disponible mais à peine lisible, et le format lui-même n'existait qu'implicitement, dans une seule fonction. Une fois retrouvé par rétro-ingénierie et mis par écrit, il donne ceci :

```
byte = floor( min( mean(|s16 sample|) / 256 , 127 ) )
```

Pour chaque fenêtre de 10 ms : on prend les échantillons 16 bits signés entrelacés, on fait la moyenne de leurs valeurs absolues, on divise par 256, on plafonne à 127 et on arrondit à l'entier inférieur. Un octet en sortie. C'est simple, et impossible à deviner sans lire le code d'origine, parce que plusieurs variantes proches (RMS au lieu de la moyenne des valeurs absolues, plafond à 128, arrondi au plus proche au lieu de l'arrondi inférieur) produisent des formes d'onde plausibles mais subtilement fausses.

## Le diable se cache dans la retenue fractionnaire

Le vrai piège, c'était la longueur des fenêtres. Une fenêtre de 10 ms en mono à 44 100 Hz fait 441 échantillons, pas de souci. En stéréo à 44 100 Hz elle en fait 882, et pour certaines combinaisons de fréquence et de nombre de canaux, la fenêtre idéale est *fractionnaire*. L'ancien binaire ne rééchantillonnait pas et n'accumulait pas de flottants : il tenait un accumulateur entier en millièmes (une variable joliment nommée `lProfilArrondi`), y ajoutait la partie fractionnaire à chaque fenêtre, et chipait un échantillon de plus chaque fois que l'accumulateur franchissait 1 000.

Il faut reproduire ça exactement. Si on se trompe, les fenêtres dérivent lentement par rapport à l'original : les deux fichiers concordent au début et divergent au bout de quelques secondes. C'est la pire forme d'erreur, celle qu'un coup d'œil rapide ne repère pas. La version Go reproduit l'accumulateur tel quel :

```
spw := float64(channels*rate) / 100.0    // samples per 10 ms window
base := int(spw)                         // whole samples per window
frac := int((spw - float64(base)) * 1000) // thousandths, carried
carry := 0

nextWindow := func() int {
    n := base
    carry += frac
    if carry > 999 {
        n++
        carry -= 1000
    }
    return n
}
```

Un accumulateur entier en millièmes, est-ce comme ça que *moi* je répartirais les échantillons fractionnaires ? Non. Mais c'est comme ça que le format fonctionne, donc c'est ce que fait le code, avec un commentaire qui renvoie à la fonction legacy qu'il reproduit. Quand on reproduit un format, la fidélité l'emporte toujours sur l'élégance ; c'est il y a vingt ans qu'il fallait être malin.

Même raisonnement pour les bords : la dernière fenêtre incomplète est jetée sans rien dire, parce que c'est ce que faisait le binaire. Une « amélioration » à cet endroit serait une régression.

## Le prouver : des tests identiques à l'octet près

Une réécriture de ce genre ne vaut que ce que vaut sa preuve de conformité. Deux niveaux de tests portent cette charge.

D'abord, des tests unitaires verrouillent chaque propriété de l'algorithme avec des cas qu'on peut calculer à la main (amplitude constante, plafonnement, dernière fenêtre jetée) :

```
// 1 second of stereo 48 kHz at amplitude 25600
// -> 100 windows, each byte = 25600/256 = 100.
in := pcm(48000*2, 25600)
Compute(&out, bytes.NewReader(in), 48000, 2)
// assert: exactly 100 bytes, all equal to 100

// Negative full scale (-32768) -> |s| = 32768 -> 32768/256 = 128 -> clamped to 127.
// (32767 would not test the clamp: floor(127.99) is already 127.)

// 661 mono samples at 44100 (1.5 windows) -> 1 byte emitted,
// trailing 220 samples discarded, like the binary.
```

Ensuite, le niveau décisif : la comparaison avec des fichiers de référence (*golden files*). On fait tourner l'ancien binaire et l'implémentation Go sur les mêmes enregistrements réels, et on compare les sorties avec `cmp`. Ni « similaires », ni « dans la tolérance » : identiques, sinon le build échoue. La première exécution de cette comparaison est aussi le meilleur outil de rétro-ingénierie qui soit. Chaque divergence pointe vers une hypothèse fausse, et la position du premier octet différent indique en général *laquelle* (un *off-by-one* provoque une dérive, une erreur de formule est fausse dès l'octet zéro).

Côté entrée, tout reste déterministe parce que les deux implémentations reçoivent le même PCM décodé : la version Go lance ffmpeg en sous-processus (`-f s16le`, fréquence d'échantillonnage et nombre de canaux de la source) exactement comme l'ancien pipeline décodait, si bien que la comparaison isole le calcul du profil de toute dérive du décodeur.

## À retenir pour votre prochaine réécriture legacy

- **Le format, c'est le comportement, pas l'intention.** Écrivez la formule que l'ancien code calcule réellement, avec ses arrondis, ses plafonds et ses bizarreries sur les cas limites, avant de toucher à la nouvelle implémentation. La compatibilité bug pour bug est une fonctionnalité.
- **Nommez le fantôme.** Gardez des références à l'ancien code dans les commentaires (`lProfilArrondi` survit dans les nôtres). Le prochain mainteneur doit savoir quelles décisions sont des murs porteurs de la compatibilité, et lesquelles il peut changer librement.
- **Identique à l'octet près, sinon rien.** Des fichiers de référence produits par le vrai binaire legacy et comparés avec `cmp` ne coûtent pas cher et mettent fin à toutes les discussions. Les comparaisons avec tolérance cachent la dérive.
- **Méfiez-vous des divergences lentes.** Les bugs dangereux sont ceux qui concordent sur les N premières fenêtres. Testez sur des entrées longues, et sur les combinaisons de fréquence et de canaux qui donnent des fenêtres fractionnaires, pas seulement sur les cas commodes.

Depuis, le démon Go a remplacé l'ancien binaire en production. Les fichiers `.ap` qu'il écrit sont impossibles à distinguer des anciens, et c'est exactement le but : la meilleure migration est celle que personne en aval ne peut détecter.
