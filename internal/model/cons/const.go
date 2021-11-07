package cons

import "time"

const (
	All         = int64(0)
	BiggestTime = int64(2000000000)
	IsSum       = 1
	IsOver      = 1
	BookFinish  = 1
	BookAbandon = 2
)

const (
	Finished        = "√"
	UnFinished      = "✕"
	RestrainFull    = "☯"
	RestrainNotFull = "☺"
	None            = "-----"
	LevelMark       = "(↥)"
	LevelMinMark    = "(↧)"
)

const (
	PushUp        = "push-up"
	SitUp         = "sit-up"
	DumbbellPress = "dumbbellPress"
)

var Newest time.Time