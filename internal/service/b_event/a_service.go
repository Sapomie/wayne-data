package b_event

import (
	"context"
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/model/resp"
	"github.com/Sapomie/wayne-data/pkg/app"
	"github.com/jinzhu/gorm"
)

type EventService struct {
	ctx context.Context
	db  *gorm.DB
}

func NewEventService(c context.Context) EventService {
	return EventService{
		ctx: c,
		db:  global.DBEngine,
	}
}

func (svc EventService) GetEventList(param *resp.EventListRequest, limit, offset int) ([]*resp.EventResponse, int, error) {
	dbReq, err := makeDbRequestParam(param)
	if err != nil {
		return nil, 0, err
	}

	events, num, err := model.NewEventModel(svc.db).ListEvents(dbReq, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	eventsResp, err := svc.makeEventListResponse(events)
	if err != nil {
		return nil, 0, err
	}
	return eventsResp, num, nil
}

func (svc EventService) makeEventListResponse(events model.Events) (resp.EventsResponse, error) {

	var eventsResp resp.EventsResponse
	for _, event := range events {
		var taskName, projectName, stuffName, tagName string
		task, err := model.NewTaskModel(svc.db).ById(event.TaskId)
		if err == nil {
			taskName = task.Name
		} else if err != gorm.ErrRecordNotFound {
			return nil, err
		}
		project, err := model.NewProjectModel(svc.db).ById(event.ProjectId)
		if err == nil {
			projectName = project.Name
		} else if err != gorm.ErrRecordNotFound {
			return nil, err
		}

		eventResp := &resp.EventResponse{
			Date:     event.Date,
			Task:     taskName,
			Comment:  event.Comment,
			Duration: event.Duration,
			Project:  projectName,
			Stuff:    stuffName,
			Tag:      tagName,
		}
		eventsResp = append(eventsResp, eventResp)
	}

	return eventsResp, nil
}

func makeDbRequestParam(p *resp.EventListRequest) (*resp.DbEventListRequest, error) {

	start, end, err := app.DateStartAndEnd(p.Date, p.Span)
	if err != nil {
		return nil, err
	}

	return &resp.DbEventListRequest{
		TaskId:    model.TaskInfoByName[p.Task].Id,
		ParentId:  model.ParentInfoByName[p.Parent].Id,
		ProjectId: model.ProjectInfoByName[p.Project].Id,
		//StuffId:   "",
		//TagId:     "",
		Start: start,
		End:   end,
	}, nil
}
