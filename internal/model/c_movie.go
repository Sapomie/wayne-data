package model

import (
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"github.com/jinzhu/gorm"
)

type Movie struct {
	Id        int    `gorm:"primary_key"`
	Date      string `gorm:"not null"`
	Name      string `gorm:"not null"`
	EnName    string
	Rate      int
	Place     int //1:own,2:theater,3:other
	Year      int
	StartTime int64 `gorm:"not null"`
	Comment   string

	*Base
}

func (e *Movie) TableName() string {
	return "c_movie"
}

type Movies []*Movie

type MovieModel struct {
	Base *BaseDbModel
}

func NewMovieModel(db *gorm.DB) *MovieModel {
	return &MovieModel{NewBaseModel(new(Movie), db)}
}

func (em *MovieModel) Exists(startTime int64) (bool, error) {
	var count int
	err := em.Base.Where("start_time = ?", startTime).Count(&count).Error
	if err != nil {
		return false, err
	}
	exists := count > 0
	return exists, nil
}

func (em *MovieModel) GetAll() (Movies, error) {
	var movies Movies
	err := em.Base.Order("start_time desc").Scan(&movies).Error
	if err != nil {
		return nil, err
	}
	return movies, nil
}

func (em *MovieModel) Timezone(zone *mtime.TimeZone) (Movies, error) {
	var movies Movies
	start, end := zone.BeginAndEnd()
	err := em.Base.Where("start_time>= ? and start_time < ?", start.Unix(), end.Unix()).Scan(&movies).Error
	if err != nil {
		return nil, err
	}
	return movies, nil
}
