package model

import (
	"fmt"
	"github.com/Sapomie/wayne-data/global"
	"github.com/jinzhu/gorm"
)

type Project struct {
	Id            int     `gorm:"primary_key"`
	Name          string  `gorm:"not null;unique"`
	EventNum      int64   `gorm:"not null" json:"tag_num"`
	TotalDuration float64 `gorm:"not null" json:"total_duration"`
	FirstTime     int64   `gorm:"not null"`
	LastTime      int64   `gorm:"not null"`
	Longest       int64   `gorm:"not null"`

	*Base
}

func (e *Project) TableName() string {
	return "b_project"
}

type Projects []*Project

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

func (em *ProjectModel) InsertAndGetProject(name string) (project *Project, info string, err error) {
	exists, err := em.Exists(name)
	if err != nil {
		return nil, "", err
	}
	if !exists {
		err = em.Base.Create(&Project{Name: name}).Error
		if err != nil {
			return nil, "", err
		}
		info = fmt.Sprintf("Add Project %v ", name)
	}
	project, err = em.ByName(name)
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

func updateProjectVariables() error {

	projects, err := NewProjectModel(global.DBEngine).GetAll()
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
