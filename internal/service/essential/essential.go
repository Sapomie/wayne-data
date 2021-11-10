package essential

import (
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"github.com/Sapomie/wayne-data/pkg/convert"
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

func (es *Essential) Decimal() {
	es.GoalPercent = convert.FloatTo(es.GoalPercent).Decimal(0)
	es.DailyPercent = convert.FloatTo(es.DailyPercent).Decimal(0)
	es.Primary = convert.FloatTo(es.Primary).Decimal(2)
	for _, v := range es.TaskInfo {
		v.Percent = convert.FloatTo(v.Percent).Decimal(0)
	}
	for _, v := range es.ParentInfo {
		v.Percent = convert.FloatTo(v.Percent).Decimal(0)
	}
	for _, v := range es.StuffInfo {
		v.Percent = convert.FloatTo(v.Percent).Decimal(0)
	}
}

func (es *Essential) giveMainColumnMapKey() {
	for _, task := range cons.MainTasks {
		_, ok := es.TaskInfo[task]
		if !ok {
			es.TaskInfo[task] = &FieldInfo{
				TenGoal: model.TaskInfoByName[task].TenGoal,
			}
		}
	}

	for _, parent := range cons.MainParents {
		_, ok := es.ParentInfo[parent]
		if !ok {
			es.ParentInfo[parent] = &FieldInfo{
				TenGoal: model.ParentInfoByName[parent].TenGoal,
			}
		}
	}

	for _, stuff := range cons.MainStuffs {
		_, ok := es.StuffInfo[stuff]
		if !ok {
			es.StuffInfo[stuff] = &FieldInfo{
				TenGoal: model.StuffInfoByName[stuff].TenGoal,
			}
		}
	}

}

type FieldInfo struct {
	Done        float64
	Percent     float64
	PercentAbs  float64
	TenGoal     float64
	DayHourType int
}

type Essentials []*Essential

func (ess Essentials) Decimal() {
	for _, es := range ess {
		es.Decimal()
	}
	return
}

func (ess Essentials) Response() {
	//decimal
	ess.Decimal()
	//running case
	for _, es := range ess {
		es.TaskInfo[cons.Running].Done *= 5
	}
	return
}
