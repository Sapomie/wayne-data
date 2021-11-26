package model

import (
	"github.com/jinzhu/gorm"
)

type Abbr struct {
	Id      int    `gorm:"primary_key"`
	TaskId  int    `gorm:"not null"`
	Abbr    string `gorm:"not null"`
	Content string `gorm:"not null"`

	*Base
}

func (e *Abbr) TableName() string {
	return "f_abbreviate"
}

type Abbrs []*Abbr

type AbbrModel struct {
	Base *BaseDbModel
}

func NewAbbrModel(db *gorm.DB) *AbbrModel {
	return &AbbrModel{NewBaseModel(new(Abbr), db)}
}

func (em *AbbrModel) GetAll() (Abbrs, error) {
	var abbrs Abbrs
	err := em.Base.Scan(&abbrs).Error
	if err != nil {
		return nil, err
	}
	return abbrs, nil
}

func (em *AbbrModel) Exists(taskId int, abbr string) (bool, error) {
	var count int
	err := em.Base.Where("task_id = ? and abbr = ?", taskId, abbr).Count(&count).Error
	if err != nil {
		return false, err
	}
	exists := count > 0
	return exists, nil
}

func InsertAbbr(db *gorm.DB, abbr *Abbr) (err error) {
	em := NewAbbrModel(db)
	exists, err := em.Exists(abbr.TaskId, abbr.Abbr)
	if err != nil {
		return err
	}
	if !exists {
		err = em.Base.Create(abbr).Error
		if err != nil {
			return err
		}
	} else {
		err = em.Base.Where("task_id = ? and abbr = ?", abbr.TaskId, abbr.Content).Update(abbr).Error
		if err != nil {
			return err
		}
	}

	return
}
