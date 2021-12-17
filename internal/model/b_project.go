package model

import (
	"fmt"
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"github.com/jinzhu/gorm"
)

const (
	ProjectFinish = 1
)

type Project struct {
	Id            int    `gorm:"primary_key"`
	Name          string `gorm:"not null"`
	TaskId        int    `gorm:"not null"`
	TagId         string
	ViaId         int
	Finish        int
	EventNum      int64   `gorm:"not null" json:"tag_num"`
	TotalDuration float64 `gorm:"not null" json:"total_duration"`
	FirstTime     int64   `gorm:"not null"`
	LastTime      int64   `gorm:"not null"`
	Longest       int64   `gorm:"not null"`

	*Base
}

func (p *Project) TableName() string {
	return "b_project"
}

func (p *Project) FieldName() string {
	return p.Name
}

func (p *Project) FieldTotalDuration() float64 {
	return p.TotalDuration
}

func (p *Project) FieldEventNum() int64 {
	return p.EventNum
}

func (p *Project) FieldFirstTimeAndLastTime() (int64, int64) {
	return p.FirstTime, p.LastTime
}

func (p *Project) FieldLongest() int64 {
	return p.Longest
}

type Projects []*Project

func (ps Projects) ToEventFields() []EventField {
	eventFields := make([]EventField, 0)
	for _, project := range ps {
		var ef EventField
		ef = project
		eventFields = append(eventFields, ef)
	}
	return eventFields
}

type ProjectModel struct {
	Base *BaseDbModel
}

func NewProjectModel(db *gorm.DB) *ProjectModel {
	return &ProjectModel{NewBaseModel(new(Project), db)}
}

func (em *ProjectModel) Exists(name string) (bool, error) {
	var count int
	err := em.Base.Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	exists := count > 0
	return exists, nil
}

func (em *ProjectModel) GetAll() (Projects, error) {
	var projects Projects
	err := em.Base.Order("last_time desc").Scan(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (em *ProjectModel) ByName(name string) (*Project, error) {
	project := new(Project)
	err := em.Base.Where("name = ?", name).Scan(project).Error
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (em *ProjectModel) ById(id int) (*Project, error) {
	project := new(Project)
	err := em.Base.Where("id = ?", id).Scan(project).Error
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (em *ProjectModel) ListProjects(limit, offset int) (Projects, int, error) {
	var (
		projects Projects
		count    int
	)
	db := em.Base.DB
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(limit).Offset(offset).Scan(&projects).Error
	if err != nil {
		return nil, 0, err
	}
	return projects, count, nil
}

func InsertAndGetProject(db *gorm.DB, project *Project) (projectDb *Project, info string, err error) {
	em := NewProjectModel(db)
	exists, err := em.Exists(project.Name)
	if err != nil {
		return nil, "", err
	}
	if !exists {
		err = em.Base.Create(project).Error
		if err != nil {
			return nil, "", err
		}
		info = fmt.Sprintf("Add Project %v ", project.Name)
	} else {
		err = em.Base.Where("name = ?", project.Name).Update(project).Error
		if err != nil {
			return nil, "", err
		}
	}
	projectDb, err = NewProjectModel(db).ByName(project.Name)
	if err != nil {
		return nil, "", err
	}
	return
}

var ProjectInfoById = make(map[int]struct {
	Name string
})

var ProjectInfoByName = make(map[string]struct {
	Id int
})

func (em *ProjectModel) UpdateProjectVariables() error {
	projects, err := em.GetAll()
	if err != nil {
		return err
	}
	for _, project := range projects {
		ProjectInfoById[project.Id] = struct {
			Name string
		}{Name: project.Name}

		ProjectInfoByName[project.Name] = struct {
			Id int
		}{Id: project.Id}
	}

	return nil
}

func UpdateProjectColumn(db *gorm.DB) (err error) {
	if err = NewProjectModel(db).UpdateProjectVariables(); err != nil {
		return err
	}

	projects, err := NewProjectModel(db).GetAll()
	if err != nil {
		return nil
	}
	for _, project := range projects {
		var data struct {
			Num            int64
			Dur            float64
			FirstTimestamp int64
			LastTimestamp  int64
		}
		err := NewEventModel(db).Base.Select("sum(duration) as dur, count(id) as num").Where("project_id = ?", project.Id).Scan(&data).Error
		if err != nil {
			return err
		}
		err = NewEventModel(db).Base.Select("start_time as first_timestamp").Where("project_id = ?", project.Id).Order("start_time asc").Limit(1).Scan(&data).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		err = NewEventModel(db).Base.Select("start_time as last_timestamp").Where("project_id = ?", project.Id).Order("start_time desc").Limit(1).Scan(&data).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		evts, err := NewEventModel(db).ByProjectName(cons.DbOldest, cons.DbNewest, project.Name)
		if err != nil {
			return err
		}
		longest := int64(0)
		if len(evts) > 0 {
			former := evts[0].StartTime
			for _, evt := range evts {
				span := evt.StartTime - former
				if span > longest {
					longest = span
				}
				former = evt.StartTime
			}
		}
		err = NewProjectModel(db).Base.Where("id = ?", project.Id).Update(map[string]interface{}{
			"event_num":      data.Num,
			"total_duration": data.Dur,
			"first_time":     data.FirstTimestamp,
			"last_time":      data.LastTimestamp,
			"longest":        longest,
		}).Error
		if err != nil {
			return err
		}
	}

	return nil
}
