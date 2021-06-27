package main

import (
	"fmt"
	//"io/ioutil"
	"io"
	"os"
	"bufio"
	"strings"
)

type Engine struct {
	cellIdx int
	cells   []int

	commandIdx int
	commands   []Command

	nestLevel int

	programFilename string

	testing bool
	stdout  []byte

	Debug bool

	Stdin io.Reader
	Stdout io.Writer
}

type BrainfuckError struct {
	err     error
	offset  int
	program string
}

type Command rune

const (
	C_INCPTR Command = '>' // Increment pointer
	C_DECPTR Command = '<' // Decrement pointer
	C_INC    Command = '+' // Increment value
	C_DEC    Command = '-' // Decrement value
	C_OUT    Command = '.' // Print value
	C_ACC    Command = ',' // Accept vaule
	C_JMPFOR Command = '[' // Jump forward
	C_JMPBAC Command = ']' // Jump backwards
	// These are not implemented yet
	C_DEBUG Command = '#' // Debug dump
	C_DATA  Command = '!' // Data section to read from for input
)

// https://github.com/MiffOttah/fuckfuck
type FFCommand string

const (
	FF_INCPTR FFCommand = "ass"
	FF_DECPTR FFCommand = "bitch"
	FF_INC    FFCommand = "cunt"
	FF_DEC    FFCommand = "damn"
	FF_OUT    FFCommand = "dick"
	FF_ACC    FFCommand = "fuck"
	FF_JMPFOR FFCommand = "shit"
	FF_JMPBAC FFCommand = "twat"
)

type TenXCommand rune

const (
	TX_INCPTR TenXCommand = '\u2715'
	TX_DECPTR TenXCommand = '\u00D7'
	TX_INC    TenXCommand = '\u0058'
	TX_DEC    TenXCommand = '\u0087'
	TX_OUT    TenXCommand = '\u2716'
	TX_ACC    TenXCommand = '\U0001D4CD'
	TX_JMPFOR TenXCommand = '\u2717'
	TX_JMPBAC TenXCommand = '\u2613'
	// These are not implemented yet
	TX_DEBUG TenXCommand = '\u24CD'
	TX_DATA  TenXCommand = '\u2612'
)

// https://www.dcode.fr/pikalang-language
type PikaCommand string

const (
	PK_INCPTR PikaCommand = "pipi"
	PK_DECPTR PikaCommand = "pichu"
	PK_INC    PikaCommand = "pi"
	PK_DEC    PikaCommand = "ka"
	PK_OUT    PikaCommand = "pikachu"
	PK_ACC    PikaCommand = "pikapi"
	PK_JMPFOR PikaCommand = "pika"
	PK_JMPBAC PikaCommand = "chu"
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

func (b *BrainfuckError) HasError() bool {
	return b.err != nil
}

func (e *Engine) reset() {
	e.cellIdx = 0
	e.commandIdx = 0
	e.nestLevel = 0
	e.cells = []int{0}
	//e.commands = []Command{}
	e.programFilename = ""
}

func Brainfuck(input io.Reader) (*Engine, error) {
	reader := bufio.NewReader(input)
	engine := &Engine{
		cells:    []int{0},
		commands: []Command{},
	}

	var err error
	for err == nil {
		var r rune
		r, _, err = reader.ReadRune()
		if err != nil && err != io.EOF {
			return nil, err
		}

		switch Command(r) {
		case C_INCPTR, C_DECPTR, C_INC, C_DEC, C_OUT, C_ACC, C_JMPFOR, C_JMPBAC:
			engine.commands = append(engine.commands, Command(r))
		}
	}

	return engine, nil
}

func FuckFuck(input io.Reader) (*Engine, error) {
	reader := bufio.NewReader(input)
	engine := &Engine{
		cells:    []int{0},
		commands: []Command{},
	}

	var err error
	for err == nil {
		var word string
		word, err = reader.ReadString(' ')
		if err != nil && err != io.EOF {
			return nil, err
		}
		word = strings.Trim(strings.ReplaceAll(word, "\r", ""), "\n\t ")

		for _, w := range strings.Split(word, "\n") {
			var cmd Command
			switch FFCommand(strings.ToLower(w)) {
			case FF_INCPTR:
				cmd = C_INCPTR
			case FF_DECPTR:
				cmd = C_DECPTR
			case FF_INC:
				cmd = C_INC
			case FF_DEC:
				cmd = C_DEC
			case FF_OUT:
				cmd = C_OUT
			case FF_ACC:
				cmd = C_ACC
			case FF_JMPFOR:
				cmd = C_JMPFOR
			case FF_JMPBAC:
				cmd = C_JMPBAC
			default:
				continue
			}
			engine.commands = append(engine.commands, cmd)
		}
	}

	return engine, nil
}

func TenX(input io.Reader) (*Engine, error) {
	reader := bufio.NewReader(input)
	engine := &Engine{
		cells:    []int{0},
		commands: []Command{},
	}

	var err error
	for err == nil {
		var r rune
		r, _, err = reader.ReadRune()
		if err != nil && err != io.EOF {
			return nil, err
		}

		var cmd Command
		switch TenXCommand(r) {
		case TX_INCPTR:
			cmd = C_INCPTR
		case TX_DECPTR:
			cmd = C_DECPTR
		case TX_INC:
			cmd = C_INC
		case TX_DEC:
			cmd = C_DEC
		case TX_OUT:
			cmd = C_OUT
		case TX_ACC:
			cmd = C_ACC
		case TX_JMPFOR:
			cmd = C_JMPFOR
		case TX_JMPBAC:
			cmd = C_JMPBAC
		default:
			continue
		}
		engine.commands = append(engine.commands, cmd)
	}

	return engine, nil
}

func Pikachu(input io.Reader) (*Engine, error) {
	reader := bufio.NewReader(input)
	engine := &Engine{
		cells:    []int{0},
		commands: []Command{},
	}

	var err error
	for err == nil {
		var word string
		word, err = reader.ReadString(' ')
		if err != nil && err != io.EOF {
			return nil, err
		}
		word = strings.Trim(strings.ReplaceAll(word, "\r", ""), "\n\t ")

		for _, w := range strings.Split(word, "\n") {
			var cmd Command
			switch PikaCommand(w) {
			case PK_INCPTR:
				cmd = C_INCPTR
			case PK_DECPTR:
				cmd = C_DECPTR
			case PK_INC:
				cmd = C_INC
			case PK_DEC:
				cmd = C_DEC
			case PK_OUT:
				cmd = C_OUT
			case PK_ACC:
				cmd = C_ACC
			case PK_JMPFOR:
				cmd = C_JMPFOR
			case PK_JMPBAC:
				cmd = C_JMPBAC
			default:
				if engine.Debug {
					fmt.Fprintf(engine.Stdout, "%q ignored", w)
				}
				continue
			}
			engine.commands = append(engine.commands, cmd)
		}
	}
	return engine, nil
}

func (e *Engine) newError(message string) *BrainfuckError {
	return &BrainfuckError{
		err:     fmt.Errorf("[%s:%d] %s", e.programFilename, e.commandIdx, message),
		offset:  e.commandIdx,
		program: e.String(),
	}
}

func (e *Engine) Run() *BrainfuckError {
	if e.Stdout == nil {
		e.Stdout = os.Stdout
	}

	if e.Stdin == nil {
		e.Stdin = os.Stdin
	}

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
				e.Debug = true
				e.status("")
				return e.newError("cellIdx below zero.")
			}

		// Increment the value
		case C_INC:
			if e.cellIdx < 0 {
				return e.newError("C_INC: negative cellIdx")
			}
			if e.cellIdx >= len(e.cells) {
				return e.newError("C_INC: cellIdx out of range")
			}
			e.cells[e.cellIdx]++

		// Decrement the value
		case C_DEC:
			if e.cellIdx < 0 {
				return e.newError("C_DEC: negative cellIdx")
			}
			if e.cellIdx >= len(e.cells) {
				return e.newError("C_DEC: cellIdx out of range")
			}
			e.cells[e.cellIdx]--

		// Print the cell's value
		case C_OUT:
			fmt.Fprintf(e.Stdout, "%c", e.cells[e.cellIdx])

		// Accept a new value
		case C_ACC:
			var v int
			// Caveat: User must hit enter to continue
			if _, err := fmt.Fscanf(e.Stdin, "%c", &v); err != nil {
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
			if e.cells[e.cellIdx] < -1000 {
				return e.newError("C_JMPBAC cell value less than -1000")
			}

			if e.cells[e.cellIdx] != 0 {
				e.status("going to loop start")
				e.gotoLoopStart()
			} else {
				e.status("Nonzero loop end")
			}

		// This shouldn't ever happen.
		default:
			return e.newError(fmt.Sprintf("Invalid command: %q", e.commands[e.commandIdx]))
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
	if !e.Debug || e.commandIdx < 62 {
		return
	}
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
