package dialects

// https://github.com/MiffOttah/fuckfuck

type FuckFuck WordMap
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
func (f FuckFuck) Type() LangType { return WordLang }
