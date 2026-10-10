---
title: 'Pourquoi j''aime Go (Golang)'
date: '2024-06-01'
lang: fr
translation: why-i-like-go-golang-EN
description: 'Pourquoi j''aime Go : compilation rapide, gofmt, bibliothèque standard riche, gestion explicite des erreurs, goroutines et channels pour la concurrence, binaires statiques multiplateformes.'
ogDescription: 'Compilation rapide, gofmt, gestion explicite des erreurs, goroutines et binaires statiques multiplateformes.'
keywords: Go, Golang, langage Go, gofmt, goroutines, channels, concurrence, gestion des erreurs, binaires statiques, multiplateforme, Fyne
summary: 'Compilation rapide, gofmt, bibliothèque standard riche, gestion explicite des erreurs, goroutines et channels pour la concurrence, binaires statiques multiplateformes.'
---

## Le confort de développement

- Une compilation vraiment rapide
- Des messages d'erreur clairs
- Le formatage Go : les règles de formatage de Go garantissent une présentation homogène dans tout le code, et chaque développeur lit et comprend plus facilement le code des autres.
- La bibliothèque standard de Go est très bien fournie, on n'a pas besoin d'un gros framework web.
- La rétrocompatibilité du langage : la promesse de compatibilité de Go 1 porte sur le code source, ce qui veut dire qu'un programme écrit pour une ancienne version 1.x compile et tourne toujours avec une version plus récente, sans modification. Passer à une version plus récente est facile, et d'ailleurs aucune version 2.x n'est prévue.

## La gestion des erreurs

En Go, la gestion des erreurs est simple et il n'y a pas d'exceptions. Une erreur est un type comme un autre.

Go pousse à une culture où l'on traite chaque erreur au moment où elle survient. Au premier abord, ça peut paraître un peu verbeux, mais on est obligé de gérer les erreurs, et elles ne finissent pas sous le tapis.

Un exemple simple de gestion d'erreur :

```
package main

import (
    "errors"
    "fmt"
)

var ErrDivideByZero = errors.New("cannot divide by zero")

func Divide(a int, b int) (int, error) {
    if b == 0 {
        return 0, ErrDivideByZero
    }
    return a / b, nil
}

func main() {
    result, err := Divide(10, 2)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Printf("10 divided by 2 is %d\n", result)
}
```

## La concurrence

Go gère la concurrence directement dans le langage, avec les goroutines et les channels, et écrire du code concurrent devient simple.

Un exemple avec le pattern *worker pool* :

```
package main

import (
   "fmt"
   "sync"
   "time"
)

func main() {
   startTime := time.Now()
   totalJobs := 5
   jobs := make(chan int, totalJobs)
   var workerGroup sync.WaitGroup

   for w := 1; w <= 2; w++ {
       workerGroup.Add(1)
       go worker(w, jobs, &workerGroup)
   }

   for job := 1; job <= totalJobs; job++ {
       jobs <- job
   }

   close(jobs)
   workerGroup.Wait()
   fmt.Printf("Total time %d\n", time.Since(startTime))
}

func worker(w int, jobs chan int, wg *sync.WaitGroup) {
   defer wg.Done()

   for job := range jobs {
       processJobs(w, job)
   }
}

func processJobs(w int, job int) {
   fmt.Printf("Worker %d : starting job %d\n", w, job)
   time.Sleep(time.Second)
   fmt.Printf("Worker %d : finished job %d\n", w, job)
}
```

## Multiplateforme

Les programmes Go tournent sur plusieurs plateformes, dont Windows, macOS et Linux, ce qui en fait un très bon choix pour écrire des applications multiplateformes. On peut les compiler en binaires statiques, c'est-à-dire que le code binaire est inclus dans l'exécutable final. Plus besoin d'édition de liens dynamique ni de dépendances à l'exécution : le programme est plus simple à déployer et à lancer.

## L'écosystème

L'écosystème Go n'est pas aussi mature que ceux de Python ou de Java. Il y a moins de bibliothèques, mais leur qualité est excellente. La seule *killer feature* qui manque à Go, c'est une interface graphique multiplateforme. Je sais que GTK et Qt peuvent faire l'affaire, mais je veux quelque chose qui soit vraiment dans l'esprit de Go. [Fyne](https://github.com/fyne-io/fyne) a l'air d'un bon début, mais lui non plus n'est pas encore mature.

Pour essayer Fyne :

```
go install fyne.io/fyne/v2/cmd/fyne_demo@latest
fyne_demo
```
