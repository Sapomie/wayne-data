package b_rawEvent

import (
	"context"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"github.com/garyburd/redigo/redis"
	"github.com/gocarina/gocsv"
	"github.com/jinzhu/gorm"
	"os"
	"time"
)

type RawEventService struct {
	ctx       context.Context
	cache     *model.Cache
	eventDb   *model.EventModel
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
		abbrDb:    model.NewAbbrModel(db),
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

func (svc RawEventService) ExportAllRawEvent() error {

	events, _, err := svc.eventDb.GetAll()
	if err != nil {
		return err
	}

	raws := make([]*RawEvent, 0)
	for _, event := range events {
		raw, err := svc.eventToRawEvent(event)
		if err != nil {
			return err
		}
		raws = append(raws, raw)
	}
	endDate := time.Unix(events.Newest().StartTime, 0).Format(mtime.TimeTemplate4)
	startDate := time.Unix(events.Oldest().StartTime, 0).Format(mtime.TimeTemplate4)

	f, err := os.Create(startDate + "-----" + endDate + ".csv")
	if err != nil {
		return err
	}

	defer f.Close()

	err = gocsv.MarshalFile(raws, f)
	if err != nil {
		return err
	}
	return nil

}
