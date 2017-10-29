package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
)

func main() {
    if len(os.Args) == 1 {
        log.Fatal("Missing input file.")
    } else if len(os.Args) > 2 {
        log.Fatal("Too many input files.")
    }

    var err error
    engine := NewEngine()
    if filepath.Ext(os.Args[1]) == ".ff" {
        err = engine.FuckFuck(os.Args[1])
    } else {
        err = engine.Load(os.Args[1])
    }

    if err != nil {
        log.Fatalf("Error loading source: %s", err)
    }

    if bferr := engine.Run(); bferr != nil {
        fmt.Println("")
        log.Fatalf("Run failed: %s\n", bferr.String())
    }
}
