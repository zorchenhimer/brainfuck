package bf

import (
    "errors"
    "fmt"
    "io/ioutil"
    "os"
)

type Engine struct {
    cellIdx     int
    cells       []int

    commandIdx  int
    commands    []Command

    nestLevel   int
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

func newError(offset int, program, message string) *BrainfuckError {
    return &BrainfuckError{
        err:      errors.New(message),
        offset:     offset,
        program:    program,
    }
}

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

func (e *Engine) Run() *BrainfuckError {
    for e.commandIdx = 0; e.commandIdx < len(e.commands); e.commandIdx++ {
        switch e.commands[e.commandIdx] {
            case C_INCPTR:
                e.cellIdx += 1

                // Add another cell if the index is larger than the number of cells.
                if e.cellIdx >= len(e.cells) {
                    e.cells = append(e.cells, 0)
                }
            case C_DECPTR:
                e.cellIdx -= 1

                // Error if cell index is below zero.
                if e.cellIdx < 0 {
                    return newError(e.commandIdx, e.String(), fmt.Sprintf("cellIdx below zero at %d.", e.commandIdx))
                }
            case C_INC:
                e.cells[e.cellIdx]++
            case C_DEC:
                e.cells[e.cellIdx]--
            case C_OUT:
                fmt.Printf("%c", e.cells[e.cellIdx])
            case C_ACC:
                var v int
                // Caveat: User must hit enter to continue
                if _, err := fmt.Scanf("%c", &v); err != nil {
                    return newError(e.commandIdx, e.String(), fmt.Sprintf("C_ACC failure: %s", err))
                }
                e.cells[e.cellIdx] = v
            case C_JMPFOR:
                // Don't jump forward unless zero
                if e.cells[e.cellIdx] != 0 {
                    continue
                }

                foundback := false
                var startIdx int

                // Look forward for C_JMPBAC
                for startIdx = e.commandIdx; e.commandIdx < len(e.commands); e.commandIdx++ {
                    if e.commands[e.commandIdx] == C_JMPBAC {
                        e.commandIdx++
                        foundback = true
                        break
                    
                    // TODO: nesting
                    } else if e.commands[e.commandIdx] == C_JMPFOR {
                        return newError(e.commandIdx, e.String(), fmt.Sprintf("Nested loops unsupported at %d", e.commandIdx))
                    }
                }

                // Didn't find C_JMPBAC
                if !foundback {
                    return newError(e.commandIdx, e.String(), fmt.Sprintf("Unmatched C_JMPFOR at %d", startIdx))
                }
            case C_JMPBAC:
                // Continue if current cell is zero
                if e.cells[e.cellIdx] == 0 {
                    continue
                }

                foundforwd := false
                var startIdx int

                // Look backwards for C_JMPFOR
                for startIdx = e.commandIdx; e.commandIdx > -1; e.commandIdx-- {
                    if e.commands[e.commandIdx] == C_JMPFOR {
                        foundforwd = true
                        break
                    }
                }

                // Went too far back
                if !foundforwd {
                    return newError(e.commandIdx, e.String(), fmt.Sprintf("Unmatched C_JMPBAC at %d", startIdx))
                }
        }
    }
    return nil
}

func (e *Engine) parseCommand() *BrainfuckError {
    switch e.commands[e.commandIdx] {
        case C_INCPTR:
            e.cellIdx += 1

            // Add another cell if the index is larger than the number of cells.
            if e.cellIdx >= len(e.cells) {
                e.cells = append(e.cells, 0)
            }
        case C_DECPTR:
            e.cellIdx -= 1

            // Error if cell index is below zero.
            if e.cellIdx < 0 {
                return newError(e.commandIdx, e.String(), fmt.Sprintf("cellIdx below zero at %d.", e.commandIdx))
            }
        case C_INC:
            e.cells[e.cellIdx]++
        case C_DEC:
            e.cells[e.cellIdx]--
        case C_OUT:
            fmt.Printf("%c", e.cells[e.cellIdx])
        case C_ACC:
            var v int
            // Caveat: User must hit enter to continue
            if _, err := fmt.Scanf("%c", &v); err != nil {
                return newError(e.commandIdx, e.String(), fmt.Sprintf("C_ACC failure: %s", err))
            }
            e.cells[e.cellIdx] = v
    }
    return nil
}

func (e *Engine) doLoop() *BrainfuckError {
    var startIdx int
    // Process loop if non-zero
    if e.commands[e.commandIdx] != 0 {
        for startIdx = e.commandIdx; e.commandIdx < len(e.commands); e.commandIdx++ {
            switch e.commands[e.commandIdx] {
                // Found another loop
                case C_JMPFOR:
                    // Down the rabbit hole
                    e.nestLevel++
                    if e.nestLevel > 20 {
                        return newError(e.commandIdx, e.String(), fmt.Sprintf("Hit max recursion level at %d", e.commandIdx))
                    }
                    e.doLoop()

                // Found the end to current loop
                case C_JMPBAC:
                    if e.commands[e.commandIdx] == 0 {
                        return nil
                    }

                    foundforwd := false
                    backsfound := 0

                    // Look backwards for C_JMPFOR
                    for startIdx = e.commandIdx; e.commandIdx > -1; e.commandIdx-- {
                        if e.commands[e.commandIdx] == C_JMPFOR {
                            if backsfound <= 0 {
                                foundforwd = true
                                break
                            } else {
                                backsfound--
                            }
                        } else if e.commands[e.commandIdx] == C_JMPBAC {
                            backsfound++
                        }
                    }

                    // Went too far back
                    if !foundforwd {
                        return newError(e.commandIdx, e.String(), fmt.Sprintf("Unmatched C_JMPBAC at %d", startIdx))
                    }

                default:
                    if err := e.parseCommand(); err != nil {
                        return err
                    }
            }
        }

    // Jump to end of loop otherwise
    } else {
        foundback := false
        var startIdx int

        // Look forward for C_JMPBAC
        for startIdx = e.commandIdx; e.commandIdx < len(e.commands); e.commandIdx++ {
            if e.commands[e.commandIdx] == C_JMPBAC {
                e.commandIdx++
                foundback = true
                break
            
            } else if e.commands[e.commandIdx] == C_JMPFOR {
                e.doLoop() // Rabbit hole time...
            }
        }

        // Didn't find C_JMPBAC
        if !foundback {
            return newError(e.commandIdx, e.String(), fmt.Sprintf("Unmatched C_JMPFOR at %d", startIdx))
        }
    }
    return nil
}

func (e *Engine) String() string {
    s := ""
    for _, c := range e.commands {
        s += string(c)
    }
    return s
}

