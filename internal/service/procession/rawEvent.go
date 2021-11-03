package procession

import (
	"fmt"
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"time"
)

type RawEvent struct {
	TaskName   string  `csv:"任务名称"`
	StartDate  string  `csv:"开始日期"`
	StartTime  string  `csv:"开始时间"`
	EndDate    string  `csv:"结束日期"`
	EndTime    string  `csv:"结束时间"`
	Duration   float64 `csv:"历时（小时）"`
	Comment    string  `csv:"注解"`
	ParentTask string  `csv:"母任务"`
}

func makeEventsByRaws(raws []*RawEvent) (model.Events, []string, error) {
	var events model.Events
	var info []string
	for _, raw := range raws {
		event, taskAndParentAddingInfo, err := raw.toEvent()
		if err != nil {
			return nil, nil, err
		}
		events = append(events, event)
		info = append(info, taskAndParentAddingInfo...)
	}
	return events, info, nil
}

func storeEvents(events model.Events) (infos []string, err error) {

	var (
		eventsInsert, eventsUpdate     model.Events
		countInsert, countUpdate       int
		durationInsert, durationUpdate float64
	)

	for _, event := range events {
		em := model.NewEventModel(global.DBEngine)
		existEvt, err := em.Exists(event.StartTime)
		if err != nil {
			return nil, err
		}
		if existEvt {
			err := em.Base.Where("start_time = ?", event.StartTime).Update(event).Error
			if err != nil {
				return nil, err
			}
			eventsUpdate = append(eventsUpdate, event)
			countUpdate++
			durationUpdate += event.Duration
		} else {
			err := em.Base.Create(event).Error
			if err != nil {
				return nil, err
			}
			durationInsert += event.Duration
			eventsInsert = append(eventsInsert, event)
			countInsert++
		}
	}

	insetInfo := fmt.Sprintf("Adding %v events to database, in all %.1f days", countInsert, durationInsert/24)
	updateInfo := fmt.Sprintf("Updating %v events in database, in all %.1f days", countUpdate, durationUpdate/24)

	infos = []string{
		insetInfo,
		updateInfo,
	}

	return
}

//通过 RawEvent 生成 Event, 并且插入 task,parentTask 等条目
func (raw *RawEvent) toEvent() (event *model.Event, info []string, err error) {
	start, end, err := raw.parseRawEventTime()
	if err != nil {
		return nil, nil, err
	}
	task, taskAddingInfo, err := model.NewTaskModel(global.DBEngine).InsertAndGetTask(raw.TaskName)
	if err != nil {
		return nil, nil, err
	}
	if taskAddingInfo != "" {
		info = append(info, taskAddingInfo)
	}
	parent, parentAddingInfo, err := model.NewParentModel(global.DBEngine).InsertAndGetParent(raw.ParentTask)
	if err != nil {
		return nil, nil, err
	}
	if parentAddingInfo != "" {
		info = append(info, parentAddingInfo)
	}

	event = &model.Event{
		Date:      start.Format(mtime.TimeTemplate5),
		Duration:  raw.Duration,
		TaskId:    task.Id,
		ParentId:  parent.Id,
		Comment:   raw.Comment,
		StartTime: start.Unix(),
		EndTime:   end.Unix(),
		StuffId:   "",
		TagId:     "",
		ProjectId: 0,
		Remark:    "",
	}

	return
}

func (raw *RawEvent) parseRawEventTime() (startTime, endTime time.Time, err error) {
	startTimeStr := raw.StartDate + " " + raw.StartTime
	endTimeStr := raw.EndDate + " " + raw.EndTime

	startTime, err = parseTime(startTimeStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endTime, err = parseTime(endTimeStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return
}

func parseTime(str string) (time.Time, error) {
	timeParsed, err := time.ParseInLocation(mtime.TimeTemplate1, str, time.Local)
	if err != nil {
		timeParsed, err = time.ParseInLocation(mtime.TimeTemplate2, str, time.Local)
		if err != nil {
			timeParsed, err = time.ParseInLocation(mtime.TimeTemplate3, str, time.Local)
			if err != nil {
				return time.Time{}, err
			}
		}
	}
	return timeParsed, nil
}
