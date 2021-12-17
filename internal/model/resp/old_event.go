package resp

import "time"

type OldEventListRequest struct {
	Task   string `form:"task"`
	Parent string `form:"parent"`
	Word   string `form:"word" binding:"min=0,max=20"`
	Date   string `form:"date" binding:"omitempty,min=8,max=8"`
	Span   int    `form:"span"`
}

type OldEventResponse struct {
	Date     string
	Task     string
	Parent   string
	Comment  string
	Duration float64
}

type DbOldEventListRequest struct {
	Task   string
	Parent string
	Word   string
	Start  time.Time
	End    time.Time
}
