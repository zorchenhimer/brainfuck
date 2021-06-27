package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexflint/go-arg"
)

type Arguments struct {
	Input string `arg:"positional,required"`
	Debug bool `arg:"-d,--debug" help:"Turn on debugging"`
}

func main() {
	var args Arguments
	arg.MustParse(&args)

	fmt.Println("Input:", args.Input)

	file, err := os.Open(args.Input)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var engine *Engine

	switch strings.ToLower(filepath.Ext(args.Input)) {
	case ".ff":
		engine, err = FuckFuck(file)
	case ".ten":
		engine, err = TenX(file)
	case ".pika":
		engine, err = Pikachu(file)
	default:
		engine, err = Brainfuck(file)
	}

	if err != nil {
		log.Fatalf("Error loading source: %s", err)
	}

	if args.Debug {
		fmt.Println("Engine:", engine)
	}

	//engine.Debug = args.debug

	if bferr := engine.Run(); bferr != nil {
		fmt.Println("")
		log.Fatalf("Run failed: %s\n", bferr.String())
	}
}
