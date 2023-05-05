package dialects

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
