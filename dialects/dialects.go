package dialects

import (
	"fmt"
)

type Dialect interface {
	Type() LangType
}

type LangType int

const (
	RuneLang LangType = iota
	WordLang
)

type RuneMap map[rune]rune
type WordMap map[string]rune

func (r RuneMap) Type() LangType { return RuneLang }
func (r WordMap) Type() LangType { return WordLang }

var Dialects map[string]Dialect

func registerDialect(name string, dialect Dialect) {
	if Dialects == nil {
		Dialects = map[string]Dialect{}
	}
	if _, exists := Dialects[name]; exists {
		panic(fmt.Sprintf("Dialect with name %q already exists!", name))
	}

	Dialects[name] = dialect
}
