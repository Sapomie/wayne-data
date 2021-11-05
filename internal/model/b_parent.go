package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Parent struct {
	Id            int     `gorm:"primary_key"`
	Name          string  `gorm:"not null;unique"`
	WeekGoal      float64 `gorm:"not null"`
	EventNum      int64   `gorm:"not null" json:"tag_num"`
	TotalDuration float64 `gorm:"not null" json:"total_duration"`
	FirstTime     int64   `gorm:"not null"`
	LastTime      int64   `gorm:"not null"`
	Longest       int64   `gorm:"not null"`

	*Model
}

func (p *Parent) TableName() string {
	return "b_parent"
}

type Parents []*Parent

type ParentModel struct {
	Base *BaseModel
}

func NewParentModel(db *gorm.DB) *ParentModel {
	return &ParentModel{NewBaseModel(new(Parent), db)}
}

func (em *ParentModel) Exists(name string) (bool, error) {
	var count int
	err := em.Base.Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	exists := count > 0
	return exists, nil
}

func (em *ParentModel) GetAll() (Parents, error) {
	var parents Parents
	err := em.Base.Order("last_time desc").Scan(&parents).Error
	if err != nil {
		return nil, err
	}
	return parents, nil
}

func (em *ParentModel) ByName(name string) (*Parent, error) {
	parent := new(Parent)
	err := em.Base.Where("name = ?", name).Scan(parent).Error
	if err != nil {
		return nil, err
	}
	return parent, nil
}

func (em *ParentModel) ListParents(limit, offset int) (Parents, int, error) {
	var (
		parents Parents
		count   int
	)
	db := em.Base.DB
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(limit).Offset(offset).Scan(&parents).Error
	if err != nil {
		return nil, 0, err
	}
	return parents, count, nil
}

func (em *ParentModel) InsertAndGetParent(name string) (parent *Parent, info string, err error) {
	exists, err := em.Exists(name)
	if err != nil {
		return nil, "", err
	}
	if !exists {
		err = em.Base.Create(&Parent{Name: name}).Error
		if err != nil {
			return nil, "", err
		}
		info = fmt.Sprintf("Add Parent %v ", name)
	}
	parent, err = em.ByName(name)
	if err != nil {
		return nil, "", err
	}
	return
}
