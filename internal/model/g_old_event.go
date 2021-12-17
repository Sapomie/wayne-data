package model

import (
	"github.com/Sapomie/wayne-data/internal/model/resp"
	"github.com/jinzhu/gorm"
)

type OldEvent struct {
	ID         uint `gorm:"primary_key"`
	Date       string
	TaskName   string
	Duration   float64
	Comment    string
	ParentTask string
	StartTime  int64
	EndTime    int64
}

func (e *OldEvent) TableName() string {
	return "g_old_event"
}

type OldEvents []*OldEvent

type OldEventModel struct {
	Base *BaseDbModel
}

func NewOldEventModel(db *gorm.DB) *OldEventModel {
	return &OldEventModel{NewBaseModel(new(OldEvent), db)}
}

func (em *OldEventModel) Exists(startTime int64) (bool, error) {
	db := em.Base
	var count int
	err := db.Where("start_time = ?", startTime).Count(&count).Error
	if err != nil {
		return false, err
	}
	exists := count != 0
	return exists, nil
}

func (em *OldEventModel) GetAll() (OldEvents, int, error) {
	var oldOldEvents OldEvents
	var count int
	err := em.Base.Scan(&oldOldEvents).Error
	if err != nil {
		return nil, 0, err
	}
	count = len(oldOldEvents)
	return oldOldEvents, count, nil
}

func (em *OldEventModel) ListOldEvents(p *resp.DbOldEventListRequest, limit, offset int) (OldEvents, int, error) {
	var oldOldEvents OldEvents
	var count int

	db := em.Base.Where("start_time >= ? and end_time <= ?", p.Start.Unix(), p.End.Unix())
	if p.Parent != "" {
		db = db.Where("parent_task = ?", p.Parent)
	}
	if p.Task != "" {
		db = db.Where("task_name = ?", p.Task)
	}
	if p.Word != "" {
		db = db.Where("comment like ?", "%"+p.Word+"%")
	}

	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(limit).Offset(offset).Order("id desc").Scan(&oldOldEvents).Error
	if err != nil {
		return nil, 0, err
	}
	return oldOldEvents, count, nil
}
