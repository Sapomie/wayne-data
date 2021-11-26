package a_procession

import (
	"context"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/service/c_book"
	"github.com/Sapomie/wayne-data/internal/service/c_project"
	"github.com/Sapomie/wayne-data/internal/service/c_series"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type ProcessionService struct {
	ctx            context.Context
	cache          *model.Cache
	db             *gorm.DB
	bookService    c_book.BookService
	seriesService  c_series.SeriesService
	projectService c_project.ProjectService
}

func NewProcessionService(c context.Context, db *gorm.DB, cache *redis.Pool) ProcessionService {
	return ProcessionService{
		ctx:            c,
		db:             db,
		bookService:    c_book.NewBookService(c, db, cache),
		seriesService:  c_series.NewSeriesService(c, db, cache),
		projectService: c_project.NewProjectService(c, db, cache),
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
	_, err = svc.projectService.ProcessProject()
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
	err := model.NewParentModel(svc.db).UpdateParentVariables()
	if err != nil {
		return err
	}
	err = model.NewProjectModel(svc.db).UpdateProjectVariables()
	if err != nil {
		return err
	}
	err = model.NewStuffModel(svc.db).UpdateStuffVariables()
	if err != nil {
		return err
	}
	err = model.NewTagModel(svc.db).UpdateTagVariables()
	if err != nil {
		return err
	}
	err = model.NewTaskModel(svc.db).UpdateTaskVariables()
	if err != nil {
		return err
	}
	err = model.NewEventModel(svc.db).UpdateNewest()
	if err != nil {
		return err
	}
	return nil
}
