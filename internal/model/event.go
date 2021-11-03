package model

import "github.com/jinzhu/gorm"

type Event struct {
	//原有属性
	ID        int64   `gorm:"primary_key" json:"id"`
	Date      string  `gorm:"not null" json:"date"`
	Duration  float64 `gorm:"not null" json:"duration"`
	TaskId    int     `gorm:"not null;default:-2" json:"task_id"`
	ParentId  int     `gorm:"not null;default:-2" json:"parent_id"`
	Comment   string  `json:"comment"`
	StartTime int64   `gorm:"not null" json:"start_time"`
	EndTime   int64   `gorm:"not null" json:"end_time"`
	//自定义属性：通过comment增加

	StuffId   string `json:"stuff_id"`
	TagId     string `json:"tag_id"`
	ProjectId int    `json:"project_id"`
	Remark    string `json:"remark"` //comment 除去自定义属性的部分

	*Model
}

type Events []*Event

type EventModel struct {
	Base *BaseModel
}

func NewEventModel(db *gorm.DB) *EventModel {
	return &EventModel{NewBaseModel(new(Event), db)}
}

func (em *EventModel) Exists(startTime int64) (bool, error) {
	db := em.Base
	var count int
	err := db.Where("start_time = ?", startTime).Count(&count).Error
	if err != nil {
		return false, err
	}
	exists := count != 0
	return exists, nil
}

func (em *EventModel) ListEvents(parentId, taskId, limit, offset int) (Events, int, error) {
	var events Events
	var count int
	db := em.Base.DB
	if parentId > 0 {
		db = em.Base.Where("parent_id = ?", parentId)
	}
	if taskId > 0 {
		db = em.Base.Where("task_id = ?", taskId)
	}

	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(limit).Offset(offset).Scan(&events).Error
	if err != nil {
		return nil, 0, err
	}
	return events, count, nil
}
