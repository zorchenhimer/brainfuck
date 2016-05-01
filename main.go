package main

import (
    "fmt"
    "log"

    "./bf"
)

func main() {
    engine := bf.NewEngine()
    err := engine.Load("hello-world3.bf")
    if err != nil {
        log.Fatalf("Error loading source: %s", err)
    }

    if err := engine.Run(); err != nil {
        fmt.Println("")
        log.Fatalf("Run failed: %s\n", err.String())
    }
    
    fmt.Println("Done.")
}
