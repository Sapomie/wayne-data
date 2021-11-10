package resp

type EssentialDayListRequest struct {
	Limit int `form:"limit,default=10" binding:"gte=0"`
}
