package progress

import (
	"fmt"
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"github.com/Sapomie/wayne-data/internal/service/essential"
	"github.com/Sapomie/wayne-data/pkg/convert"
	"strconv"
	"time"
)

type Progress struct {
	Date       string
	Days       string
	DaysPassed float64
	Task       map[string]*summaryField
	Parent     map[string]*summaryField
	Stuff      map[string]*summaryField
	DailyLimit *summaryField
	//Daily      *summaryField
	GoalLefts []*goalLeft
	Primary   float64
	GcRunning *GcRunning
}

type summaryField struct {
	Done       float64
	DonePerDay float64
	Percent    float64 //relative
	PercentAbs float64 //absolute
	AheadLevel int8
}

type goalLeft struct {
	Name       string
	TaskInfo   map[string]*goalLeftField
	ParentInfo map[string]*goalLeftField
	StuffInfo  map[string]*goalLeftField
	DailyLimit *goalLeftField
}

type goalLeftField struct {
	Left        string
	LeftF       float64
	LeftFPerDay float64
	Finish      int8
}

type GcRunning struct {
	RunningDistance float64
	GcAccumulation  int
	GcUsed          int
	GcLeft          int
}

func MakeProgress(es *essential.Essential, progressStartTime time.Time) *Progress {
	if es.StartTime.Unix() < progressStartTime.Unix() {
		es.StartTime = progressStartTime
	}

	fullPct := 1.00
	goalNowPct := es.GoalPercent / 100
	goalMaxPct := es.GoalMaxPercent / 100

	nowPctSummary, dailyLimit := makeGoalLeft(es, goalNowPct)
	maxPctSummary, _ := makeGoalLeft(es, goalMaxPct)
	fullPctSummary, _ := makeGoalLeft(es, fullPct)

	nowPctSummary.Name += cons.LevelUpMark
	maxPctSummary.Name += cons.LevelDownMark

	goalLefts := make([]*goalLeft, 0)
	goalLefts = append(
		goalLefts,
		nowPctSummary,
		maxPctSummary,
		fullPctSummary,
	)

	taskMap := makeMapValueFloatToString(es.TaskInfo, es.Duration, goalNowPct)
	parentMap := makeMapValueFloatToString(es.ParentInfo, es.Duration, goalNowPct)
	stuffMap := makeStuffMapValueFloatToString(es.StuffInfo, es.Duration, es.DurationTotal, goalNowPct)

	totalDays := es.EndTime.Sub(es.StartTime).Hours() / 24

	progress := &Progress{
		Date:       es.Date,
		Days:       fmt.Sprint(es.Duration) + "/" + fmt.Sprint(totalDays),
		DaysPassed: convert.FloatTo(es.Duration / es.DurationTotal * 100).Decimal(0),
		Task:       taskMap,
		Parent:     parentMap,
		Stuff:      stuffMap,
		DailyLimit: dailyLimit,
		GoalLefts:  goalLefts,
		Primary:    convert.FloatTo(es.Primary).Decimal(2),
	}

	return progress
}

func makeMapValueFloatToString(mp map[string]*essential.FieldInfo, dur float64, goalNowPct float64) (mpO map[string]*summaryField) {
	mpO = make(map[string]*summaryField)
	for k, v := range mp {
		if k == cons.Running {
			v.Done *= 5
		}
		mpO[k] = &summaryField{
			Done:       v.Done,
			DonePerDay: convert.FloatTo(v.Done / dur).Decimal(2),
			Percent:    convert.FloatTo(v.Percent).Decimal(0),
			PercentAbs: convert.FloatTo(v.PercentAbs / goalNowPct).Decimal(0),
		}
		if v.Percent >= goalNowPct*100 {
			mpO[k].AheadLevel = 1
		}
	}
	return
}

func makeStuffMapValueFloatToString(mp map[string]*essential.FieldInfo, dur, durTotal float64, goalNowPct float64) (mpO map[string]*summaryField) {
	mpO = make(map[string]*summaryField)

	goalMin, restrainMax := goalNowPct*100, 10000/(goalNowPct*100)

	for k, v := range mp {
		var pct, pctAbs, goalCutTotal, goalCutPresent, left float64
		if v.TenGoal > 0 {
			goalCutTotal = convert.FloatTo(v.TenGoal / 10.0 * durTotal).Decimal(0)
			if !cons.IsRestrainStuff(k) {
				goalCutPresent = convert.FloatTo(v.TenGoal / 10.0 * durTotal * goalNowPct).Decimal(0)
			} else {
				goalCutPresent = convert.FloatTo(v.TenGoal / 10.0 * durTotal / goalNowPct).Decimal(0)
			}

			left = goalCutPresent - v.Done
			if goalCutTotal == 0 {
				pct = 0
			} else {
				pct = convert.FloatTo(v.Done / goalCutTotal * 100 * durTotal / dur).Decimal(0)
			}

			if goalCutPresent == 0 {
				pctAbs = 100
			} else {
				pctAbs = convert.FloatTo(v.Done / goalCutPresent * 100).Decimal(0)
			}
		}

		mpO[k] = &summaryField{
			Done:       v.Done,
			DonePerDay: convert.FloatTo(v.Done / dur).Decimal(2),
			Percent:    convert.FloatTo(pct).Decimal(0),
			PercentAbs: pctAbs,
			AheadLevel: 0,
		}

		if cons.IsRestrainStuff(k) {
			if goalCutTotal == 0 {
				mpO[k].AheadLevel = 1
			} else {
				if pct <= restrainMax || (left >= 0 && durTotal-dur == 0) {
					mpO[k].AheadLevel = 1
				}
			}
		} else {
			if goalCutTotal == 0 {
				mpO[k].AheadLevel = 1
			} else {
				if left <= 0 || pct >= goalMin {
					mpO[k].AheadLevel = 1
				}
			}
		}
	}

	return

}

func makeGoalLeft(es *essential.Essential, pct float64) (mrc *goalLeft, dailyLimitField *summaryField) {
	taskInfo := make(map[string]*goalLeftField)
	parentInfo := make(map[string]*goalLeftField)
	stuffInfo := make(map[string]*goalLeftField)

	totalDays := es.EndTime.Sub(es.StartTime).Hours() / 24
	//leftDays := totalDays - es.Duration/24

	var dailyDoneLimit,
		dailyTaskGoalTotal,
		dailyLimitPct,
		dailyTaskGoalPresent float64

	//task info
	for task, fieldInfo := range es.TaskInfo {

		if cons.IsDailyTask(task) {
			taskDone := fieldInfo.Done
			if fieldInfo.Done > fieldInfo.TenGoal/10.0*es.DurationTotal*pct {
				taskDone = fieldInfo.TenGoal / 10.0 * es.DurationTotal * pct
			}
			dailyDoneLimit += taskDone
			dailyTaskGoalTotal += fieldInfo.TenGoal / 10.0 * es.DurationTotal
			dailyTaskGoalPresent += fieldInfo.TenGoal / 10.0 * es.Duration
		}

		monthGoal := fieldInfo.TenGoal / 10.0 * totalDays

		left := monthGoal*pct - fieldInfo.Done
		leftPerDay := 0.0
		if totalDays-es.Duration != 0 {
			leftPerDay = left / (totalDays - es.Duration)
		}
		if task == cons.Running {
			left *= 5
			leftPerDay *= 5
		}

		taskInfo[task] = &goalLeftField{
			Left:        fmt.Sprintf("%.2f/%.2f", leftPerDay, left),
			LeftF:       left,
			LeftFPerDay: leftPerDay,
			Finish:      0,
		}
		if left <= 0 {
			taskInfo[task] = &goalLeftField{
				Left:        "0.00/0",
				LeftF:       0,
				LeftFPerDay: 0,
				Finish:      1,
			}
		}
	}

	//code ,code input
	dailyTaskGoalTotal += es.ParentInfo[cons.Code].TenGoal / 10.0 * es.DurationTotal
	dailyTaskGoalPresent += es.ParentInfo[cons.Code].TenGoal / 10.0 * es.Duration
	croaTenGoal := es.ParentInfo[cons.Code].TenGoal - es.TaskInfo[cons.CodeInput].TenGoal
	croaDone := es.TaskInfo[cons.CodeOutput].Done + es.TaskInfo[cons.CodeInfoAndArrange].Done
	if croaDone > (croaTenGoal/10.0)*es.DurationTotal*pct {
		croaDone = (croaTenGoal / 10.0) * es.DurationTotal * pct
	}
	crDone := croaDone + es.TaskInfo[cons.CodeInput].Done
	if crDone > es.ParentInfo[cons.Code].TenGoal/10.0*es.DurationTotal*pct {
		dailyDoneLimit += es.ParentInfo[cons.Code].TenGoal / 10.0 * es.DurationTotal * pct
	} else {
		dailyDoneLimit += crDone
	}
	dailyLimitLeft := dailyTaskGoalTotal*pct - dailyDoneLimit

	if dailyTaskGoalPresent == 0 {
		dailyLimitPct = 0
	} else {
		dailyLimitPct = dailyDoneLimit / dailyTaskGoalPresent
	}

	dailyLimitLeftPerDay := 0.0
	if totalDays-es.Duration != 0 {
		dailyLimitLeftPerDay = dailyLimitLeft / (totalDays - es.Duration)
	}

	dailyLimitGoalLeft := &goalLeftField{
		Left:        fmt.Sprintf("%.2f/%.2f", dailyLimitLeftPerDay, dailyLimitLeft),
		LeftF:       dailyLimitPct,
		LeftFPerDay: dailyLimitLeftPerDay,
		Finish:      0,
	}
	if dailyLimitLeft <= 0 {
		dailyLimitGoalLeft.Finish = 1
	}

	dailyLimitPercent := convert.FloatTo(dailyLimitPct * 100).Decimal(0)

	abs := convert.FloatTo(dailyLimitPercent * (es.Duration / es.DurationTotal) / pct).Decimal(0)
	dailyLimitField = &summaryField{
		Done:       convert.FloatTo(dailyDoneLimit).Decimal(2),
		DonePerDay: convert.FloatTo(dailyDoneLimit / es.Duration).Decimal(2),
		Percent:    dailyLimitPercent,
		PercentAbs: abs,
		AheadLevel: 0,
	}

	if dailyLimitPercent >= pct*100 {
		dailyLimitField.AheadLevel = 1
	}

	//parent info
	for parent, fieldInfo := range es.ParentInfo {
		monthGoal := fieldInfo.TenGoal / 10.0 * totalDays
		left := monthGoal*pct - fieldInfo.Done
		leftPerDay := 0.0
		if totalDays-es.Duration != 0 {
			leftPerDay = left / (totalDays - es.Duration)
		}
		parentInfo[parent] = &goalLeftField{
			Left:        fmt.Sprintf("%.2f/%.2f", leftPerDay, left),
			LeftF:       left,
			LeftFPerDay: leftPerDay,
			Finish:      0,
		}
		if left <= 0 {
			parentInfo[parent] = &goalLeftField{
				Left:        "0.00/0",
				LeftF:       0,
				LeftFPerDay: 0,
				Finish:      1,
			}
		}
	}

	//stuff info
	for stuff, fieldInfo := range es.StuffInfo {
		var pctStuff float64
		if cons.IsRestrainStuff(stuff) {
			pctStuff = 1 / pct
		} else {
			pctStuff = pct
		}

		goal := fieldInfo.TenGoal / 10.0 * totalDays
		goalNow := goal * pctStuff
		goalNowCut, _ := strconv.ParseFloat(fmt.Sprintf("%.0f", goalNow), 0)

		left := goalNowCut - fieldInfo.Done

		leftPerDay := 0.0
		if totalDays-es.Duration != 0 {
			leftPerDay = left / (totalDays - es.Duration)
		}

		if !cons.IsRestrainStuff(stuff) {
			stuffInfo[stuff] = &goalLeftField{
				Left:        fmt.Sprintf("%.2f/%.2f", leftPerDay, left),
				LeftF:       left,
				LeftFPerDay: leftPerDay,
				Finish:      0,
			}
			if left <= 0 {
				stuffInfo[stuff] = &goalLeftField{
					Left:        "0.00/0",
					LeftF:       0,
					LeftFPerDay: 0,
					Finish:      1,
				}
			}
		} else { //restrain case
			if left >= 0 {
				stuffInfo[stuff] = &goalLeftField{
					Left:        fmt.Sprintf("%.2f/%.0f", leftPerDay, left),
					LeftF:       left,
					LeftFPerDay: leftPerDay,
					Finish:      1,
				}
			} else {
				stuffInfo[stuff] = &goalLeftField{
					Left:        "0.00/0",
					LeftF:       0,
					LeftFPerDay: 0,
					Finish:      0,
				}
			}
		}
	}

	left := &goalLeft{
		Name:       fmt.Sprintf("%.0f", pct*100) + "%",
		TaskInfo:   taskInfo,
		ParentInfo: parentInfo,
		StuffInfo:  stuffInfo,
		DailyLimit: dailyLimitGoalLeft,
	}

	return left, dailyLimitField
}
