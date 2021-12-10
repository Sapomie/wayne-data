package c_anaerobic

import (
	"fmt"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"strconv"
	"strings"
)

func (svc AnaerobicService) ProcessAnaerobic() ([]string, error) {
	anaerobicS, infos, err := svc.makeAnaerobicS()
	if err != nil {
		return nil, err
	}

	err = svc.storeAnaerobicS(anaerobicS)
	if err != nil {
		return nil, err
	}

	return infos, nil
}

func (svc AnaerobicService) makeAnaerobicS() (anaerobicS model.AnaerobicS, infos []string, err error) {

	start, end := mtime.NewTimeZone(mtime.TypeYear, 2021, 1).BeginAndEnd()
	events, err := model.NewEventModel(svc.db).ByTaskName(start, end, cons.Anaerobic)

	for _, event := range events {
		strs := strings.Split(event.Comment, "，")
		name := strs[0]
		group, err := strconv.Atoi(strs[1])
		if err != nil {
			infos = append(infos, fmt.Sprintf("make anaerobic error,event start time: %v,coment: %v", event.Start(), event.Comment))
			continue
		}
		number, err := strconv.Atoi(strs[2])
		if err != nil {
			infos = append(infos, fmt.Sprintf("make anaerobic error,event start time: %v,coment: %v", event.Start(), event.Comment))
			continue
		}

		anaerobic := &model.Anaerobic{
			StartTime: event.StartTime,
			EndTime:   event.EndTime,
			Name:      name,
			Group:     group,
			Times:     number,
		}
		if len(strs) >= 4 {
			addition, err := strconv.ParseFloat(strs[3], 64)
			if err != nil {
				infos = append(infos, fmt.Sprintf("make book error,event start time: %v,coment: %v", event.Start(), event.Comment))
				continue
			}
			anaerobic.Addition = addition
		}
		anaerobicS = append(anaerobicS, anaerobic)
	}
	return
}

func (svc AnaerobicService) storeAnaerobicS(anaerobicS model.AnaerobicS) error {
	for _, anaerobic := range anaerobicS {
		mm := model.NewAnaerobicModel(svc.db)
		exist, err := mm.Exists(anaerobic.StartTime)
		if err != nil {
			return err
		}
		if exist {
			err := mm.Base.Where("start_time = ?", anaerobic.StartTime).Update(anaerobic).Error
			if err != nil {
				return err
			}
		} else {
			err := mm.Base.Create(anaerobic).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}
