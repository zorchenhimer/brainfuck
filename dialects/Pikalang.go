package dialects

// https://www.dcode.fr/pikalang-language

type Pikalang WordMap
var pl WordMap = WordMap{
	"pipi": '>',
	"pichu": '<',
	"pi": '+',
	"ka": '-',
	"pikachu": '.',
	"pikapi": ',',
	"pika": '[',
	"chu": ']',
}

func init() { registerDialect("Pikalang", pl) }
func (p Pikalang) Type() LangType { return WordLang }
