---
title: Why I like Go / Golang
date: '2024-06-01'
lang: en
description: 'Why I like Go: fast compilation, gofmt, a rich standard library, explicit error handling, goroutines and channels for concurrency, and static cross-platform binaries.'
ogDescription: Fast compilation, gofmt, explicit error handling, goroutines and static cross-platform binaries.
summary: Fast compilation, gofmt, a rich standard library, explicit error handling, goroutines and channels for concurrency, and static cross-platform binaries.
---

## Development experience

- Really fast compilation
- Clear error messages
- Go format: Go's formatting rules ensure consistency throughout the codebase, making it easier for developers to read and understand each other's code.
- The standard library in Go is very well provided, and does not require a heavy-duty web framework.
- Backward compatibility of the language: Go 1.x maintains binary compatibility with previous versions, which means that programs compiled with older Go versions can still run without modification. Migration to newer versions is easy, and actually there is no plan for a version 2.x.

## Error management

Go error management is simple and has no exceptions. Errors are just a type like any other.

Go encourages a culture of error management that forces you to handle errors as they arise. At first glance, it may seem a bit verbose, but it forces error handling and prevents errors from being swept under the rug.

A simple example of error handling:

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

## Concurrency

Concurrency: Go has built-in concurrency support through its goroutine and channel mechanisms, which makes it easy to write concurrent code.

An example with a worker pool pattern:

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

## Cross compatibility

Go programs can run on multiple platforms, including Windows, macOS, and Linux, making it a great choice for building cross-platform applications. They can be compiled into static binaries, which means that the binary code is included in the final executable. This eliminates the need for dynamic linking or runtime dependencies, making it easier to deploy and run.

## The ecosystem

While the Go ecosystem is not as mature as Python or Java, there are fewer libraries, but the quality is awesome. The only killer feature missing from Go is a multiplatform GUI. I know GTK and Qt can be options, but I want something truly in the spirit of Go. [Fyne](https://github.com/fyne-io/fyne) looks like a good start, but it is not yet mature either.

Just try Fyne:

```
go install fyne.io/fyne/v2/cmd/fyne_demo@latest
fyne_demo
```
