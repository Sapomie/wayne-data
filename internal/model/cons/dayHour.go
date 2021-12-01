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
	OtherDailyCoefficient = 2.0
	SelfDailyCoefficient  = 3.0
	CountGoalBase         = 0.7
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
	GoalBase  float64
)
