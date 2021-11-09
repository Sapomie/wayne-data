package model

import (
	"fmt"
	"github.com/Sapomie/wayne-data/global"
	"github.com/jinzhu/gorm"
)

type Tag struct {
	Id            int     `gorm:"primary_key"`
	Name          string  `gorm:"not null;unique"`
	EventNum      int64   `gorm:"not null" json:"tag_num"`
	TotalDuration float64 `gorm:"not null" json:"total_duration"`
	FirstTime     int64   `gorm:"not null"`
	LastTime      int64   `gorm:"not null"`
	Longest       int64   `gorm:"not null"`

	*Base
}

func (e *Tag) TableName() string {
	return "b_tag"
}

type Tags []*Tag

type TagModel struct {
	Base *BaseDbModel
}

func NewTagModel(db *gorm.DB) *TagModel {
	return &TagModel{NewBaseModel(new(Tag), db)}
}

func (em *TagModel) Exists(name string) (bool, error) {
	var count int
	err := em.Base.Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	exists := count > 0
	return exists, nil
}

func (em *TagModel) GetAll() (Tags, error) {
	var tags Tags
	err := em.Base.Order("last_time desc").Scan(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (em *TagModel) ByName(name string) (*Tag, error) {
	tag := new(Tag)
	err := em.Base.Where("name = ?", name).Scan(tag).Error
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (em *TagModel) ListTags(limit, offset int) (Tags, int, error) {
	var (
		tags  Tags
		count int
	)
	db := em.Base.DB
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(limit).Offset(offset).Scan(&tags).Error
	if err != nil {
		return nil, 0, err
	}
	return tags, count, nil
}

func (em *TagModel) InsertAndGetTag(name string) (tag *Tag, info string, err error) {
	exists, err := em.Exists(name)
	if err != nil {
		return nil, "", err
	}
	if !exists {
		err = em.Base.Create(&Tag{Name: name}).Error
		if err != nil {
			return nil, "", err
		}
		info = fmt.Sprintf("Add Tag %v ", name)
	}
	tag, err = em.ByName(name)
	if err != nil {
		return nil, "", err
	}
	return
}

var TagInfoById = make(map[int]struct {
	Name string
})

var TagInfoByName = make(map[string]struct {
	Id int
})

func updateTagVariables() error {

	tags, err := NewTagModel(global.DBEngine).GetAll()
	if err != nil {
		return err
	}
	for _, tag := range tags {
		TagInfoById[tag.Id] = struct {
			Name string
		}{Name: tag.Name}

		TagInfoByName[tag.Name] = struct {
			Id int
		}{Id: tag.Id}
	}

	return nil
}
