package constant

const (
	DayHourOther = 1 + iota
	DayHourDaily
	DayHourSelfEntertain
	DayHourSleep
	DayHourBlank
	DayHourRoutine
)

const (
	OtherFull     = 11.0
	SelfFull      = 17.0
	CountGoalBase = 0.75
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
