package model

import (
	"fmt"
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"github.com/jinzhu/gorm"
)

type Stuff struct {
	Id            int     `gorm:"primary_key"`
	Name          string  `gorm:"not null;unique"`
	EventNum      int64   `gorm:"not null" json:"tag_num"`
	TenGoal       float64 `gorm:"not null"`
	TotalDuration float64 `gorm:"not null" json:"total_duration"`
	FirstTime     int64   `gorm:"not null"`
	LastTime      int64   `gorm:"not null"`
	Longest       int64   `gorm:"not null"`

	*Base
}

func (s *Stuff) TableName() string {
	return "b_stuff"
}

func (s *Stuff) FieldName() string {
	return s.Name
}

func (s *Stuff) FieldTotalDuration() float64 {
	return s.TotalDuration
}

func (s *Stuff) FieldEventNum() int64 {
	return s.EventNum
}

func (s *Stuff) FieldFirstTimeAndLastTime() (int64, int64) {
	return s.FirstTime, s.LastTime
}

func (s *Stuff) FieldLongest() int64 {
	return s.Longest
}

type Stuffs []*Stuff

func (ss Stuffs) ToEventFields() []EventField {
	eventFields := make([]EventField, 0)
	for _, stuff := range ss {
		var ef EventField
		ef = stuff
		eventFields = append(eventFields, ef)
	}
	return eventFields
}

type StuffModel struct {
	Base *BaseDbModel
}

func NewStuffModel(db *gorm.DB) *StuffModel {
	return &StuffModel{NewBaseModel(new(Stuff), db)}
}

func (em *StuffModel) Exists(name string) (bool, error) {
	var count int
	err := em.Base.Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	exists := count > 0
	return exists, nil
}

func (em *StuffModel) GetAll() (Stuffs, error) {
	var stuffs Stuffs
	err := em.Base.Order("last_time desc").Scan(&stuffs).Error
	if err != nil {
		return nil, err
	}
	return stuffs, nil
}

func (em *StuffModel) ByName(name string) (*Stuff, error) {
	stuff := new(Stuff)
	err := em.Base.Where("name = ?", name).Scan(stuff).Error
	if err != nil {
		return nil, err
	}
	return stuff, nil
}

func (em *StuffModel) ListStuffs(limit, offset int) (Stuffs, int, error) {
	var (
		stuffs Stuffs
		count  int
	)
	db := em.Base.DB
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(limit).Offset(offset).Scan(&stuffs).Error
	if err != nil {
		return nil, 0, err
	}
	return stuffs, count, nil
}

func InsertAndGetStuff(db *gorm.DB, stuff *Stuff) (stuffDb *Stuff, info string, err error) {

	em := NewStuffModel(db)
	exists, err := em.Exists(stuff.Name)
	if err != nil {
		return nil, "", err
	}
	if !exists {
		err = em.Base.Create(stuff).Error
		if err != nil {
			return nil, "", err
		}
		info = fmt.Sprintf("Add Stuff %v ", stuff.Name)
	} else {
		err = em.Base.Where("name = ?", stuff.Name).Update(stuff).Error
		if err != nil {
			return nil, "", err
		}
	}

	stuffDb, err = NewStuffModel(db).ByName(stuff.Name)
	if err != nil {
		return nil, "", err
	}
	return
}

var StuffInfoById = make(map[int]struct {
	Name    string
	TenGoal float64
})

var StuffInfoByName = make(map[string]struct {
	Id      int
	TenGoal float64
})

func (em *StuffModel) UpdateStuffVariables() error {

	stuffs, err := em.GetAll()
	if err != nil {
		return err
	}
	for _, stuff := range stuffs {
		StuffInfoById[stuff.Id] = struct {
			Name    string
			TenGoal float64
		}{Name: stuff.Name, TenGoal: stuff.TenGoal}

		StuffInfoByName[stuff.Name] = struct {
			Id      int
			TenGoal float64
		}{Id: stuff.Id, TenGoal: stuff.TenGoal}
	}

	return nil
}

func UpdateStuffColumn(db *gorm.DB) (err error) {
	if err = NewStuffModel(db).UpdateStuffVariables(); err != nil {
		return err
	}

	stuffs, err := NewStuffModel(db).GetAll()
	if err != nil {
		return nil
	}
	for _, stuff := range stuffs {
		var data struct {
			Num            int64
			Dur            float64
			FirstTimestamp int64
			LastTimestamp  int64
		}
		err := NewEventModel(db).Base.Select("sum(duration) as dur, count(id) as num").Where("stuff_id = ?", stuff.Id).Scan(&data).Error
		if err != nil {
			return err
		}
		err = NewEventModel(db).Base.Select("start_time as first_timestamp").Where("stuff_id = ?", stuff.Id).Order("start_time asc").Limit(1).Scan(&data).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		err = NewEventModel(db).Base.Select("start_time as last_timestamp").Where("stuff_id = ?", stuff.Id).Order("start_time desc").Limit(1).Scan(&data).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		evts, err := NewEventModel(db).ByStuffName(cons.Oldest, cons.Newest, stuff.Name)
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
		err = NewStuffModel(db).Base.Where("id = ?", stuff.Id).Update(map[string]interface{}{
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
