package essential

import (
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"time"
)

type Essential struct {
	Date          string
	Type          mtime.TimeType
	DateNumber    int
	StartTime     time.Time
	EndTime       time.Time
	Duration      float64
	DurationTotal float64
	Primary       float64 //todo: rename progressPoint

	DayHour        map[string]float64
	GoalPercent    float64
	GoalMaxPercent float64
	DailyPercent   float64

	TaskInfo    map[string]*FieldInfo
	ParentInfo  map[string]*FieldInfo
	StuffInfo   map[string]*FieldInfo
	ProjectInfo map[string]*FieldInfo
	TagInfo     map[string]*FieldInfo
}

type FieldInfo struct {
	Done        float64
	Percent     float64
	PercentAbs  float64
	WeekGoal    float64
	DayHourType int
}

type Essentials []*Essential

//events 为数据库里的所有
func MakeEssentials(events model.Events, start, end time.Time, typ mtime.TimeType) (Essentials, error) {
	start, end = mtime.TrimTime(start, end)
	year := start.Year()
	startNumber := mtime.NewMTime(start).TimeZone(typ).Number()
	endNumber := mtime.NewMTime(end).TimeZone(typ).Number()

	var essentials Essentials
	for i := startNumber; i <= endNumber; i++ {
		essential, err := makeEssential(events, start, year, i, typ)
		if err != nil {
			return nil, err
		}
		if essential != nil && essential.Duration > 0 {
			essentials = append(essentials, essential)
		}
	}

	return nil, nil
}

//start 为自定开始日期，默认应该为初始年第一天
func makeEssential(events model.Events, start time.Time, year, num int, typ mtime.TimeType) (ess *Essential, err error) {
	zone := mtime.NewTimeZone(typ, year, num)
	zoneStart, zoneEnd := zone.BeginAndEnd()
	eventsZone := events.Between(zoneStart, zoneEnd)

	var durTotal float64
	if start.Unix() > zoneStart.Unix() {
		durTotal = float64((zoneEnd.Unix() - start.Unix()) / 60 / 60 / 24)
	} else {
		durTotal = float64((zoneEnd.Unix() - zoneStart.Unix()) / 60 / 60 / 24)
	}
	dur := eventsZone.Duration()

	return nil, nil
}

func columnInfo(events model.Events, dur, durTotal float64) (taskInfo, parentInfo, stuffInfo, projectInfo, categoryInfo, tagInfo map[string]*FieldInfo) {
	taskInfo = make(map[string]*FieldInfo)
	parentInfo = make(map[string]*FieldInfo)
	stuffInfo = make(map[string]*FieldInfo)
	projectInfo = make(map[string]*FieldInfo)
	categoryInfo = make(map[string]*FieldInfo)
	tagInfo = make(map[string]*FieldInfo)

	for _, evt := range events {
		if evt.TaskId > 0 {
			taskName := fd.TaskNameInfo[evt.TaskNameId].Name
			_, ok := taskInfo[taskName]
			if !ok {
				taskInfo[taskName] = new(models.FieldInfo)
				taskInfo[taskName].WeekGoal = fd.TaskNameInfo[evt.TaskNameId].WeekGoal
				taskInfo[taskName].DayHourType = fd.TaskNameInfo[evt.TaskNameId].DayHourType
			}
			taskInfo[taskName].Done += evt.Duration
		}

		if evt.ParentTaskId > 0 {
			parent := fd.ParentTaskInfo[evt.ParentTaskId].Name
			_, ok := parentInfo[parent]
			if !ok {
				parentInfo[parent] = new(models.FieldInfo)
				parentInfo[parent].WeekGoal = fd.ParentTaskInfo[evt.ParentTaskId].WeekGoal
			}
			parentInfo[parent].Done += evt.Duration
		}

		if evt.StuffId > 0 {
			stuff := fd.StuffInfo[evt.StuffId].Name
			_, ok := stuffInfo[stuff]
			if !ok {
				stuffInfo[stuff] = new(models.FieldInfo)
				stuffInfo[stuff].WeekGoal = fd.StuffInfo[evt.StuffId].WeekGoal
			}
			stuffInfo[stuff].Done++
		}

		if evt.ProjectId > 0 {
			project := fd.ProjectInfo[evt.ProjectId]
			_, ok := projectInfo[project]
			if !ok {
				projectInfo[project] = new(models.FieldInfo)
			}
			projectInfo[project].Done += evt.Duration
		}

		if evt.CategoryId > 0 {
			category := fd.CategoryInfo[evt.CategoryId]
			_, ok := categoryInfo[category]
			if !ok {
				categoryInfo[category] = new(models.FieldInfo)
			}
			categoryInfo[category].Done += evt.Duration
		}

		if evt.TagId > 0 {
			tag := fd.TagInfo[evt.TagId]
			_, ok := taskInfo[tag]
			if !ok {
				tagInfo[tag] = new(models.FieldInfo)
			}
			tagInfo[tag].Done += evt.Duration
		}

	}

	for k, v := range taskInfo {
		if v.WeekGoal > 0 {
			taskInfo[k].Percent = v.Done / (v.WeekGoal / 7.0 * dur) * 100
			taskInfo[k].PercentAbs = v.Done / (v.WeekGoal / 7.0 * durTotal) * 100
		}
	}

	for k, v := range stuffInfo {
		if v.WeekGoal > 0 {
			stuffInfo[k].Percent = v.Done / (v.WeekGoal / 7.0 * dur) * 100
			stuffInfo[k].PercentAbs = v.Done / (v.WeekGoal / 7.0 * durTotal) * 100
		}
	}

	for k, v := range parentInfo {
		if v.WeekGoal > 0 {
			parentInfo[k].Percent = v.Done / (v.WeekGoal / 7.0 * dur) * 100
			parentInfo[k].PercentAbs = v.Done / (v.WeekGoal / 7.0 * durTotal) * 100
		}
	}

	for _, parent := range tk.MainParentTasks {
		_, ok := parentInfo[parent]
		if !ok {
			parentInfo[parent] = &models.FieldInfo{
				WeekGoal: fd.ParentTaskInfo[fd.ReverseParentTaskInfo[parent]].WeekGoal,
			}
		}
	}

	for _, stuff := range tk.MainStuff {
		_, ok := stuffInfo[stuff]
		if !ok {
			stuffInfo[stuff] = &models.FieldInfo{
				WeekGoal: fd.StuffInfo[fd.ReverseStuffInfo[stuff]].WeekGoal,
			}
		}
	}

	for _, task := range tk.MainTasks {
		_, ok := taskInfo[task]
		if !ok {
			taskInfo[task] = &models.FieldInfo{
				WeekGoal: fd.TaskNameInfo[fd.ReverseTaskNameInfo[task]].WeekGoal,
			}
		}
	}

	return
}
