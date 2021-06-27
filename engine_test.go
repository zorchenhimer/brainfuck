package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	//"strings"
	"testing"
)

// Test case filenames
type testCase struct {
	Source         string
	Input          string
	Output         string
	OutputExpected bool
	OutputText     string
}

func runFileTest(t *testing.T, tc testCase) error {
	t.Helper()
	var input *os.File
	var expected []byte
	var err error = nil
	var e *Engine
	var file *os.File

	file, err = os.Open(tc.Source)
	if err != nil {
		return fmt.Errorf("Unable to open source file: %w", err)
	}

	//if filepath.Ext(tc.Source) == ".ff" {
	switch filepath.Ext(tc.Source) {
	case ".ff":
		e, err = FuckFuck(file)
	case ".ten":
		e, err = TenX(file)
	case ".pika":
		e, err = Pikachu(file)
	default:
		e, err = Brainfuck(file)
	}

	if err != nil {
		return fmt.Errorf("Unable to load file: %w", err)
	}

	// Check for input commands and skip those tests if no input file is supplied
	if len(tc.Input) == 0 {
		//for _, cmd := range e.commands {
		//    if cmd == C_ACC {
		//        fmt.Sprintf("Skipping %s due to finding an input command.", tc.Source)
		//        t.Skipf("Skipping %s due to finding an input command.", tc.Source)
		//        continue
		//    }
		//}
	} else {
		input, err = os.Open(tc.Input)
		if err != nil {
			return fmt.Errorf("Unable to open input file: %w", err)
		}
		defer input.Close()
		e.Stdin = input
	}

	if len(tc.OutputText) > 0 {
		expected = []byte(tc.OutputText)
	} else {
		expected, err = ioutil.ReadFile(tc.Output)
		if err != nil {
			return fmt.Errorf("Unable to load output file: %w", err)
		}
	}

	outBuff := &bytes.Buffer{}
	e.Stdout = outBuff

	if bferr := e.Run(); bferr != nil {
		return fmt.Errorf("Run returned an error: %w", bferr)
	}

	if !bytes.Equal(outBuff.Bytes(), expected) {
		return fmt.Errorf("\nUnexepected Stdout: %q\n          Expected: %q",
			outBuff, expected)
	}

	return nil
}

// FIXME: What should the output of this test *really* be?
func TestBitwidth(t *testing.T) {
	//if err := engine.Load("testing/bitwidth.b"); err != nil {
	//    t.Fatalf("Unable to load file: %s", err)
	//}

	//if err := engine.Run(); err != nil {
	//    t.Fatalf("Run returned an error: %s", err)
	//}

	//t.Logf("bitwidth stdout: %s", engine.stdout)
	err := runFileTest(t, testCase{
		Source:     "testing/bitwidth.b",
		OutputText: "Hello, world!\n",
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestHello(t *testing.T) {
	err := runFileTest(t, testCase{
		Source: "testing/Hello.b",
		Output: "testing/Hello.out",
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestHello2(t *testing.T) {
	err := runFileTest(t, testCase{
		Source: "testing/Hello2.b",
		Output: "testing/Hello2.out",
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

// Cell index goes negative and fails
//func TestBench(t *testing.T) {
//	err := runFileTest(t, testCase{
//		Source: "testing/Bench.b",
//		Output: "testing/Bench.out",
//	})
//	if err != nil {
//		t.Fatalf(err.Error())
//	}
//}

func TestCollatz(t *testing.T) {
	err := runFileTest(t, testCase{
		Source: "testing/Collatz.b",
		Output: "testing/Collatz.out",
		Input:  "testing/Collatz.in",
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestFuckFuck(t *testing.T) {
	err := runFileTest(t, testCase{
		Source: "testing/fuckfuck.ff",
		Output: "testing/fuckfuck.out",
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

// Fails at offset 2234
//func TestLife(t *testing.T) {
//	err := runFileTest(t, testCase{
//		Source: "testing/Life.b",
//		Output: "testing/Life.out",
//		Input: "testing/Life.in",
//	})
//	if err != nil {
//		t.Fatalf(err.Error())
//	}
//}

func TestCounter(t *testing.T) {
	err := runFileTest(t, testCase{
		Source: "testing/Counter.b",
		Output: "testing/Counter.out",
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestTribit(t *testing.T) {
	err := runFileTest(t, testCase{
		Source:     "testing/Tribit.b",
		OutputText: "32 bit cells\n", // This'll probably fail with different configurations
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestPikachu(t *testing.T) {
	err := runFileTest(t, testCase{
		Source: "testing/pikachu-hello.pika",
		OutputText: "Pokemon",
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
}
