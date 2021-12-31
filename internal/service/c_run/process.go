package c_run

import (
	"errors"
	"fmt"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"strconv"
	"strings"
	"time"
)

func (svc RunService) ProcessRun() ([]string, error) {
	runs, infos, err := svc.makeRuns()
	if err != nil {
		return nil, err
	}

	err = svc.storeRuns(runs)
	if err != nil {
		return nil, err
	}

	return infos, nil
}

func (svc RunService) makeRuns() (runs model.Runs, infos []string, err error) {

	events, err := model.NewEventModel(svc.db).ByTaskName(cons.DbOldest, cons.DbNewest, cons.Running)
	if err != nil {
		return nil, nil, err
	}

	for _, event := range events {
		run, err := makeRun(event)
		if err != nil {
			info := fmt.Sprintf("make run error,event start time: %v,coment: %v", event.Start(), event.Comment)
			infos = append(infos, info)
			continue
		}
		runs = append(runs, run)
	}

	return
}

func makeRun(event *model.Event) (*model.Run, error) {
	fieldsRaw := strings.Split(event.Comment, "，")

	var fields []string
	for _, field := range fieldsRaw {
		fields = append(fields, strings.TrimSpace(field))
	}

	if len(fields) < 3 {
		return nil, errors.New("not enough fields")
	}

	distance, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, errors.New("run distance parse error")
	}
	if distance < 0 || distance > 100000 {
		return nil, errors.New("run distance out of range")
	}

	minutes, err := strconv.Atoi(fields[1])
	if err != nil {
		return nil, errors.New("run minutes parse error")
	}
	if minutes < 0 || minutes > 1200 {
		return nil, errors.New("run minutes out of range")
	}

	seconds, err := strconv.Atoi(fields[2])
	if err != nil {
		return nil, errors.New("run seconds parse error")
	}
	if seconds < 0 || seconds >= 60 {
		return nil, errors.New("run seconds out of range")
	}

	altitude := -model.RunMinAltitude
	if len(fields) >= 4 {
		altitude, err = strconv.Atoi(fields[3])
		if err != nil {
			return nil, errors.New("run altitude parse error")
		}
		if altitude < -415 || altitude >= 5000 {
			return nil, errors.New("run altitude out of range")
		}
	}

	temperature := -model.RunMinTemperature
	if len(fields) >= 5 {
		temperature, err = strconv.Atoi(fields[4])
		if err != nil {
			return nil, errors.New("run temperature parse error")
		}
		if temperature < -273 || temperature >= 50 {
			return nil, errors.New("run temperature out of range")
		}
	}

	rate := model.RunMinRate
	if len(fields) >= 6 {
		rate, err = strconv.Atoi(fields[5])
		if err != nil {
			return nil, errors.New("run rate parse error")
		}
		if rate < 100 || rate >= 200 {
			return nil, errors.New("run rate out of range")
		}
	}

	comment := ""
	if len(fields) >= 7 {
		comment = fields[6]
	}

	timeCo := minutes*60 + seconds
	run := &model.Run{
		StartTime:   event.StartTime,
		Date:        time.Unix(event.StartTime, 0).Format(mtime.TimeTemplate5),
		Distance:    distance,
		TimeCost:    timeCo,
		Pace:        timeCo * 1000 / distance,
		Rate:        rate,
		Temperature: temperature,
		Altitude:    altitude,
		Comment:     comment,
	}

	return run, nil

}

func (svc RunService) storeRuns(runs model.Runs) error {
	for _, run := range runs {
		mm := model.NewRunModel(svc.db)
		exist, err := mm.Exists(run.StartTime)
		if err != nil {
			return err
		}
		if exist {
			err := mm.Base.Where("start_time = ?", run.StartTime).Update(run).Error
			if err != nil {
				return err
			}
		} else {
			err := mm.Base.Create(run).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}
