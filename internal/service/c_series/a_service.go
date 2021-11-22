package c_series

import (
	"context"
	"errors"
	"fmt"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

type SeriesService struct {
	ctx      context.Context
	cache    *model.Cache
	seriesDb *model.SeriesModel
	eventDb  *model.EventDbModel
}

func NewSeriesService(c context.Context, db *gorm.DB, cache *redis.Pool) SeriesService {
	return SeriesService{
		ctx:      c,
		cache:    model.NewCache(cache),
		seriesDb: model.NewSeriesModel(db),
		eventDb:  model.NewEventModel(db),
	}
}

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
	mm := svc.seriesDb

	for _, series := range seriesS {
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
	events, err := svc.eventDb.ByTaskName(start, end, cons.AnimationAndEpisode)
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
				series.FirstReadingTime = event.StartTime
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
				series.Finish = model.BookSeriesFinish
			}

			if event.StartTime > series.LastReadingTime {
				series.LastReadingTime = event.StartTime
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
