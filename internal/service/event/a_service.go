package event

import (
	"context"
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/model/resp"
)

type ServiceEvent struct {
	ctx context.Context
	db  *model.EventDbModel
}

func NewEventService(c context.Context) ServiceEvent {
	return ServiceEvent{
		ctx: c,
		db:  model.NewEventModel(global.DBEngine),
	}
}

func (svc *ServiceEvent) GetEventList(param *resp.EventListRequest, limit, offset int) ([]*resp.EventResponse, int, error) {
	events, num, err := svc.db.ListEvents(param.ParentId, param.TaskId, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	eventsResp := makeEventListResponse(events)
	return eventsResp, num, nil
}

func makeEventListResponse(events model.Events) resp.EventsResponse {

	var eventsResp resp.EventsResponse
	for _, event := range events {
		eventResp := &resp.EventResponse{
			Date:    event.Date,
			Comment: event.Comment,
		}
		eventsResp = append(eventsResp, eventResp)
	}

	return eventsResp
}
