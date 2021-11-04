package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Stuff struct {
	Id            int     `gorm:"primary_key"`
	Name          string  `gorm:"not null;unique"`
	EventNum      int64   `gorm:"not null" json:"tag_num"`
	TotalDuration float64 `gorm:"not null" json:"total_duration"`
	FirstTime     int64   `gorm:"not null"`
	LastTime      int64   `gorm:"not null"`
	Longest       int64   `gorm:"not null"`

	*Model
}

func (e *Stuff) TableName() string {
	return "b_stuff"
}

type Stuffs []*Stuff

type StuffModel struct {
	Base *BaseModel
}

func NewStuffModel(db *gorm.DB) *StuffModel {
	return &StuffModel{NewBaseModel(new(Stuff), db)}
}

func (em *StuffModel) Exists(name string) (bool, error) {
	var count int
	err := em.Base.Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	exists := count > 0
	return exists, nil
}

func (em *StuffModel) GetAll() (Stuffs, error) {
	var stuffs Stuffs
	err := em.Base.Order("last_time desc").Scan(&stuffs).Error
	if err != nil {
		return nil, err
	}
	return stuffs, nil
}

func (em *StuffModel) ByName(name string) (*Stuff, error) {
	stuff := new(Stuff)
	err := em.Base.Where("name = ?", name).Scan(stuff).Error
	if err != nil {
		return nil, err
	}
	return stuff, nil
}

func (em *StuffModel) ListStuffs(limit, offset int) (Stuffs, int, error) {
	var (
		stuffs Stuffs
		count  int
	)
	db := em.Base.DB
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(limit).Offset(offset).Scan(&stuffs).Error
	if err != nil {
		return nil, 0, err
	}
	return stuffs, count, nil
}

func (em *StuffModel) InsertAndGetStuff(name string) (stuff *Stuff, info string, err error) {
	exists, err := em.Exists(name)
	if err != nil {
		return nil, "", err
	}
	if !exists {
		err = em.Base.Create(&Stuff{Name: name}).Error
		if err != nil {
			return nil, "", err
		}
		info = fmt.Sprintf("Add Stuff %v ", name)
	}
	stuff, err = em.ByName(name)
	if err != nil {
		return nil, "", err
	}
	return
}
