package model

import (
	"fmt"
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"github.com/jinzhu/gorm"
)

type Task struct {
	Id            int     `gorm:"primary_key"`
	Name          string  `gorm:"not null;unique"`
	TenGoal       float64 `gorm:"not null"`
	Point         float64 `gorm:"not null"`
	DayHourType   int     `gorm:"not null"`
	ParentTaskId  int     `gorm:"not null"`
	EventNum      int64   `gorm:"not null" json:"tag_num"`
	TotalDuration float64 `gorm:"not null" json:"total_duration"`
	FirstTime     int64   `gorm:"not null"`
	LastTime      int64   `gorm:"not null"`
	Longest       int64   `gorm:"not null"`

	*Model
}

func (e *Task) TableName() string {
	return "b_task"
}

type Tasks []*Task

type TaskModel struct {
	Base *BaseModel
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

func (em *TaskModel) InsertAndGetTask(name string) (task *Task, info string, err error) {
	exists, err := em.Exists(name)
	if err != nil {
		return nil, "", err
	}
	if !exists {
		err = em.Base.Create(&Task{Name: name}).Error
		if err != nil {
			return nil, "", err
		}
		info = fmt.Sprintf("Add Task %v ", name)
	}
	task, err = em.ByName(name)
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

func updateTaskVariables() error {

	tasks, err := NewTaskModel(global.DBEngine).GetAll()
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
