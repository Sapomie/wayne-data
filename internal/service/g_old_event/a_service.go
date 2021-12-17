package g_old_event

import (
	"context"
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/model/resp"
	"github.com/Sapomie/wayne-data/pkg/app"
	"github.com/jinzhu/gorm"
)

type OldEventService struct {
	ctx context.Context
	db  *gorm.DB
}

func NewOldEventService(c context.Context) OldEventService {
	return OldEventService{
		ctx: c,
		db:  global.DBEngine,
	}
}

func (svc OldEventService) GetOldEventList(param *resp.OldEventListRequest, limit, offset int) ([]*resp.OldEventResponse, int, error) {
	dbReq, err := makeDbRequestParam(param)
	if err != nil {
		return nil, 0, err
	}
	oldEvents, num, err := model.NewOldEventModel(svc.db).ListOldEvents(dbReq, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	oldEventsResp, err := svc.makeOldEventListResponse(oldEvents)
	if err != nil {
		return nil, 0, err
	}
	return oldEventsResp, num, nil
}

func (svc OldEventService) makeOldEventListResponse(oldEvents model.OldEvents) ([]*resp.OldEventResponse, error) {

	oldEventResponses := make([]*resp.OldEventResponse, 0)
	for _, e := range oldEvents {
		oldEventResp := &resp.OldEventResponse{
			Date:     e.Date,
			Task:     e.TaskName,
			Comment:  e.Comment,
			Duration: e.Duration,
			Parent:   e.ParentTask,
		}
		oldEventResponses = append(oldEventResponses, oldEventResp)
	}

	return oldEventResponses, nil
}

func makeDbRequestParam(p *resp.OldEventListRequest) (*resp.DbOldEventListRequest, error) {

	start, end, err := app.DateStartAndEnd(p.Date, p.Span)
	if err != nil {
		return nil, err
	}

	return &resp.DbOldEventListRequest{
		Task:   p.Task,
		Parent: p.Parent,
		Word:   p.Word,
		Start:  start,
		End:    end,
	}, nil
}
