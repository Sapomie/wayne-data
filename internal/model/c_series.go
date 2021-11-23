package model

import (
	"github.com/jinzhu/gorm"
)

const (
	SeriesFinish = 1
)

type Series struct {
	Id               int64  `gorm:"primary_key"`
	NameSeason       string `gorm:"not null"`
	Name             string `gorm:"not null"`
	NameOrigin       string
	Season           int
	Category         string
	Year             int
	EpisodeNumber    int
	Rate             int
	Duration         float64
	FirstReadingTime int64
	LastReadingTime  int64
	Finish           int8
	CreatedTime      int64 `gorm:"not null" json:"created_time"`
	UpdatedTime      int64 `gorm:"not null" json:"updated_time"`
}

func (e *Series) TableName() string {
	return "c_series"
}

type SeriesS []*Series

type SeriesModel struct {
	Base *BaseDbModel
}

func NewSeriesModel(db *gorm.DB) *SeriesModel {
	return &SeriesModel{NewBaseModel(new(Series), db)}
}

func (em *SeriesModel) Exists(nameSeason string) (bool, error) {
	db := em.Base
	var count int
	err := db.Where("name_season = ?", nameSeason).Count(&count).Error
	if err != nil {
		return false, err
	}
	exists := count != 0

	return exists, nil
}

func (em *SeriesModel) GetAll() (SeriesS, error) {
	var seriesS SeriesS
	err := em.Base.Order("last_time desc").Scan(&seriesS).Error
	if err != nil {
		return nil, err
	}
	return seriesS, nil
}

func (em *SeriesModel) ListSeriesS(limit, offset int) (SeriesS, int, error) {
	var (
		seriesS SeriesS
		count   int
	)
	db := em.Base.DB
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(limit).Offset(offset).Scan(&seriesS).Error
	if err != nil {
		return nil, 0, err
	}
	return seriesS, count, nil
}
