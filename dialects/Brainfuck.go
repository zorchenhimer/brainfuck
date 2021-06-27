package dialects

type Brainfuck RuneMap
var bf RuneMap = RuneMap{
	'>': '>',
	'<': '<',
	'+': '+',
	'-': '-',
	'.': '.',
	',': ',',
	'[': '[',
	']': ']',
	//'#': '#',
	//'!': '!',
}

func init() {
	registerDialect("Brainfuck", bf)
}

func (b Brainfuck) Type() LangType { return RuneLang }
