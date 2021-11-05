package model

import "github.com/Sapomie/wayne-data/global"

var TaskInfoById = make(map[int]struct {
	Name        string
	WeekGoal    float64
	DayHourType int
})

var TaskInfoByName = make(map[string]struct {
	Id          int
	WeekGoal    float64
	DayHourType int
})

func updateVariable() error {
	tasks, err := NewTaskModel(global.DBEngine).GetAll()
	if err != nil {
		return err
	}
	for _, task := range tasks {
		TaskInfoById[task.Id] = struct {
			Name        string
			WeekGoal    float64
			DayHourType int
		}{Name: task.Name, WeekGoal: task.WeekGoal, DayHourType: task.DayHourType}

		TaskInfoByName[task.Name] = struct {
			Id          int
			WeekGoal    float64
			DayHourType int
		}{Id: task.Id, WeekGoal: task.WeekGoal, DayHourType: task.DayHourType}
	}

	return nil
}
