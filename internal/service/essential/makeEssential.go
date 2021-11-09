package essential

import (
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/model/constant"
	"github.com/Sapomie/wayne-data/pkg/convert"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"time"
)

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

	return essentials, nil
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
	dur := eventsZone.Duration() / 24

	taskInfo, parentInfo, stuffInfo, projectInfo, tagInfo := columnInfo(eventsZone, dur, durTotal)

	dayHour := dayHourInfos(taskInfo, dur)
	primary := getPrimary(dayHour)
	goalPct := countGoal(dayHour[constant.DHOther], dayHour[constant.DHSelfEntertain], typ) * 100
	goalMaxPct := countGoal(dayHour[constant.DHOther]*dur/durTotal, dayHour[constant.DHSelfEntertain]*dur/durTotal, typ) * 100
	dailyPct := dayHour[constant.DHDaily] / constant.DailyFull * 100

	essential := &Essential{
		Date:           zone.DateString(),
		Type:           typ,
		DateNumber:     num,
		StartTime:      zoneStart,
		EndTime:        zoneEnd,
		Duration:       dur,
		DurationTotal:  durTotal,
		Primary:        primary,
		DayHour:        dayHour,
		GoalPercent:    goalPct,
		GoalMaxPercent: goalMaxPct,
		DailyPercent:   dailyPct,
		TaskInfo:       taskInfo,
		ParentInfo:     parentInfo,
		StuffInfo:      stuffInfo,
		ProjectInfo:    projectInfo,
		TagInfo:        tagInfo,
	}

	return essential, nil
}

func columnInfo(events model.Events, dur, durTotal float64) (taskInfo, parentInfo, stuffInfo, projectInfo, tagInfo map[string]*FieldInfo) {
	taskInfo = make(map[string]*FieldInfo)
	parentInfo = make(map[string]*FieldInfo)
	projectInfo = make(map[string]*FieldInfo)

	stuffInfo = make(map[string]*FieldInfo)
	tagInfo = make(map[string]*FieldInfo)

	for _, evt := range events {
		if evt.TaskId > 0 {
			taskName := model.TaskInfoById[evt.TaskId].Name
			_, ok := taskInfo[taskName]
			if !ok {
				taskInfo[taskName] = new(FieldInfo)
				taskInfo[taskName].TenGoal = model.TaskInfoById[evt.TaskId].TenGoal
				taskInfo[taskName].DayHourType = model.TaskInfoById[evt.TaskId].DayHourType
			}
			taskInfo[taskName].Done += evt.Duration
		}

		if evt.ParentId > 0 {
			parent := model.ParentInfoById[evt.ParentId].Name
			_, ok := parentInfo[parent]
			if !ok {
				parentInfo[parent] = new(FieldInfo)
				parentInfo[parent].TenGoal = model.ParentInfoById[evt.ParentId].TenGoal
			}
			parentInfo[parent].Done += evt.Duration
		}

		if evt.ProjectId > 0 {
			project := model.ProjectInfoById[evt.ProjectId].Name
			_, ok := projectInfo[project]
			if !ok {
				projectInfo[project] = new(FieldInfo)
			}
			projectInfo[project].Done += evt.Duration
		}

		for _, stuffId := range evt.StuffIds() {
			if stuffId > 0 {
				stuff := model.StuffInfoById[stuffId].Name
				_, ok := stuffInfo[stuff]
				if !ok {
					stuffInfo[stuff] = new(FieldInfo)
					stuffInfo[stuff].TenGoal = model.StuffInfoById[stuffId].TenGoal
				}
				stuffInfo[stuff].Done += 1
			}
		}

		for _, tagId := range evt.TagIds() {
			if tagId > 0 {
				tag := model.TagInfoById[tagId].Name
				_, ok := tagInfo[tag]
				if !ok {
					tagInfo[tag] = new(FieldInfo)
				}
				tagInfo[tag].Done += evt.Duration
			}
		}
	}

	//percent,percentAbs
	for k, v := range taskInfo {
		if v.TenGoal > 0 {
			taskInfo[k].Percent = v.Done / (v.TenGoal / 10.0 * dur) * 100
			taskInfo[k].PercentAbs = v.Done / (v.TenGoal / 10.0 * durTotal) * 100
		}
	}

	for k, v := range parentInfo {
		if v.TenGoal > 0 {
			parentInfo[k].Percent = v.Done / (v.TenGoal / 10.0 * dur) * 100
			parentInfo[k].PercentAbs = v.Done / (v.TenGoal / 10.0 * durTotal) * 100
		}
	}

	for k, v := range stuffInfo {
		if v.TenGoal > 0 {
			stuffInfo[k].Percent = v.Done / (v.TenGoal / 10.0 * dur) * 100
			stuffInfo[k].PercentAbs = v.Done / (v.TenGoal / 10.0 * durTotal) * 100
		}
	}

	//make sure main field not nil
	for _, task := range constant.MainTasks {
		_, ok := taskInfo[task]
		if !ok {
			taskInfo[task] = &FieldInfo{
				TenGoal: model.TaskInfoByName[task].TenGoal,
			}
		}
	}

	for _, parent := range constant.MainParents {
		_, ok := parentInfo[parent]
		if !ok {
			parentInfo[parent] = &FieldInfo{
				TenGoal: model.ParentInfoByName[parent].TenGoal,
			}
		}
	}

	for _, stuff := range constant.MainStuffs {
		_, ok := stuffInfo[stuff]
		if !ok {
			stuffInfo[stuff] = &FieldInfo{
				TenGoal: model.StuffInfoByName[stuff].TenGoal,
			}
		}
	}

	return
}

func getPrimary(dayHour map[string]float64) float64 {

	return dayHour[constant.DHDaily] + dayHour[constant.DHOther]*(constant.DailyFull*constant.CountGoalBase/constant.OtherFull) + dayHour[constant.DHSelfEntertain]*(constant.DailyFull*constant.CountGoalBase/constant.SelfFull)
}

func dayHourInfos(taskInfo map[string]*FieldInfo, dur float64) map[string]float64 {

	dayHour := make(map[string]float64)
	dayHourDecimal := make(map[string]float64)

	for _, info := range taskInfo {
		switch info.DayHourType {
		case constant.DayHourOther:
			dayHour[constant.DHOther] += info.Done / dur
		case constant.DayHourDaily:
			dayHour[constant.DHDaily] += info.Done / dur
		case constant.DayHourSelfEntertain:
			dayHour[constant.DHSelfEntertain] += info.Done / dur
		case constant.DayHourSleep:
			dayHour[constant.DHSleep] += info.Done / dur
		case constant.DayHourRoutine:
			dayHour[constant.DHRoutine] += info.Done / dur
		case constant.DayHourBlank:
			dayHour[constant.DHBlank] += info.Done / dur
		}
	}

	for k, v := range dayHour {
		dayHourDecimal[k] = convert.FloatTo(v).Decimal(2)
	}

	for _, dayHourName := range constant.DayHourNames {
		_, ok := dayHour[dayHourName]
		if !ok {
			dayHourDecimal[dayHourName] = 0.0
		}
	}

	return dayHourDecimal
}

func countGoal(dayHourOther, dayHourSelf float64, typ mtime.TimeType) (goalPct float64) {

	var addition float64
	baseDaily := constant.CountGoalBase * constant.DailyFull

	switch typ {
	case mtime.TypeTen:
		addition = 0.0
	case mtime.TypeMonth:
		addition = 0.05
	case mtime.TypeYear, mtime.TypeHalf, mtime.TypeDay, mtime.TypeQuarter:
		addition = 0.15
	}

	dayOtherPct := dayHourOther / constant.OtherFull
	daySelfPct := dayHourSelf / constant.SelfFull

	goalPct = (baseDaily-baseDaily*dayOtherPct-baseDaily*daySelfPct)/baseDaily*constant.CountGoalBase + addition
	if goalPct < addition {
		goalPct = addition
	}

	return
}
