package model

import (
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"github.com/jinzhu/gorm"
)

const (
	RunMinRate        = 60
	RunMinAltitude    = -415
	RunMinTemperature = -273
)

type Run struct {
	Id          int    `gorm:"primary_key" json:"id"`
	StartTime   int64  `gorm:"not null;unique" `
	Date        string `gorm:"not null" json:"date"`
	Distance    int    `gorm:"not null" json:"distance"`
	TimeCost    int    `gorm:"not null" json:"time_cost"`
	Pace        int    `gorm:"not null" json:"pace"`
	Rate        int    `json:"rate"`
	Temperature int    `json:"temperature"`
	Altitude    int    `json:"altitude"`
	Comment     string `json:"comment"`
	CreatedTime int64  `gorm:"not null" json:"created_time"`
	UpdatedTime int64  `gorm:"not null" json:"updated_time"`
}

func (e *Run) TableName() string {
	return "c_run"
}

type Runs []*Run

type RunModel struct {
	Base *BaseDbModel
}

func NewRunModel(db *gorm.DB) *RunModel {
	return &RunModel{NewBaseModel(new(Run), db)}
}

func (em *RunModel) Exists(startTime int64) (bool, error) {
	var count int
	err := em.Base.Where("start_time = ?", startTime).Count(&count).Error
	if err != nil {
		return false, err
	}
	exists := count > 0
	return exists, nil
}

func (em *RunModel) GetAll() (Runs, error) {
	var runs Runs
	err := em.Base.Order("start_time desc").Scan(&runs).Error
	if err != nil {
		return nil, err
	}
	return runs, nil
}

func (em *RunModel) Timezone(zone *mtime.TimeZone) (Runs, error) {
	var runs Runs
	start, end := zone.BeginAndEnd()
	err := em.Base.Where("start_time>= ? and start_time < ?", start.Unix(), end.Unix()).Scan(&runs).Error
	if err != nil {
		return nil, err
	}
	return runs, nil
}
