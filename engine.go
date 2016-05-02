package main

import (
    "fmt"
    "io/ioutil"
    "os"
)

var debug bool = false

type Engine struct {
    cellIdx     int
    cells       []int

    commandIdx  int
    commands    []Command

    nestLevel   int

    programFilename string
}

type BrainfuckError struct {
    err         error
    offset      int
    program     string
}

type Command rune

const (
    C_INCPTR    Command = '>'   // Increment pointer
    C_DECPTR    Command = '<'   // Decrement pointer
    C_INC       Command = '+'   // Increment value
    C_DEC       Command = '-'   // Decrement value
    C_OUT       Command = '.'   // Print value
    C_ACC       Command = ','   // Accept vaule
    C_JMPFOR    Command = '['   // Jump forward
    C_JMPBAC    Command = ']'   // Jump backwards
)

func (b *BrainfuckError) String() string {
    return b.Error() + "\n" + b.HelpString()
}

func (b *BrainfuckError) Error() string {
    return b.err.Error()
}

// FIXME: Breaks with programs longer than the terminal is wide.
func (b *BrainfuckError) HelpString() string {
    s := ""
    for i := 0; i < b.offset; i++ {
        s += " "
    }

    return b.program + "\n" + s + "^ here"
}

func NewEngine() *Engine {
    e := &Engine{
        cellIdx:    0,
        commandIdx: 0,
        nestLevel:  0,

        cells:      []int{0},
        commands:   []Command{},
    }
    return e
}

func (e *Engine) Load(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    rawBytes, err := ioutil.ReadAll(file)
    if err != nil {
        return err
    }
    e.programFilename = filename

    for _, b := range rawBytes {
        switch Command(b) {
            case C_INCPTR, C_DECPTR, C_INC, C_DEC, C_OUT, C_ACC, C_JMPFOR, C_JMPBAC:
                e.commands = append(e.commands, Command(b))
            default:
                // Ignore all other characters
                continue
        }
    }
    return nil
}

func (e *Engine) newError(message string) *BrainfuckError {
    return &BrainfuckError{
        err:        fmt.Errorf("[%s:%d] %s", e.programFilename, e.commandIdx, message),
        offset:     e.commandIdx,
        program:    e.String(),
    }
}

func (e *Engine) Run() *BrainfuckError {
    for e.commandIdx = 0; e.commandIdx < len(e.commands); e.commandIdx++ {
        switch e.commands[e.commandIdx] {
            
            // Increment the pointer
            case C_INCPTR:
                e.cellIdx += 1

                // Add another cell if the index is larger than the number of cells.
                if e.cellIdx >= len(e.cells) {
                    e.cells = append(e.cells, 0)
                }

            // Decrement the pointer
            case C_DECPTR:
                e.cellIdx -= 1

                // Error if cell index is below zero.
                if e.cellIdx < 0 {
                    debug = true
                    e.status("")
                    return e.newError("cellIdx below zero.")
                }

            // Increment the value
            case C_INC:
                e.cells[e.cellIdx]++

            // Decrement the value
            case C_DEC:
                e.cells[e.cellIdx]--

            // Print the cell's value
            case C_OUT:
                fmt.Printf("%c", e.cells[e.cellIdx])

            // Accept a new value
            case C_ACC:
                var v int
                // Caveat: User must hit enter to continue
                if _, err := fmt.Scanf("%c", &v); err != nil {
                    return e.newError(fmt.Sprintf("C_ACC failure: %s", err))
                }
                e.cells[e.cellIdx] = v

            // Jump forwards to matching close bracket if current cell is zero
            case C_JMPFOR:
                if e.cells[e.cellIdx] == 0 {
                    e.status("going to loop end")
                    e.gotoLoopEnd()
                    e.status("found loop end")
                } else {
                    e.status("continuing loop")
                }

            // Jump backwards to matching close open if current cell not zero
            case C_JMPBAC:
                if e.cells[e.cellIdx] != 0 {
                    e.status("going to loop start")
                    e.gotoLoopStart()
                } else {
                    e.status("Nonzero loop end")
                }

            // This shouldn't ever happen.
            default:
                return e.newError("Invalid command.")
        }
    }
    return nil
}

// Go to the end of the current loop
func (e *Engine) gotoLoopEnd() *BrainfuckError {
    lvl := 0
    tlvl := 0
    for e.commandIdx += 1; e.commandIdx < len(e.commands); e.commandIdx++ {
        switch e.commands[e.commandIdx] {
            case C_JMPFOR:
                lvl++
                tlvl++
            case C_JMPBAC:
                if lvl > 0 {
                    lvl--
                } else if lvl < 0 {
                    return e.newError("nest level too low in gotoLoopStart()")
                } else {
                    return nil
                }
        }
    }
    return e.newError(fmt.Sprintf("gotoLoopEnd finished without finding C_JMPBAC [%d]", tlvl))
}

// Go to the start of the current loop
func (e *Engine) gotoLoopStart() *BrainfuckError {
    lvl := 0
    tlvl := 0
    for e.commandIdx -= 1; e.commandIdx > -1; e.commandIdx-- {
        switch e.commands[e.commandIdx] {
            case C_JMPBAC:
                lvl++
                tlvl++
            case C_JMPFOR:
                if lvl > 0 {
                    lvl--
                } else if lvl < 0 {
                    return e.newError("nest level too low in gotoLoopStart()")
                } else {
                    return nil
                }
        }
    }
    return e.newError(fmt.Sprintf("gotoLoopStart finished without finding C_JMPFOR [%d]", tlvl))
}

// Print debugging information
func (e *Engine) status(message string) {
    if !debug || e.commandIdx < 62 { return }
    fmt.Printf("{%s} commandIdx: %d; command: %c; cellIdx: %d; cells: %v\n",
        message, e.commandIdx, e.commands[e.commandIdx], e.cellIdx, e.cells)
}

// Print the loaded program
func (e *Engine) String() string {
    s := ""
    for _, c := range e.commands {
        s += string(c)
    }
    return s
}

