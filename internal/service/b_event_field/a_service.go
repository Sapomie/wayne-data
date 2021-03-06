package b_event

import (
	"context"
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"github.com/Sapomie/wayne-data/internal/model/resp"
	"github.com/Sapomie/wayne-data/pkg/convert"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"github.com/jinzhu/gorm"
	"time"
)

type EvtFieldService struct {
	ctx   context.Context
	db    *gorm.DB
	cache *model.Cache
}

func NewEvtFieldService(c context.Context) EvtFieldService {
	return EvtFieldService{
		ctx:   c,
		db:    global.DBEngine,
		cache: model.NewCache(global.CacheEngine),
	}
}

func (svc EvtFieldService) GetFieldList(typ int) (response []*resp.EventFieldResponse, err error) {
	resp := make([]*resp.EventFieldResponse, 0)
	key := getEventFieldKey(typ)
	exists, err := svc.cache.Get(key, &resp)
	if err != nil {
		return nil, err
	}
	if !exists {
		resp, err = svc.GetFieldListFromDB(typ)
		if err != nil {
			return nil, err
		}
		err = svc.cache.Set(key, resp, 0)
		if err != nil {
			return nil, err
		}
	}
	//redis 存储的From now勘误
	for _, r := range resp {
		r.FromNow = convert.FloatTo(time.Now().Sub(time.Unix(r.LastTimeT, 0)).Hours() / 24).Decimal(1)
	}
	return resp, nil
}

func (svc EvtFieldService) GetFieldListFromDB(typ int) (response []*resp.EventFieldResponse, err error) {
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
		fieldResp := toFieldResponse(field)
		response = append(response, fieldResp)
	}
	return response, nil
}

func toFieldResponse(field model.EventField) *resp.EventFieldResponse {
	first, last := field.FieldFirstTimeAndLastTime()
	fromNow := time.Now().Sub(time.Unix(last, 0)).Hours() / 24
	longest := float64(field.FieldLongest()) / 3600 / 24
	if fromNow > longest {
		longest = fromNow
	}
	fieldResp := &resp.EventFieldResponse{
		Name:      field.FieldName(),
		Duration:  field.FieldTotalDuration(),
		Times:     field.FieldEventNum(),
		FirstTime: time.Unix(first, 0).Format(mtime.TimeTemplate5),
		LastTime:  time.Unix(last, 0).Format(mtime.TimeTemplate5),
		FromNow:   convert.FloatTo(fromNow).Decimal(1),
		Longest:   convert.FloatTo(longest).Decimal(1),
		LastTimeT: last,
	}
	return fieldResp
}

func getEventFieldKey(typ int) string {
	switch typ {
	case model.TypeParent:
		return cons.RedisKeyParent
	case model.TypeTask:
		return cons.RedisKeyTask
	case model.TypeStuff:
		return cons.RedisKeyStuff
	case model.TypeTag:
		return cons.RedisKeyTag
	case model.TypeProject:
		return cons.RedisKeyProject
	}
	return ""
}
