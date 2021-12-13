package resp

import "time"

type EventListRequest struct {
	Task    string `form:"task"`
	Parent  string `form:"parent"`
	Project string `form:"project"`
	Stuff   string `form:"stuff"`
	Tag     string `form:"tag"`
	Date    string `form:"date" binding:"omitempty,min=8,max=8"`
	Span    int    `form:"span"`
}

type EventResponse struct {
	Date     string
	Task     string
	Comment  string
	Duration float64
	Stuff    string
	Project  string
	Tag      string
}

type EventsResponse []*EventResponse

type DbEventListRequest struct {
	TaskId    int
	ParentId  int
	ProjectId int
	StuffId   int
	TagId     int
	Start     time.Time
	End       time.Time
}
