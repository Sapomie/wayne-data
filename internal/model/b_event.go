package model

import (
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"github.com/Sapomie/wayne-data/pkg/convert"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

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

func (e *Event) TableName() string {
	return "b_event"
}

func (e *Event) Start() time.Time {
	return time.Unix(e.StartTime, 0)
}

func (e *Event) End() time.Time {
	return time.Unix(e.EndTime, 0)
}

func (e *Event) StuffIds() (stuffIds []int) {
	ids := strings.Split(e.StuffId, ",")
	for _, stuffId := range ids {
		stuffIds = append(stuffIds, convert.StrTo(stuffId).MustInt())
	}
	return
}

func (e *Event) TagIds() (tagIds []int) {
	ids := strings.Split(e.TagId, ",")
	for _, tagId := range ids {
		tagIds = append(tagIds, convert.StrTo(tagId).MustInt())
	}
	return
}

type Events []*Event

func (ets Events) Between(start, end time.Time) (events Events) {
	for _, event := range ets {
		if event.StartTime >= start.Unix() && event.EndTime <= end.Unix() {
			events = append(events, event)
		}
	}
	return
}

func (ets Events) Duration() (duration float64) {
	for _, event := range ets {
		duration += event.Duration
	}
	return
}

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

func (em *EventModel) GetAll() (Events, int, error) {
	var events Events
	var count int
	err := em.Base.Scan(&events).Error
	if err != nil {
		return nil, 0, err
	}
	count = len(events)
	return events, count, nil
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

func (em *EventModel) Newest() (*Event, error) {
	db := em.Base
	event := new(Event)
	err := db.Order("end_time desc").First(event).Error
	if err != nil {
		return nil, err
	}
	return event, nil
}

func updateNewest() error {
	em := NewEventModel(global.DBEngine)
	evt, err := em.Newest()
	if err != nil {
		return err
	}
	cons.Newest = time.Unix(evt.EndTime, 0)
	return nil
}