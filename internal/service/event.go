package service

import (
	"github.com/Sapomie/wayne-data/pkg/app"
)

type EventListRequest struct {
	TaskId   int `form:"tag_id" binding:"gte=1"`
	ParentId int `form:"state,default=1" binding:"oneof=0 1"`
}

type EventResponse struct {
}

func (svc *Service) GetEventList(param *EventListRequest, pager *app.Pager) ([]*EventResponse, int, error) {

	return nil, 0, nil
}
