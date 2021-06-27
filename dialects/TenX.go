package dialects

type TenX RuneMap
var tx RuneMap = RuneMap{
	'\u2715': '>',
	'\u00D7': '<',
	'\u0058': '+',
	'\u0087': '-',
	'\u2716': '.',
	'\U0001D4CD': ',',
	'\u2717': '[',
	'\u2613': ']',
	'\u24CD': '#',
	'\u2612': '!',
}

func init() { registerDialect("TenX", tx) }
func (t TenX) Type() LangType { return RuneLang }
