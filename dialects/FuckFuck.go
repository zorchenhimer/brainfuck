package dialects

// https://github.com/MiffOttah/fuckfuck

var ff WordMap = WordMap{
	"ass": '>',
	"bitch": '<',
	"cunt": '+',
	"damn": '-',
	"dick": '.',
	"fuck": ',',
	"shit": '[',
	"twat": ']',
}

func init() { registerDialect("FuckFuck", ff) }
