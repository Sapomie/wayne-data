package a_procession

import (
	"context"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/service/b_project"
	"github.com/Sapomie/wayne-data/internal/service/c_book"
	"github.com/Sapomie/wayne-data/internal/service/c_run"
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
	projectService b_project.ProjectService
	runService     c_run.RunService
}

func NewProcessionService(c context.Context, db *gorm.DB, cache *redis.Pool) ProcessionService {
	return ProcessionService{
		ctx:            c,
		db:             db,
		bookService:    c_book.NewBookService(c, db, cache),
		seriesService:  c_series.NewSeriesService(c, db, cache),
		projectService: b_project.NewProjectService(c, db, cache),
		runService:     c_run.NewRunService(c, db, cache),
	}
}

func (svc ProcessionService) ProcessAll() (info gin.H, err error) {
	err = svc.UpdateFieldVariables()
	if err != nil {
		return nil, err
	}

	bookInfo, err := svc.bookService.ProcessBook()
	if err != nil {
		return nil, err
	}
	seriesInfo, err := svc.seriesService.ProcessSeries()
	if err != nil {
		return nil, err
	}
	runInfo, err := svc.runService.ProcessRun()
	if err != nil {
		return nil, err
	}
	_, err = svc.projectService.ProcessProject()
	if err != nil {
		return nil, err
	}

	return gin.H{
		"book":   bookInfo,
		"series": seriesInfo,
		"runs":   runInfo,
	}, nil
}

func (svc ProcessionService) UpdateFieldVariables() error {
	err := model.NewEventModel(svc.db).UpdateNewestAndOldest()
	if err != nil {
		return err
	}
	err = model.UpdateParentColumn(svc.db)
	if err != nil {
		return err
	}
	err = model.UpdateTaskColumn(svc.db)
	if err != nil {
		return err
	}
	err = model.UpdateProjectColumn(svc.db)
	if err != nil {
		return err
	}
	err = model.UpdateTagColumn(svc.db)
	if err != nil {
		return err
	}
	err = model.UpdateStuffColumn(svc.db)
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

	return nil
}
