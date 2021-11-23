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
