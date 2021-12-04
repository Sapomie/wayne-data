package b_event

import (
	"context"
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/model/resp"
	"github.com/Sapomie/wayne-data/pkg/convert"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"github.com/jinzhu/gorm"
	"time"
)

type EvtFieldService struct {
	ctx context.Context
	db  *gorm.DB
}

func NewEvtFieldService(c context.Context) EvtFieldService {
	return EvtFieldService{
		ctx: c,
		db:  global.DBEngine,
	}
}

func (svc EvtFieldService) GetFieldList(typ int) (response []*resp.EventFieldResponse, err error) {
	var eventFields model.EventFields
	switch typ {
	case model.TypeParent:
		eventFields, err = model.NewParentModel(svc.db).GetAll()
	case model.TypeTask:
		eventFields, err = model.NewTaskModel(svc.db).GetAll()
	case model.TypeStuff:
		eventFields, err = model.NewStuffModel(svc.db).GetAll()
	case model.TypeTag:
		eventFields, err = model.NewTagModel(svc.db).GetAll()
	case model.TypeProject:
		eventFields, err = model.NewProjectModel(svc.db).GetAll()
	}
	if err != nil {
		return nil, err
	}

	for _, field := range eventFields.ToEventFields() {
		first, last := field.FieldFirstTimeAndLastTime()
		fromNow := time.Now().Sub(time.Unix(last, 0)).Hours() / 24
		longest := float64(field.FieldLongest()) / 3600 / 24
		if fromNow > longest {
			longest = fromNow
		}
		stuffResp := &resp.EventFieldResponse{
			Name:      field.FieldName(),
			Duration:  field.FieldTotalDuration(),
			Times:     field.FieldEventNum(),
			FirstTime: time.Unix(first, 0).Format(mtime.TimeTemplate5),
			LastTime:  time.Unix(last, 0).Format(mtime.TimeTemplate5),
			FromNow:   convert.FloatTo(fromNow).Decimal(1),
			Longest:   convert.FloatTo(longest).Decimal(1),
		}
		response = append(response, stuffResp)
	}
	return response, nil
}
