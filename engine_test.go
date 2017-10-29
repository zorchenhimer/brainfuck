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

var oldstdin *os.File = os.Stdin
var engine *testEngine = newTestEngine()

type testEngine struct {
    Engine
}

func newTestEngine() *testEngine {
    e := &testEngine{}

    e.cellIdx       = 0
    e.commandIdx    = 0
    e.nestLevel     = 0

    e.cells         = []int{0}
    e.commands      = []Command{}
    e.testing       = true
    e.stdout        = []byte{}

    return e
}

// Test case filenames
type testCase struct {
    Source          string
    Input           string
    Output          string
    OutputExpected  bool
}

func exists(path string) bool {
    _, err := os.Stat(path)
    if err == nil { return true }
    if os.IsNotExist(err) { return false }
    return true
}

func runFileTest(t *testing.T, tc testCase) {
    e := newTestEngine()
    t.Helper()
    var input *os.File
    var expected []byte
    var err error = nil

    if filepath.Ext(tc.Source) == ".ff" {
        err = e.FuckFuck(tc.Source)
    } else {
        err = e.Load(tc.Source)
    }

    if err != nil {
        t.Fatalf("Unable to load file: %s", err)
    }

    // Check for input commands and skip those tests
    if len(tc.Input) == 0 {
        for _, cmd := range e.commands {
            if cmd == C_ACC {
                fmt.Sprintf("Skipping %s due to finding an input command.", tc.Source)
                t.Skipf("Skipping %s due to finding an input command.", tc.Source)
                continue
            }
        }
    } else {
        input, err = os.Open(tc.Input)
        if err != nil {
            t.Fatalf("Unable to open input file: %s", err)
        }
        defer input.Close()
        os.Stdin = input
    }

    expected, err = ioutil.ReadFile(tc.Output)
    if err != nil {
        t.Fatalf("Unable to load output file: %s", err)
    }

    if bferr := e.Run(); bferr != nil {
        t.Fatalf("Run returned an error: %s", bferr)
    }

    t.Logf("stdout: %s", e.stdout)
    if !bytes.Equal(e.stdout, expected) {
        t.Fatalf("\nUnexepected stdout: %q\n          Expected: %q", e.stdout, expected)
    }

    if input != nil {
        os.Stdin = oldstdin
    }
}

func TestBitwidth(t *testing.T) {
    if err := engine.Load("testing/bitwidth.b"); err != nil {
        t.Fatalf("Unable to load file: %s", err)
    }

    if err := engine.Run(); err != nil {
        t.Fatalf("Run returned an error: %s", err)
    }

    t.Logf("bitwidth stdout: %s", engine.stdout)
}

func TestHello(t *testing.T) {
    runFileTest(t, testCase{
        Source: "testing/Hello.b",
        Output: "testing/Hello.out",
    })
}

func TestHello2(t *testing.T) {
    runFileTest(t, testCase{
        Source: "testing/Hello2.b",
        Output: "testing/Hello2.out",
    })
}

func TestCollatz(t *testing.T) {
    runFileTest(t, testCase{
        Source: "testing/Collatz.b",
        Output: "testing/Collatz.out",
        Input: "testing/Collatz.in",
    })
}

func TestCounter(t *testing.T) {
    runFileTest(t, testCase{
        Source: "testing/Counter.b",
        Output: "testing/Counter.out",
    })
}
