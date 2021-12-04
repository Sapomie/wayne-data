package model

import (
	"fmt"
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"github.com/jinzhu/gorm"
)

type Parent struct {
	Id            int     `gorm:"primary_key"`
	Name          string  `gorm:"not null;unique"`
	TenGoal       float64 `gorm:"not null"`
	EventNum      int64   `gorm:"not null" json:"tag_num"`
	TotalDuration float64 `gorm:"not null" json:"total_duration"`
	FirstTime     int64   `gorm:"not null"`
	LastTime      int64   `gorm:"not null"`
	Longest       int64   `gorm:"not null"`

	*Base
}

func (p *Parent) TableName() string {
	return "b_parent"
}

type Parents []*Parent

type ParentModel struct {
	Base *BaseDbModel
}

func NewParentModel(db *gorm.DB) *ParentModel {
	return &ParentModel{NewBaseModel(new(Parent), db)}
}

func (em *ParentModel) Exists(name string) (bool, error) {
	var count int
	err := em.Base.Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	exists := count > 0
	return exists, nil
}

func (em *ParentModel) GetAll() (Parents, error) {
	var parents Parents
	err := em.Base.Order("last_time desc").Scan(&parents).Error
	if err != nil {
		return nil, err
	}
	return parents, nil
}

func (em *ParentModel) ByName(name string) (*Parent, error) {
	parent := new(Parent)
	err := em.Base.Where("name = ?", name).Scan(parent).Error
	if err != nil {
		return nil, err
	}
	return parent, nil
}

func (em *ParentModel) ById(id int) (*Parent, error) {
	parent := new(Parent)
	err := em.Base.Where("id = ?", id).Scan(parent).Error
	if err != nil {
		return nil, err
	}
	return parent, nil
}

func (em *ParentModel) ListParents(limit, offset int) (Parents, int, error) {
	var (
		parents Parents
		count   int
	)
	db := em.Base.DB
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(limit).Offset(offset).Scan(&parents).Error
	if err != nil {
		return nil, 0, err
	}
	return parents, count, nil
}

var ParentInfoById = make(map[int]struct {
	Name    string
	TenGoal float64
})

var ParentInfoByName = make(map[string]struct {
	Id      int
	TenGoal float64
})

func (em *ParentModel) UpdateParentVariables() error {

	parents, err := em.GetAll()
	if err != nil {
		return err
	}
	for _, parent := range parents {
		ParentInfoById[parent.Id] = struct {
			Name    string
			TenGoal float64
		}{Name: parent.Name, TenGoal: parent.TenGoal}

		ParentInfoByName[parent.Name] = struct {
			Id      int
			TenGoal float64
		}{Id: parent.Id, TenGoal: parent.TenGoal}
	}

	return nil
}

func InsertAndGetParent(db *gorm.DB, parent *Parent) (parentDb *Parent, info string, err error) {

	em := NewParentModel(db)
	exists, err := em.Exists(parent.Name)
	if err != nil {
		return nil, "", err
	}
	if !exists {
		err = em.Base.Create(parent).Error
		if err != nil {
			return nil, "", err
		}
		info = fmt.Sprintf("Add Parent %v ", parent.Name)
	} else {
		err = em.Base.Where("name = ?", parent.Name).Update(parent).Error
		if err != nil {
			return nil, "", err
		}
	}

	parentDb, err = NewParentModel(db).ByName(parent.Name)
	if err != nil {
		return nil, "", err
	}
	return
}

func UpdateParentColumn(db *gorm.DB) (err error) {
	if err = NewParentModel(db).UpdateParentVariables(); err != nil {
		return err
	}

	parents, err := NewParentModel(db).GetAll()
	if err != nil {
		return nil
	}
	for _, parent := range parents {
		var data struct {
			Num            int64
			Dur            float64
			FirstTimestamp int64
			LastTimestamp  int64
		}
		err := NewEventModel(db).Base.Select("sum(duration) as dur, count(id) as num").Where("parent_id = ?", parent.Id).Scan(&data).Error
		if err != nil {
			return err
		}
		err = NewEventModel(db).Base.Select("start_time as first_timestamp").Where("parent_id = ?", parent.Id).Order("start_time asc").Limit(1).Scan(&data).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		err = NewEventModel(db).Base.Select("start_time as last_timestamp").Where("parent_id = ?", parent.Id).Order("start_time desc").Limit(1).Scan(&data).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		evts, err := NewEventModel(db).ByParentName(cons.Oldest, cons.Newest, parent.Name)
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
		err = NewParentModel(db).Base.Where("id = ?", parent.Id).Update(map[string]interface{}{
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
