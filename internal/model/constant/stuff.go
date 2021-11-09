package constant

//main stuff
var MainStuffs = []string{
	StuGc,
	StuProtein1,
	StuMovie,
	StuKeganmin,
	StuWine,
	StuMask,
}

//stuff name
const (
	StuGc       = "gc"
	StuMovie    = "mv"
	StuWine     = "wine"
	StuKeganmin = "keganmin"
	StuProtein1 = "protein1"
	StuMask     = "mask"
	WinningDays = "winningDays"
)

var RestrainStuff = []string{
	StuGc,
	StuKeganmin,
	StuWine,
}

func IsRestrain(stuffName string) bool {
	for _, v := range RestrainStuff {
		if v == stuffName {
			return true
		}
	}
	return false
}
