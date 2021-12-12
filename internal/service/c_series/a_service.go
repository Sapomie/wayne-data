package c_series

import (
	"context"
	"fmt"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"github.com/Sapomie/wayne-data/internal/model/resp"
	"github.com/Sapomie/wayne-data/pkg/convert"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	"time"
)

type SeriesService struct {
	ctx   context.Context
	cache *model.Cache
	db    *gorm.DB
}

func NewSeriesService(c context.Context, db *gorm.DB, cache *redis.Pool) SeriesService {
	return SeriesService{
		ctx:   c,
		cache: model.NewCache(cache),
		db:    db,
	}
}

func (svc SeriesService) ListSeries() (*resp.Series, error) {
	seriesS := new(resp.Series)
	key := cons.RedisKeySeries
	exists, err := svc.cache.Get(key, &seriesS)
	if err != nil {
		return nil, err
	}
	if !exists {
		seriesS, err = svc.GetSeriesFromDB()
		if err != nil {
			return nil, err
		}
		err = svc.cache.Set(key, seriesS, 0)
		if err != nil {
			return nil, err
		}
	}

	return seriesS, nil
}

func (svc SeriesService) GetSeriesFromDB() (*resp.Series, error) {
	seriesS, err := model.NewSeriesModel(svc.db).GetAll()
	if err != nil {
		return nil, err
	}

	bookResponses := make([]*resp.SeriesItem, 0)
	for _, series := range seriesS {
		bookResp := toSeriesResponse(series)
		bookResponses = append(bookResponses, bookResp)
	}

	return &resp.Series{
		Item: bookResponses,
		Sum:  toSeriesSum(seriesS),
	}, nil
}

func toSeriesResponse(s *model.Series) *resp.SeriesItem {
	var finishMark string
	switch s.Finish {
	case model.BookFinish:
		finishMark = "Finish"
	case model.BookAbandon:
		finishMark = "Abandon"
	}

	return &resp.SeriesItem{
		Name:          s.Name,
		Category:      s.Category,
		Season:        fmt.Sprintf("第%v季", s.Season),
		Year:          s.Year,
		EpisodeNumber: s.EpisodeNumber,
		Duration:      s.Duration,
		Rate:          s.Rate,
		Finish:        finishMark,
		FirstTime:     time.Unix(s.FirstTime, 0).Format(mtime.TimeTemplate4),
		LastTime:      time.Unix(s.LastTime, 0).Format(mtime.TimeTemplate4),
	}

}

func toSeriesSum(seriesS model.SeriesS) *resp.SeriesSum {
	var (
		finishNum         int
		durationFinishSum float64
		durationSum       float64
		rateSum           int
	)

	for _, series := range seriesS {
		if series.Finish == model.BookFinish {
			finishNum++
			durationFinishSum += series.Duration
			rateSum += series.Rate
		}
		durationSum += series.Duration
	}

	var durationAvg float64
	var rateAvg int
	if finishNum > 0 {
		durationAvg = durationSum / float64(finishNum)
		rateAvg = rateSum / finishNum
	}

	return &resp.SeriesSum{
		SeriesNumber: len(seriesS),
		DurationAvg:  convert.FloatTo(durationAvg).Decimal(2),
		RateAvg:      rateAvg,
		Finish:       finishNum,
	}

}
