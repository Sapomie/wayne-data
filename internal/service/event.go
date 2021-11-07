package service

import (
	"context"
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/resp"
)

type EventService struct {
	ctx context.Context
	*model.EventModel
}

func NewEventService(c context.Context) EventService {
	return EventService{
		ctx:        c,
		EventModel: model.NewEventModel(global.DBEngine),
	}
}

func (svc *EventService) GetEventList(param *resp.EventListRequest, limit, offset int) ([]*resp.EventResponse, int, error) {
	events, num, err := svc.ListEvents(param.ParentId, param.TaskId, limit, offset)
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