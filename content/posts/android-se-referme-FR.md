---
title: 'Android se referme : 12 testeurs nécessaires, et une vérification des devs'
date: '2026-09-17'
lang: fr
description: 'En rouvrant mon compte Google Play, j''ai découvert les nouvelles règles : 12 testeurs pendant 14 jours avant de publier, une vérification d''identité arbitraire, un sideloading bridé et un argument de sécurité qui vérifie l''auteur plutôt que le code. Mises bout à bout, ces mesures referment Android.'
ogDescription: '12 testeurs pendant 14 jours, vérification d''identité, sideloading bridé, bootloaders verrouillés : mises bout à bout, ces mesures referment Android.'
image: https://rvier.fr/images/android-se-referme.png
keywords: Android, Google Play, Play Console, 12 testeurs, test fermé, vérification des développeurs, sideloading, F-Droid, bootloader, Samsung, Xiaomi, Pixel, GrapheneOS, Play Integrity, attestation matérielle, applications bancaires
summary: 'Retour d''un développeur Android publiant depuis 2013 sur les nouvelles règles de Google : 12 testeurs pendant 14 jours, vérification d''identité des développeurs, sideloading limité avec biométrie et délai de 24 h, verrouillage des bootloaders chez Samsung et Xiaomi, et Play Integrity qui bloque les applications bancaires sur ROM alternative.'
---

Je développe des applications Android depuis longtemps, la première a dû être publiée vers 2013. Des applications personnelles, comme ForRunners ou MOrg, dans divers stores, F-Droid, Google Play Store, mais aussi de nombreuses applications professionnelles, certaines publiées, d'autres utilisées en interne.

<img src="../images/android-se-referme.png" alt="Vue du dessus d'un cercueil en bois avec le robot Android allongé dedans, le couvercle qu'on glisse sur ses pieds et un marteau qui enfonce un clou ; en légende, les quatre clous : 12 testeurs, vérification d'identité, sideloading bridé, l'auteur vérifié plutôt que le code" loading="lazy" width="1200" height="627">

Publier était simple, 25 $ servant principalement à identifier la personne, un APK, une description, et hop, c'était dans le store.

Sauf que mon ancien compte Play Store a été clos, pour une raison inconnue. À sa réouverture, je dois maintenant respecter les nouvelles règles du Google Play Store.

## 1er clou

Depuis décembre 2024, pour publier une application sur le Play Store, il faut d'abord trouver 12 personnes qui la testent pendant 14 jours consécutifs avant qu'elle puisse être disponible pour tous. Compliqué pour une application de niche ou à usage ponctuel.

## 2e clou

La vérification d'identité est au bon vouloir de Google, qui accepte ou non vos pièces d'identité, parfois refusées sans aucun motif.

## 3e clou

Pour installer des apps non certifiées, hors Play Store, la procédure limite le sideloading à 20 « amis », nécessite une identification biométrique sur le téléphone pour installer l'application, et impose un délai de 24 h pour valider, avec des messages qui dissuadent brutalement l'utilisateur.

## 4e clou

L'argument de sécurité. Google vérifie qui est la source des publications, pas le code. Contrairement à F-Droid, qui compile et signe les applications à partir des sources publiques, ce qui garantit une plus grande sécurité, ou tout du moins permet d'auditer facilement le code (surtout à l'heure des LLM).

## Le cercueil

L'étau se resserre. Chaque mesure prise isolément semble anodine, mais ensemble elles scellent définitivement le cercueil d'Android. Seul Google détiendra les clefs de qui est autorisé ou non à publier pour sa plateforme. Vous ne disposerez plus de votre matériel que si Google vous y autorise.

Et si vous avez la chance d'avoir un téléphone encore ouvert avec une ROM ouverte, chose de plus en plus difficile puisque les bootloaders sont verrouillés sur les smartphones récents (Samsung, Pixel, Xiaomi), vous aurez toujours cette horreur de service Play Integrity, requis par la plupart des applications bancaires, de paiement et de transport. Et là, ces applications refuseront tout simplement de fonctionner. « Pourquoi est-ce un problème, il reste les navigateurs et les sites web ? » Eh bien, certaines banques n'acceptent la double authentification que via leur application mobile.

Et non, il n'existe plus vraiment d'alternative solide.

### Sources

- [Test fermé, règle actuelle (12 testeurs, 14 jours, comptes créés après le 13 nov. 2023, questionnaire)](https://support.google.com/googleplay/android-developer/answer/14151465)
- [Guide communautaire Google « Everything about the 12 testers requirement »](https://support.google.com/googleplay/android-developer/community-guide/255621488)
- [Passage de 20 à 12 testeurs (déc. 2024), Android Authority](https://www.androidauthority.com/google-play-app-testing-requirement-3510580/)
- [Annonce de la vérification, 25 août 2025 (« 50 fois plus de malwares », aéroport, calendrier initial)](https://android-developers.googleblog.com/2025/08/elevating-android-security.html)
- [Concessions de mars 2026 (flux avancé, comptes distribution limitée)](https://android-developers.googleblog.com/2026/03/android-developer-verification.html)
- [Ouverture à tous les développeurs, mars 2026](https://android-developers.googleblog.com/2026/03/android-developer-verification-rolling-out-to-all-developers.html)
- [Calendrier détaillé de juin 2026 (30 septembre, quatre pays, sept stores, 2027)](https://android-developers.googleblog.com/2026/06/android-developer-verification.html)
- [Android developer verification, page officielle (25 $, pièce d'identité, D-U-N-S, adb exempté, blocage des apps non enregistrées, flux avancé avec 24 h et biométrie)](https://developer.android.com/developer-verification) et sa [FAQ](https://developer.android.com/developer-verification/guides/faq)
- [Aide Android Developer Console](https://support.google.com/android-developer-console/answer/16561738)
- [Android Authority (analyse XDA du code retiré du bootloader Samsung)](https://www.androidauthority.com/samsung-bootloader-unlocking-disabled-one-ui-8-3581366/)
- [Annonce officielle Xiaomi Community](https://new.c.mi.com/global/post/2251454)
- [GrapheneOS sur l'origine du partenariat avec Motorola et l'avenir incertain des Pixel](https://piunikaweb.com/2026/03/12/grapheneos-explains-motorola-partnership-origin-the-uncertain-future-of-pixels/)
- [Android Authority, mai 2025, attestation matérielle par défaut](https://www.androidauthority.com/google-play-integrity-hardware-attestation-3561592/)
- [Guide de compatibilité par attestation matérielle (à destination des banques)](https://grapheneos.org/articles/attestation-compatibility-guide)
- [Forum GrapheneOS, banque sans Play Integrity](https://discuss.grapheneos.org/d/29413-banking-without-the-play-integrity-api)
- [Fil du forum GrapheneOS (14 pages, témoignages)](https://discuss.grapheneos.org/d/17714-revolut-mobile-finance-not-supported-on-devices-with-custom-firmware-problem)

