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
