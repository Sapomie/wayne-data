package c_project

import (
	"context"
	"errors"
	"fmt"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

type ProjectService struct {
	ctx       context.Context
	cache     *model.Cache
	projectDb *model.ProjectModel
	eventDb   *model.EventModel
	pTagDb    *model.PTagModel
}

func NewProjectService(c context.Context, db *gorm.DB, cache *redis.Pool) ProjectService {
	return ProjectService{
		ctx:       c,
		cache:     model.NewCache(cache),
		projectDb: model.NewProjectModel(db),
		pTagDb:    model.NewPTagModel(db),
		eventDb:   model.NewEventModel(db),
	}
}

func (svc ProjectService) ProcessProject() ([]string, error) {
	projects, infos, err := svc.makeProjects()
	if err != nil {
		return nil, err
	}

	err = svc.updateProjects(projects)
	if err != nil {
		return nil, err
	}

	return infos, nil
}

func (svc ProjectService) makeProjects() (projects model.Projects, infos []string, err error) {

	start, end := mtime.NewTimeZone(mtime.TypeYear, 2021, 1).BeginAndEnd()
	events, err := svc.eventDb.ByTaskNames(start, end,
		cons.CodeInput,
		cons.CodeOutput,
		cons.EnglishInput,
	)
	events = events.WithProject()
	if err != nil {
		return nil, nil, err
	}

	projectMap := make(map[string]model.Events, 0)
	for _, event := range events {
		strs := strings.Split(event.Comment, "，")
		name := strs[0]
		projectMap[name] = append(projectMap[name], event)
	}

	for name, projectEvents := range projectMap {
		project := &model.Project{Name: name}
		for _, event := range projectEvents {

			if isProjectFirstTime(event) {
				var ptagIds string
				tagNames, viaName, err := projectInfo(event)
				if err != nil {
					infos = append(infos, fmt.Sprintf("make project error,event start time: %v,coment: %v", event.Start(), event.Comment))
					continue
				}
				project.FirstTime = event.StartTime
				if viaName != "" {
					via, _, err := svc.pTagDb.InsertAndGetPTag(viaName, model.TypeProjectVia)
					if err != nil {
						return nil, nil, err
					}
					project.ViaId = via.Id
				}

				for _, tagName := range tagNames {
					ptag, _, err := svc.pTagDb.InsertAndGetPTag(tagName, model.TypeProjectTag)
					if err != nil {
						return nil, nil, err
					}
					if ptagIds == "" {
						ptagIds = fmt.Sprint(ptag.Id)
					} else {
						ptagIds = ptagIds + "," + fmt.Sprint(ptag.Id)
					}
					project.TagId = ptagIds
				}
			}

			if isProjectLastTime(event) {
				project.Finish = model.ProjectFinish
			}

			if event.StartTime > project.LastTime {
				project.LastTime = event.StartTime
			}

			project.TotalDuration += event.Duration
		}
		projects = append(projects, project)
	}

	return
}

func (svc ProjectService) updateProjects(projects model.Projects) error {
	mm := svc.projectDb
	for _, project := range projects {
		exist, err := mm.Exists(project.Name)
		if err != nil {
			return err
		}
		if exist {
			err := mm.Base.Where("name = ?", project.Name).Update(project).Error
			if err != nil {
				return err
			}
		} else {
			err := mm.Base.Create(project).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isProjectFirstTime(event *model.Event) bool {
	if strings.Contains(event.Comment, "、s") {
		return true
	}
	return false
}

func isProjectLastTime(event *model.Event) bool {
	if strings.Contains(event.Comment, "、e") {
		return true
	}
	return false
}

func projectInfo(event *model.Event) (tags []string, via string, err error) {
	strs := strings.Split(event.Comment, "，")
	if len(strs) < 4 {
		return nil, "", errors.New("wrong length of project comment")
	}
	via = strs[1]
	tags = strings.Split(strs[2], "；")
	return
}

func projectRate(event *model.Event) (rate int, err error) {
	strs := strings.Split(event.Comment, "，")
	if len(strs) < 3 {
		return 0, errors.New("wrong length of project comment")
	}
	rate, err = strconv.Atoi(strs[1])
	if err != nil {
		return 0, err
	}
	return
}
