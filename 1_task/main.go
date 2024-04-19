package main

import (
    "fmt"
    "sync"
)

func factorial(n int, wg *sync.WaitGroup, ch chan int) {
    defer wg.Done()
    result := 1
    for i := 1; i <= n; i++ {
        result *= i
    }
    ch <- result
}

func main() {
    number := 10
    numWorkers := 4

    wg := &sync.WaitGroup{}
    ch := make(chan int, numWorkers)

    step := number / numWorkers
    for i := 0; i < numWorkers; i++ {
        start := i*step + 1
        end := (i + 1) * step
        if i == numWorkers-1 {
            end = number
        }
        wg.Add(1)
        go factorial(end-start+1, wg, ch)
    }

    go func() {
        wg.Wait()
        close(ch)
    }()

    result := 1
    for v := range ch {
        result *= v
    }

    fmt.Printf("%d ning faktorialini hisoblab chiqildi: %d\n", number, result)
}
