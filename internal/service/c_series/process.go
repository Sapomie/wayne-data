package c_series

import (
	"errors"
	"fmt"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"strconv"
	"strings"
)

func (svc SeriesService) ProcessSeries() ([]string, error) {
	seriesS, infos, err := svc.makeTvSeriesS()
	if err != nil {
		return nil, err
	}

	err = svc.storeSeriesS(seriesS)
	if err != nil {
		return nil, err
	}

	return infos, nil
}

func (svc SeriesService) storeSeriesS(seriesS model.SeriesS) error {

	for _, series := range seriesS {
		mm := model.NewSeriesModel(svc.db)

		exist, err := mm.Exists(series.NameSeason)
		if err != nil {
			return err
		}
		if exist {
			err := mm.Base.Where("name_season = ?", series.NameSeason).Update(series).Error
			if err != nil {
				return err
			}
		} else {
			err := mm.Base.Create(series).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (svc SeriesService) makeTvSeriesS() (seriesS model.SeriesS, infos []string, err error) {
	start, end := mtime.NewTimeZone(mtime.TypeYear, 2021, 1).BeginAndEnd()
	events, err := model.NewEventModel(svc.db).ByTaskName(start, end, cons.AnimationAndEpisode)
	if err != nil {
		return nil, nil, err
	}

	seriesMap := make(map[string]model.Events, 0)

	for _, event := range events {
		strs := strings.Split(event.Comment, "，")
		name := strs[0]
		seriesMap[name] = append(seriesMap[name], event)
	}

	for name, seriesEvents := range seriesMap {
		series := &model.Series{Name: name}
		for _, event := range seriesEvents {

			if isSeriesFirstTime(event) {
				series.NameOrigin, series.Category, series.Year, series.Season, series.EpisodeNumber, err = seriesInfo(event)
				series.NameSeason = series.Name + "_" + fmt.Sprintf("第%v季", series.Season)
				series.FirstTime = event.StartTime
				if err != nil {
					infos = append(infos, fmt.Sprintf("make series error,event start time: %v,coment: %v", event.Start(), event.Comment))
					continue
				}
			}

			if isSeriesLastTime(event) {
				series.Rate, err = seriesRate(event)
				if err != nil {
					infos = append(infos, fmt.Sprintf("make series error,event start time: %v,coment: %v", event.Start(), event.Comment))
					continue
				}
				series.Finish = model.ProjectFinish
			}

			if event.StartTime > series.LastTime {
				series.LastTime = event.StartTime
			}

			series.Duration += event.Duration
		}
		seriesS = append(seriesS, series)
	}
	return
}

func seriesInfo(event *model.Event) (originName, category string, year, season, episodeNumber int, err error) {
	strs := strings.Split(event.Comment, "，")
	if len(strs) < 7 {
		return "", "", 0, 0, 0, errors.New("wrong length of series comment")
	}
	category = strs[3]
	originName = strs[2]
	season, err = strconv.Atoi(strs[1])
	if err != nil {
		return "", "", 0, 0, 0, err
	}
	year, err = strconv.Atoi(strs[4])
	if err != nil {
		return "", "", 0, 0, 0, err
	}
	episodeNumber, err = strconv.Atoi(strs[5])
	if err != nil {
		return "", "", 0, 0, 0, err
	}

	return
}

func isSeriesFirstTime(event *model.Event) bool {
	if strings.Contains(event.Comment, "、s") {
		return true
	}
	return false
}

func isSeriesLastTime(event *model.Event) bool {
	if strings.Contains(event.Comment, "、e") {
		return true
	}
	return false
}

func seriesRate(event *model.Event) (rate int, err error) {
	strs := strings.Split(event.Comment, "，")
	if len(strs) < 4 {
		return 0, errors.New("wrong length of series comment")
	}
	rate, err = strconv.Atoi(strs[2])
	if err != nil {
		return 0, err
	}
	return
}
