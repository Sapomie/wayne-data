package b_rawEvent

import (
	"fmt"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"strings"
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

func (svc RawEventService) makeEvents(raws []*RawEvent) (model.Events, []string, error) {
	var events model.Events
	var info []string
	for _, raw := range raws {
		event, taskAndParentAddingInfo, err := svc.makeEvent(raw)
		if err != nil {
			return nil, nil, err
		}
		events = append(events, event)
		info = append(info, taskAndParentAddingInfo...)
	}
	return events, info, nil
}

func (svc RawEventService) storeEvents(events model.Events) (infos []string, err error) {

	var (
		eventsInsert, eventsUpdate     model.Events
		countInsert, countUpdate       int
		durationInsert, durationUpdate float64
	)

	for _, event := range events {
		em := svc.eventDb
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
func (svc RawEventService) makeEvent(raw *RawEvent) (event *model.Event, info []string, err error) {

	start, end, err := raw.parseRawEventTime()
	if err != nil {
		return nil, nil, err
	}

	task, taskAddingInfo, err := svc.taskDb.InsertAndGetTask(raw.TaskName)
	if err != nil {
		return nil, nil, err
	}
	if taskAddingInfo != "" {
		info = append(info, taskAddingInfo)
	}
	parent, parentAddingInfo, err := svc.parentDb.InsertAndGetParent(raw.ParentTask)
	if err != nil {
		return nil, nil, err
	}
	if parentAddingInfo != "" {
		info = append(info, parentAddingInfo)
	}

	err = svc.translateAbbreviate(raw, task.Id)
	if err != nil {
		return nil, nil, err
	}

	stuffIds, tagIds, remark, projectId, commentPropertyInfo, err := svc.processCommentProperty(raw, task.Id)
	if err != nil {
		return nil, nil, err
	}
	if commentPropertyInfo != nil {
		info = append(info, commentPropertyInfo...)
	}

	event = &model.Event{
		Date:      start.Format(mtime.TimeTemplate5),
		Duration:  raw.Duration,
		TaskId:    task.Id,
		ParentId:  parent.Id,
		Comment:   raw.Comment,
		StartTime: start.Unix(),
		EndTime:   end.Unix(),
		StuffId:   stuffIds,
		TagId:     tagIds,
		ProjectId: projectId,
		Remark:    remark,
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

func (svc RawEventService) translateAbbreviate(raw *RawEvent, taskId int) (err error) {
	abbrs, err := svc.abbrDb.GetAll()
	if err != nil {
		return err
	}
	for _, abbr := range abbrs {
		if taskId == abbr.TaskId {
			raw.Comment = strings.ReplaceAll(raw.Comment, abbr.Abbr, abbr.Content)
		}
	}
	return nil
}

func (svc RawEventService) eventToRawEvent(event *model.Event) (*RawEvent, error) {

	task, err := svc.taskDb.ById(event.TaskId)
	if err != nil {
		return nil, err
	}
	parent, err := svc.parentDb.ById(event.ParentId)
	if err != nil {
		return nil, err
	}

	raw := &RawEvent{
		TaskName:   task.Name,
		StartDate:  time.Unix(event.StartTime, 0).Format(mtime.TimeTemplate10),
		StartTime:  time.Unix(event.StartTime, 0).Format(mtime.TimeTemplate9),
		EndDate:    time.Unix(event.EndTime, 0).Format(mtime.TimeTemplate10),
		EndTime:    time.Unix(event.EndTime, 0).Format(mtime.TimeTemplate9),
		Duration:   event.Duration,
		Comment:    event.Comment,
		ParentTask: parent.Name,
	}

	return raw, nil
}
