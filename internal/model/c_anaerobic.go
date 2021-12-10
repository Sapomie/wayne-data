package model

import (
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"github.com/jinzhu/gorm"
)

const (
	PushUp        = "push-up"
	SitUp         = "sit-up"
	DumbbellPress = "dumbbellPress"
)

type Anaerobic struct {
	Id        int64 `gorm:"primary_key"`
	StartTime int64
	EndTime   int64
	Name      string `gorm:"not null"`
	Group     int
	Times     int
	Addition  float64

	*Base
}

func (e *Anaerobic) TableName() string {
	return "c_anaerobic"
}

type AnaerobicS []*Anaerobic

type AnaerobicModel struct {
	Base *BaseDbModel
}

func NewAnaerobicModel(db *gorm.DB) *AnaerobicModel {
	return &AnaerobicModel{NewBaseModel(new(Anaerobic), db)}
}

func (em *AnaerobicModel) Exists(startTime int64) (bool, error) {
	var count int
	err := em.Base.Where("start_time = ?", startTime).Count(&count).Error
	if err != nil {
		return false, err
	}
	exists := count > 0
	return exists, nil
}

func (em *AnaerobicModel) GetAll() (AnaerobicS, error) {
	var anaerobicS AnaerobicS
	err := em.Base.Order("start_time desc").Scan(&anaerobicS).Error
	if err != nil {
		return nil, err
	}
	return anaerobicS, nil
}

func (em *AnaerobicModel) Timezone(zone *mtime.TimeZone) (AnaerobicS, error) {
	var anaerobicS AnaerobicS
	start, end := zone.BeginAndEnd()
	err := em.Base.Where("start_time>= ? and start_time < ?", start.Unix(), end.Unix()).Scan(&anaerobicS).Error
	if err != nil {
		return nil, err
	}
	return anaerobicS, nil
}
