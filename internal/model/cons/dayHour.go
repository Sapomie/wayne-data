package cons

const (
	DayHourOther = 1 + iota
	DayHourDaily
	DayHourSelfEntertain
	DayHourSleep
	DayHourBlank
	DayHourRoutine
)

const (
	OtherFull = 10.0
	SelfFull  = 15.0
)

const (
	DHOther         = "Other"
	DHDaily         = "Daily"
	DHSelfEntertain = "Self"
	DHSleep         = "Sleep"
	DHBlank         = "Blank"
	DHRoutine       = "Routine"
)

var DayHourNames = []string{
	DHOther,
	DHDaily,
	DHSelfEntertain,
	DHSleep,
	DHBlank,
	DHRoutine,
}

var (
	DailyFull float64
)
