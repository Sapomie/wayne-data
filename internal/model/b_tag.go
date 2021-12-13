package model

import (
	"fmt"
	"github.com/Sapomie/wayne-data/internal/model/cons"
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

func (t *Tag) TableName() string {
	return "b_tag"
}

func (t *Tag) FieldName() string {
	return t.Name
}

func (t *Tag) FieldTotalDuration() float64 {
	return t.TotalDuration
}

func (t *Tag) FieldEventNum() int64 {
	return t.EventNum
}

func (t *Tag) FieldFirstTimeAndLastTime() (int64, int64) {
	return t.FirstTime, t.LastTime
}

func (t *Tag) FieldLongest() int64 {
	return t.Longest
}

type Tags []*Tag

func (ts Tags) ToEventFields() []EventField {
	eventFields := make([]EventField, 0)
	for _, tag := range ts {
		var ef EventField
		ef = tag
		eventFields = append(eventFields, ef)
	}
	return eventFields
}

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

func (em *TagModel) ById(id int) (*Tag, error) {
	tag := new(Tag)
	err := em.Base.Where("id = ?", id).Scan(tag).Error
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

func InsertAndGetTag(db *gorm.DB, name string) (tag *Tag, info string, err error) {
	em := NewTagModel(db)
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

	tag, err = NewTagModel(db).ByName(name)
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

func (em *TagModel) UpdateTagVariables() error {

	tags, err := em.GetAll()
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

func UpdateTagColumn(db *gorm.DB) (err error) {
	if err = NewTagModel(db).UpdateTagVariables(); err != nil {
		return err
	}

	tags, err := NewTagModel(db).GetAll()
	if err != nil {
		return nil
	}
	for _, tag := range tags {
		var data struct {
			Num            int64
			Dur            float64
			FirstTimestamp int64
			LastTimestamp  int64
		}
		err := NewEventModel(db).Base.Select("sum(duration) as dur, count(id) as num").Where("tag_id = ?", tag.Id).Scan(&data).Error
		if err != nil {
			return err
		}
		err = NewEventModel(db).Base.Select("start_time as first_timestamp").Where("tag_id = ?", tag.Id).Order("start_time asc").Limit(1).Scan(&data).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		err = NewEventModel(db).Base.Select("start_time as last_timestamp").Where("tag_id = ?", tag.Id).Order("start_time desc").Limit(1).Scan(&data).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		evts, err := NewEventModel(db).ByTagName(cons.Oldest, cons.Newest, tag.Name)
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
		err = NewTagModel(db).Base.Where("id = ?", tag.Id).Update(map[string]interface{}{
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
