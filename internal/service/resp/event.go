package resp

type EventListRequest struct {
	TaskId   int `form:"tag_id" binding:"gte=-2"`
	ParentId int `form:"state,default=1" binding:"gte=-2"`
}

type EventResponse struct {
	Date    string
	Comment string
}

type EventsResponse []*EventResponse
