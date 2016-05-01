package main

import (
    "fmt"
    "log"
    "os"

    "./bf"
)

func main() {
    if len(os.Args) == 1 {
        log.Fatal("Missing input file.")
    } else if len(os.Args) > 2 {
        log.Fatal("Too many input files.")
    }

    engine := bf.NewEngine()
    err := engine.Load(os.Args[1])
    if err != nil {
        log.Fatalf("Error loading source: %s", err)
    }

    if err := engine.Run(); err != nil {
        fmt.Println("")
        log.Fatalf("Run failed: %s\n", err.String())
    }
}
