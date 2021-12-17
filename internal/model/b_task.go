package model

import (
	"fmt"
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"github.com/jinzhu/gorm"
)

type Task struct {
	Id            int     `gorm:"primary_key"`
	Name          string  `gorm:"not null;unique"`
	TenGoal       float64 `gorm:"not null"`
	DayHourType   int     `gorm:"not null"`
	EventNum      int64   `gorm:"not null" json:"tag_num"`
	TotalDuration float64 `gorm:"not null" json:"total_duration"`
	FirstTime     int64   `gorm:"not null"`
	LastTime      int64   `gorm:"not null"`
	Longest       int64   `gorm:"not null"`

	*Base
}

func (t *Task) TableName() string {
	return "b_task"
}

func (t *Task) FieldName() string {
	return t.Name
}

func (t *Task) FieldTotalDuration() float64 {
	return t.TotalDuration
}

func (t *Task) FieldEventNum() int64 {
	return t.EventNum
}

func (t *Task) FieldFirstTimeAndLastTime() (int64, int64) {
	return t.FirstTime, t.LastTime
}

func (t *Task) FieldLongest() int64 {
	return t.Longest
}

type Tasks []*Task

func (ts Tasks) ToEventFields() []EventField {
	eventFields := make([]EventField, 0)
	for _, task := range ts {
		var ef EventField
		ef = task
		eventFields = append(eventFields, ef)
	}
	return eventFields
}

type TaskModel struct {
	Base *BaseDbModel
}

func NewTaskModel(db *gorm.DB) *TaskModel {
	return &TaskModel{NewBaseModel(new(Task), db)}
}

func (em *TaskModel) Exists(name string) (bool, error) {
	var count int
	err := em.Base.Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	exists := count > 0
	return exists, nil
}

func (em *TaskModel) GetAll() (Tasks, error) {
	var tasks Tasks
	err := em.Base.Order("last_time desc").Scan(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (em *TaskModel) ByName(name string) (*Task, error) {
	task := new(Task)
	err := em.Base.Where("name = ?", name).Scan(task).Error
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (em *TaskModel) ById(id int) (*Task, error) {
	task := new(Task)
	err := em.Base.Where("id = ?", id).Scan(task).Error
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (em *TaskModel) ListTasks(limit, offset int) (Tasks, int, error) {
	var (
		tasks Tasks
		count int
	)
	db := em.Base.DB
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(limit).Offset(offset).Scan(&tasks).Error
	if err != nil {
		return nil, 0, err
	}
	return tasks, count, nil
}

func InsertAndGetTask(db *gorm.DB, task *Task) (taskDb *Task, info string, err error) {

	em := NewTaskModel(db)
	exists, err := em.Exists(task.Name)
	if err != nil {
		return nil, "", err
	}
	if !exists {
		err = em.Base.Create(task).Error
		if err != nil {
			return nil, "", err
		}
		info = fmt.Sprintf("Add Task %v ", task.Name)
	} else {
		//update struct 且struct 带有主键（id）时，em.Base.DB变量受影响（之后的查询会带有where id = *）
		err = em.Base.Where("name = ?", task.Name).Update(task).Error
		if err != nil {
			return nil, "", err
		}
	}

	taskDb, err = NewTaskModel(db).ByName(task.Name)
	if err != nil {
		return nil, "", err
	}
	return
}

var TaskInfoById = make(map[int]struct {
	Name        string
	TenGoal     float64
	DayHourType int
})

var TaskInfoByName = make(map[string]struct {
	Id          int
	TenGoal     float64
	DayHourType int
})

func (em *TaskModel) UpdateTaskVariables() error {

	tasks, err := em.GetAll()
	if err != nil {
		return err
	}
	cons.DailyFull = 0
	for _, task := range tasks {
		TaskInfoById[task.Id] = struct {
			Name        string
			TenGoal     float64
			DayHourType int
		}{Name: task.Name, TenGoal: task.TenGoal, DayHourType: task.DayHourType}

		TaskInfoByName[task.Name] = struct {
			Id          int
			TenGoal     float64
			DayHourType int
		}{Id: task.Id, TenGoal: task.TenGoal, DayHourType: task.DayHourType}
		if task.DayHourType == cons.DayHourDaily {
			cons.DailyFull += task.TenGoal / 10
		}
	}

	return nil
}

func UpdateTaskColumn(db *gorm.DB) (err error) {
	if err = NewTaskModel(db).UpdateTaskVariables(); err != nil {
		return err
	}

	tasks, err := NewTaskModel(db).GetAll()
	if err != nil {
		return nil
	}
	for _, task := range tasks {
		var data struct {
			Num            int64
			Dur            float64
			FirstTimestamp int64
			LastTimestamp  int64
		}
		err := NewEventModel(db).Base.Select("sum(duration) as dur, count(id) as num").Where("task_id = ?", task.Id).Scan(&data).Error
		if err != nil {
			return err
		}
		err = NewEventModel(db).Base.Select("start_time as first_timestamp").Where("task_id = ?", task.Id).Order("start_time asc").Limit(1).Scan(&data).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		err = NewEventModel(db).Base.Select("start_time as last_timestamp").Where("task_id = ?", task.Id).Order("start_time desc").Limit(1).Scan(&data).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		evts, err := NewEventModel(db).ByTaskName(cons.DbOldest, cons.DbNewest, task.Name)
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
		err = NewTaskModel(db).Base.Where("id = ?", task.Id).Update(map[string]interface{}{
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
