package main

import (
	"fmt"
	"io"
	"os"
	"bufio"
	"strings"
	"unicode"

	"github.com/zorchenhimer/brainfuck/dialects"
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

func Load(reader io.Reader, dialect string) (*Engine, error) {
	d, ok := dialects.Dialects[dialect]
	if !ok {
		return nil, fmt.Errorf("Dialect %q doesn't exist", dialect)
	}

	if d.Type() == dialects.RuneLang {
		return runeLang(reader, d)
	}
	return wordLang(reader, d)
}

func Translate(reader io.Reader, writer io.Writer, source, dest string) error {
	engine, err := Load(reader, source)
	if err != nil {
		return fmt.Errorf("Unable to load source program: %w", err)
	}

	if source == dest {
		return fmt.Errorf("Both dialects are the same!")
	}

	d, ok := dialects.Dialects[dest]
	if !ok {
		return fmt.Errorf("Destination dialect %q doesn't exist", dest)
	}

	newMap := map[Command]string{}

	if d.Type() == dialects.WordLang {
		wordMap := d.(dialects.WordMap)
		for key, val := range wordMap {
			newMap[Command(val)] = key
		}
	} else {
		runeMap := d.(dialects.RuneMap)
		for key, val := range runeMap {
			newMap[Command(val)] = string(key)
		}
	}

	tl := []string{}
	for _, c := range engine.commands {
		tl = append(tl, newMap[c])
		//fmt.Fprintf(writer, "%s", newMap[c])
	}

	joinstr := ""
	if d.Type() == dialects.WordLang {
		joinstr = " "
	}
	_, err = fmt.Fprintf(writer, strings.Join(tl, joinstr))
	return err
}

func runeLang(input io.Reader, dialect dialects.Dialect) (*Engine, error) {
	reader := bufio.NewReader(input)
	engine := &Engine{
		cells:    []int{0},
		commands: []Command{},
	}

	runeMap := dialect.(dialects.RuneMap)

	var err error
	for err == nil {
		var r rune
		r, _, err = reader.ReadRune()
		if err != nil && err != io.EOF {
			return nil, err
		}

		if cmd, ok := runeMap[r]; ok {
			engine.commands = append(engine.commands, Command(cmd))
		}
	}

	if len(engine.commands) == 0 {
		return nil, fmt.Errorf("No commands read from source. Was the correct dialect used?")
	}
	return engine, nil
}

func wordLang(input io.Reader, dialect dialects.Dialect) (*Engine, error) {
	reader := bufio.NewReader(input)
	engine := &Engine{
		cells:    []int{0},
		commands: []Command{},
	}

	wordMap := dialect.(dialects.WordMap)

	var err error
	for err == nil {
		word := strings.Builder{}
		for {
			var r rune
			r, _, err = reader.ReadRune()
			if err != nil {
				if err != io.EOF {
					return nil, err
				}
				break
			}

			if unicode.IsSpace(r) {
				break
			}

			_, err = word.WriteRune(r)
			if err != nil {
				return nil, err
			}
		}

		if cmd, ok := wordMap[word.String()]; ok {
			engine.commands = append(engine.commands, Command(cmd))
		}
	}

	if len(engine.commands) == 0 {
		return nil, fmt.Errorf("No commands read from source. Was the correct dialect used?")
	}
	return engine, nil
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
