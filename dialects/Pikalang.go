package dialects

// https://www.dcode.fr/pikalang-language

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
