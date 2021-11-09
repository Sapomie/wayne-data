package resp

type EventListRequest struct {
	TaskId   int `form:"tag_id,default=-1" binding:"gte=-2"`
	ParentId int `form:"parent_id,default=-1" binding:"gte=-2"`
}

type EventResponse struct {
	Date    string
	Comment string
}

type EventsResponse []*EventResponse
