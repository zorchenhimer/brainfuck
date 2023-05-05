package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/zorchenhimer/brainfuck/dialects"
)

type Arguments struct {
	Input string `arg:"positional" help:"Input filename."`
	Debug bool `arg:"-d,--debug" help:"Turn on debugging."`
	Language string `arg:"-l,--language" help:"Specific brainfuck dialect to use.  See --dialects for full list."`
	Dialects bool `arg:"--dialects" help:"List all available dialects."`
	Translate string `arg:"-t,--translate" help:"Translate to new dialect.  Outputs to STDOUT."`
}

func main() {
	var args Arguments
	arg.MustParse(&args)

	if args.Dialects {
		d := []string{}
		for name, _ := range dialects.Dialects {
			d = append(d, name)
		}

		fmt.Println("Available dialects:", strings.Join(d, ", "))
		os.Exit(0)
	}

	if args.Input == "" {
		fmt.Println("No input file given")
		os.Exit(1)
	}

	file, err := os.Open(args.Input)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lang := "Brainfuck"
	if args.Language != "" {
		lang = args.Language
	} else {
		switch strings.ToLower(filepath.Ext(args.Input)) {
		case ".ff":
			lang = "FuckFuck"
		case ".ten":
			lang = "TenX"
		case ".pika":
			lang = "Pikalang"
		}
	}

	if args.Translate != "" {
		err = Translate(file, os.Stdout, lang, args.Translate)
		if err != nil {
			fmt.Println("Translation error: ", err)
			os.Exit(1)
		}
		fmt.Println("")
		return
	}

	var engine *Engine
	engine, err = Load(file, lang)

	if err != nil {
		log.Fatalf("Error loading source: %s", err)
	}

	if args.Debug {
		fmt.Println("Engine:", engine)
	}

	engine.Debug = args.Debug

	if bferr := engine.Run(); bferr != nil {
		fmt.Println("")
		log.Fatalf("Run failed: %s\n", bferr.String())
	}
}
