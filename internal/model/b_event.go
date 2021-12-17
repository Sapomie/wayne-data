package model

import (
	"fmt"
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"github.com/Sapomie/wayne-data/internal/model/resp"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"github.com/jinzhu/gorm"
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

	StuffId   int    `json:"stuff_id"`
	TagId     int    `json:"tag_id"`
	ProjectId int    `json:"project_id"`
	Remark    string `json:"remark"` //comment 除去自定义属性的部分

	*Base
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

type Events []*Event

func (ets Events) Between(start, end time.Time) (events Events) {
	for _, event := range ets {
		if event.StartTime >= start.Unix() && event.EndTime <= end.Unix() {
			events = append(events, event)
		}
	}
	return
}

func (ets Events) WithProject() (events Events) {
	for _, event := range ets {
		if event.ProjectId > 0 {
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

func (ets Events) Newest() (newest *Event) {
	newest = new(Event)
	for _, event := range ets {
		if event.StartTime >= newest.StartTime {
			newest = event
		}
	}
	return
}

func (ets Events) Oldest() (oldest *Event) {
	oldest = new(Event)
	oldest.StartTime = cons.BiggestTime
	for _, event := range ets {
		if event.StartTime <= oldest.StartTime {
			oldest = event
		}
	}
	return
}

type EventModel struct {
	Base *BaseDbModel
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

func (em *EventModel) YearEvents(year int) (Events, error) {
	var events Events
	start, end := mtime.NewTimeZone(mtime.TypeYear, year, 1).BeginAndEnd()
	err := em.Base.Where("start_time>= ? and start_time < ?", start.Unix(), end.Unix()).Scan(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (em *EventModel) Timezone(zone *mtime.TimeZone) (Events, error) {
	var events Events
	start, end := zone.BeginAndEnd()
	err := em.Base.Where("start_time>= ? and start_time < ?", start.Unix(), end.Unix()).Scan(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (em *EventModel) ListEvents(p *resp.DbEventListRequest, limit, offset int) (Events, int, error) {
	var events Events
	var count int
	db := em.Base.Where("start_time >= ? and end_time <= ?", p.Start.Unix(), p.End.Unix())
	if p.ParentId > 0 {
		db = db.Where("parent_id = ?", p.ParentId)
	}
	if p.TaskId > 0 {
		db = db.Where("task_id = ?", p.TaskId)
	}
	if p.StuffId > 0 {
		db = db.Where("stuff_id = ?", p.StuffId)
	}
	if p.TagId > 0 {
		db = db.Where("tag_id = ?", p.TagId)
	}
	if p.ProjectId > 0 {
		db = db.Where("project_id = ?", p.ProjectId)
	}
	if p.Word != "" {
		db = db.Where("comment like ?", "%"+p.Word+"%")
	}

	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(limit).Offset(offset).Order("id desc").Scan(&events).Error
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

func (em *EventModel) Oldest() (*Event, error) {
	db := em.Base
	event := new(Event)
	err := db.Order("end_time asc").First(event).Error
	if err != nil {
		return nil, err
	}
	return event, nil
}

//get events during start time to end time
func (em *EventModel) ByTaskName(start, end time.Time, name string) (Events, error) {
	db := em.Base
	var events Events
	err := db.
		Where("start_time >= ? and end_time <= ?", start.Unix(), end.Unix()).
		Where("task_id = ?", TaskInfoByName[name].Id).
		Scan(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}

//get events during start time to end time
func (em *EventModel) ByParentName(start, end time.Time, name string) (Events, error) {
	db := em.Base
	var events Events
	err := db.
		Where("start_time >= ? and end_time <= ?", start.Unix(), end.Unix()).
		Where("parent_id = ?", ParentInfoByName[name].Id).
		Order("start_time asc").
		Scan(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}

//get events during start time to end time
func (em *EventModel) ByProjectName(start, end time.Time, name string) (Events, error) {
	db := em.Base
	var events Events
	err := db.
		Where("start_time >= ? and end_time <= ?", start.Unix(), end.Unix()).
		Where("project_id = ?", ProjectInfoByName[name].Id).
		Order("start_time asc").
		Scan(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}

//get events during start time to end time
func (em *EventModel) ByStuffName(start, end time.Time, name string) (Events, error) {
	db := em.Base
	var events Events
	stuffId := StuffInfoByName[name].Id
	sql := fmt.Sprintf(`select * from b_event where start_time >= %v and end_time <= %v and (stuff_id = '%v' or stuff_id like "%v,%%" or stuff_id like "%%,%v,%%" or stuff_id like "%%,%v")`,
		start.Unix(), end.Unix(), stuffId, stuffId, stuffId, stuffId)
	err := db.
		Raw(sql).
		Order("start_time asc").
		Scan(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}

//get events during start time to end time
func (em *EventModel) ByTagName(start, end time.Time, name string) (Events, error) {
	db := em.Base
	var events Events
	stuffId := TagInfoByName[name].Id
	sql := fmt.Sprintf(`select * from b_event where start_time >= %v and end_time <= %v and (tag_id = '%v' or tag_id like "%v,%%" or tag_id like "%%,%v,%%" or tag_id like "%%,%v")`,
		start.Unix(), end.Unix(), stuffId, stuffId, stuffId, stuffId)
	err := db.
		Raw(sql).
		Order("start_time asc").
		Scan(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}

//get events during start time to end time
func (em *EventModel) ByTaskNames(start, end time.Time, names ...string) (Events, error) {
	db := em.Base
	var events Events
	var ids []int
	for _, name := range names {
		ids = append(ids, TaskInfoByName[name].Id)
	}

	err := db.
		Where("start_time >= ? and end_time <= ?", start.Unix(), end.Unix()).
		Where("task_id in (?)", ids).
		Order("start_time asc").
		Scan(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}

//get events during start time to end time
func (em *EventModel) WithProject(start, end time.Time) (Events, error) {
	db := em.Base
	var events Events
	err := db.
		Where("start_time >= ? and end_time <= ?", start.Unix(), end.Unix()).
		Where("project_id > 0").
		Scan(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (em *EventModel) UpdateNewestAndOldest() error {
	newest, err := em.Newest()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	cons.DbNewest = time.Unix(newest.StartTime, 0)

	oldest, err := em.Oldest()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	cons.DbOldest = time.Unix(oldest.StartTime, 0)
	return nil
}
