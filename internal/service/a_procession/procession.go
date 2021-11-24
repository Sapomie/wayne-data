package a_procession

import (
	"context"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/service/c_book"
	"github.com/Sapomie/wayne-data/internal/service/c_series"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type ProcessionService struct {
	ctx           context.Context
	cache         *model.Cache
	eventDb       *model.EventModel
	taskDb        *model.TaskModel
	parentDb      *model.ParentModel
	stuffDb       *model.StuffModel
	tagDb         *model.TagModel
	projectDb     *model.ProjectModel
	bookService   c_book.BookService
	seriesService c_series.SeriesService
}

func NewProcessionService(c context.Context, db *gorm.DB, cache *redis.Pool) ProcessionService {
	return ProcessionService{
		ctx:           c,
		cache:         model.NewCache(cache),
		eventDb:       model.NewEventModel(db),
		taskDb:        model.NewTaskModel(db),
		parentDb:      model.NewParentModel(db),
		stuffDb:       model.NewStuffModel(db),
		tagDb:         model.NewTagModel(db),
		projectDb:     model.NewProjectModel(db),
		bookService:   c_book.NewBookService(c, db, cache),
		seriesService: c_series.NewSeriesService(c, db, cache),
	}
}

func (svc ProcessionService) ProcessAll() (info gin.H, err error) {
	bookInfo, err := svc.bookService.ProcessBook()
	if err != nil {
		return nil, err
	}
	seriesInfo, err := svc.seriesService.ProcessSeries()
	if err != nil {
		return nil, err
	}
	err = svc.UpdateFieldVariables()
	if err != nil {
		return nil, err
	}

	return gin.H{
		"book":   bookInfo,
		"series": seriesInfo,
	}, nil
}

func (svc ProcessionService) UpdateFieldVariables() error {
	err := svc.parentDb.UpdateParentVariables()
	if err != nil {
		return err
	}
	err = svc.projectDb.UpdateProjectVariables()
	if err != nil {
		return err
	}
	err = svc.stuffDb.UpdateStuffVariables()
	if err != nil {
		return err
	}
	err = svc.tagDb.UpdateTagVariables()
	if err != nil {
		return err
	}
	err = svc.taskDb.UpdateTaskVariables()
	if err != nil {
		return err
	}
	err = svc.eventDb.UpdateNewest()
	if err != nil {
		return err
	}
	return nil
}
