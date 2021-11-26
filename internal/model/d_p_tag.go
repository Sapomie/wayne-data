package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

const (
	TypeProjectTag = iota + 1
	TypeProjectVia
)

type PTag struct {
	Id            int     `gorm:"primary_key"`
	Name          string  `gorm:"not null;unique"`
	Typ           int     `gorm:"not null"`
	EventNum      int64   `gorm:"not null"`
	TotalDuration float64 `gorm:"not null" json:"total_duration"`
	FirstTime     int64   `gorm:"not null"`
	LastTime      int64   `gorm:"not null"`
	Longest       int64   `gorm:"not null"`

	*Base
}

func (e *PTag) TableName() string {
	return "d_project_tag"
}

type PTags []*PTag

type PTagModel struct {
	Base *BaseDbModel
}

func NewPTagModel(db *gorm.DB) *PTagModel {
	return &PTagModel{NewBaseModel(new(PTag), db)}
}

func (em *PTagModel) Exists(name string) (bool, error) {
	var count int
	err := em.Base.Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	exists := count > 0
	return exists, nil
}

func (em *PTagModel) GetAll() (PTags, error) {
	var pTags PTags
	err := em.Base.Order("last_time desc").Scan(&pTags).Error
	if err != nil {
		return nil, err
	}
	return pTags, nil
}

func (em *PTagModel) ByName(name string) (*PTag, error) {
	pTag := new(PTag)
	err := em.Base.Where("name = ?", name).Scan(pTag).Error
	if err != nil {
		return nil, err
	}
	return pTag, nil
}

func (em *PTagModel) ListPTags(limit, offset int) (PTags, int, error) {
	var (
		pTags PTags
		count int
	)
	db := em.Base.DB
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(limit).Offset(offset).Scan(&pTags).Error
	if err != nil {
		return nil, 0, err
	}
	return pTags, count, nil
}

func (em *PTagModel) InsertAndGetPTag(name string, typ int) (pTag *PTag, info string, err error) {
	exists, err := em.Exists(name)
	if err != nil {
		return nil, "", err
	}
	if !exists {
		err = em.Base.Create(&PTag{Name: name, Typ: typ}).Error
		if err != nil {
			return nil, "", err
		}
		info = fmt.Sprintf("Add PTag %v ", name)
	}
	pTag, err = em.ByName(name)
	if err != nil {
		return nil, "", err
	}
	return
}
