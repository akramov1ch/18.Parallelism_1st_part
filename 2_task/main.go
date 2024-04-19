package main

import (
    "fmt"
    "os"
    "sync"
)

func readFile(filePath string, results chan<- []byte, errors chan<- error, wg *sync.WaitGroup) {
    defer wg.Done()
    file, err := os.Open(filePath)
    if err != nil {
        errors <- err
        return
    }
    defer file.Close()

    fileInfo, err := file.Stat()
    if err != nil {
        errors <- err
        return
    }

    fileSize := fileInfo.Size()
    buffer := make([]byte, fileSize)

    bytesRead, err := file.Read(buffer)
    if err != nil {
        errors <- err
        return
    }

    if int64(bytesRead) != fileSize {
        errors <- fmt.Errorf("could not read the entire file: %s", filePath)
        return
    }

    results <- buffer
}

func main() {
    filePaths := []string{
        "2_task/file1.txt",
        "2_task/file2.txt",
        "2_task/file3.txt",
    }

    results := make(chan []byte)
    errors := make(chan error)
    wg := &sync.WaitGroup{}

    for _, filePath := range filePaths {
        wg.Add(1)
        go readFile(filePath, results, errors, wg)
    }

    go func() {
        wg.Wait()
        close(results)
        close(errors)
    }()

    allResults := []byte{}
    collectedErrors := []error{}

    for {
        select {
        case result, ok := <-results:
            if ok {
                allResults = append(allResults, result...)
            } else {
                results = nil
            }
        case err, ok := <-errors:
            if ok {
                collectedErrors = append(collectedErrors, err)
            } else {
                errors = nil
            }
        }
        if results == nil && errors == nil {
            break
        }
    }

    fmt.Println("All results:")
    fmt.Println(string(allResults))

    if len(collectedErrors) > 0 {
        fmt.Println("Errors:")
        for _, err := range collectedErrors {
            fmt.Println(err)
        }
    }
}
