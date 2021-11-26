package b_rawEvent

import (
	"encoding/json"
	"github.com/Sapomie/wayne-data/internal/model"
	"os"
)

type task struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	TenGoal     float64 `json:"ten_goal"`
	DayHourType int     `json:"day_hour_type"`
}

type parent struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	TenGoal float64 `json:"ten_goal"`
}

type stuff struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	TenGoal float64 `json:"ten_goal"`
}

type abbr struct {
	Id      int    `json:"id"`
	TaskId  int    `json:"task_id"`
	Abbr    string `json:"abbr"`
	Content string `json:"content"`
}

type DefaultValue struct {
	Tasks   []*task   `json:"tasks"`
	Parents []*parent `json:"parents"`
	Stuffs  []*stuff  `json:"stuffs"`
	Abbrs   []*abbr   `json:"abbrs"`
}

func (svc RawEventService) ReadDefaultTaskValue() error {
	byt, err := os.ReadFile(svc.appSetting.DefaultValuePath)
	if err != nil {
		return err
	}
	data := new(DefaultValue)

	err = json.Unmarshal(byt, data)
	if err != nil {
		return err
	}

	for _, task := range data.Tasks {
		//svc.taskDb = model.NewTaskModel(global.DBEngine)
		_, _, err = model.InsertAndGetTask(svc.db, &model.Task{
			Id:          task.Id,
			Name:        task.Name,
			TenGoal:     task.TenGoal,
			DayHourType: task.DayHourType,
		})
		if err != nil {
			return err
		}
	}
	for _, parent := range data.Parents {
		_, _, err = model.InsertAndGetParent(svc.db, &model.Parent{
			Name:    parent.Name,
			TenGoal: parent.TenGoal,
		})
		if err != nil {
			return err
		}
	}
	for _, stuff := range data.Stuffs {
		_, _, err = model.InsertAndGetStuff(svc.db, &model.Stuff{
			Name:    stuff.Name,
			TenGoal: stuff.TenGoal,
		})
		if err != nil {
			return err
		}
	}
	for _, abbr := range data.Abbrs {
		err = model.InsertAbbr(svc.db, &model.Abbr{
			TaskId:  abbr.TaskId,
			Abbr:    abbr.Abbr,
			Content: abbr.Content,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
