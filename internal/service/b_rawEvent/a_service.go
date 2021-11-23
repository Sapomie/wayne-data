package b_rawEvent

import (
	"context"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
)

type RawEventService struct {
	ctx       context.Context
	cache     *model.Cache
	eventDb   *model.EventDbModel
	taskDb    *model.TaskModel
	parentDb  *model.ParentModel
	stuffDb   *model.StuffModel
	tagDb     *model.TagModel
	projectDb *model.ProjectModel
	abbrDb    *model.AbbrModel
}

func NewRawEventService(c context.Context, db *gorm.DB, cache *redis.Pool) RawEventService {
	return RawEventService{
		ctx:       c,
		cache:     model.NewCache(cache),
		eventDb:   model.NewEventModel(db),
		taskDb:    model.NewTaskModel(db),
		parentDb:  model.NewParentModel(db),
		stuffDb:   model.NewStuffModel(db),
		tagDb:     model.NewTagModel(db),
		projectDb: model.NewProjectModel(db),
	}
}

func (svc RawEventService) ImportCsvData() (model.Events, map[string]interface{}, error) {
	rawEvents, _, err := getRawEventFromCsvFile()
	if err != nil {
		return nil, nil, err
	}

	events, taskAndParentAddingInfo, err := svc.makeEvents(rawEvents)
	if err != nil {
		return nil, nil, err
	}

	eventsStoreInfo, err := svc.storeEvents(events)
	if err != nil {
		return nil, nil, err
	}

	err = svc.cache.FlushDb()
	if err != nil {
		return nil, nil, err
	}

	info := make(map[string]interface{})
	info["Task and Parent"] = taskAndParentAddingInfo
	info["Events"] = eventsStoreInfo

	return events, info, nil
}
