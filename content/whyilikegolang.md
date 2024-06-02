Title: Why i like Go / Golang
Date: 2024-06-01
Lang: EN

## Development experience

- Really fast compilation
- Clear error message
- Go format : Go's formatting rules ensure consistency throughout the codebase, making it easier for developers to read and understand each other's code.
- The standard library in Go is very well-provided for, and does not require a heavy-duty web framework.
- Retro compatibility of the language : Go 1.x maintains binary compatibility with previous versions,
which means that programs compiled with older Go versions can still run without modification.
Migration to newer version is easy, and actually there is no plan for a version 2.x.

## Error management

Go error management is simple and have no exception. Errors are just a type like others.

Go encourages a culture of error management that forces you to handle as they arise.
At first glance, it may seem a bit verbose, but it forces error handling and prevents errors from being swept under the rug.

A sample exemple of error handling :

    :::golang
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



## Concurrency

The concurrency : Go has built-in concurrency support through its goroutine and channel mechanisms, which makes it easy to write concurrent code.

A example with a worker pool pattern :

    :::golang
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
       fmt.Printf("Worker %d : starting job %d\n", worker, job)
       time.Sleep(time.Second)
       fmt.Printf("Worker %d : finished job %d\n" , worker , job)
    }

## Cross compatibility

Go programs can run on multiple platforms, including Windows, macOS, and Linux, making it a great choice for building cross-platform applications.
They can be compiled into static binaries, which means that the binary code is included in the final executable.
This eliminates the need for dynamic linking or runtime dependencies, making it easier to deploy and run.

## The ecosystem

While the Goang ecosystem is not as mature as Python or Java, there are fewer libraries, but the quality is awesome.
The only killer feature missing to Go is a multiplatform GUI, i know GTK and Qt can be an option, but something truely
in the spirit of GO. [https://github.com/fyne-io/fyne](Fyne) look like a good start, but also not yet mature.

Just try fyne :

    go install fyne.io/fyne/v2/cmd/fyne_demo@latest
    fyne_demo
