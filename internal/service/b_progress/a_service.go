package b_progress

import (
	"context"
	"fmt"
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/service/b_essential"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"time"
)

const (
	RedisTenProgressKey     = "WayneDataProgressTen"
	RedisMonthProgressKey   = "WayneDataProgressMonth"
	RedisQuarterProgressKey = "WayneDataProgressQuarter"
	RedisYearProgressKey    = "WayneDataProgressYear"
)

type ProgressService struct {
	ctx     context.Context
	cache   *model.Cache
	eventDb *model.EventModel
}

func NewProgressService(c context.Context) ProgressService {
	return ProgressService{
		ctx:     c,
		cache:   model.NewCache(global.CacheEngine),
		eventDb: model.NewEventModel(global.DBEngine),
	}
}

func (svc *ProgressService) GetProgress(zone *mtime.TimeZone, progressStart time.Time) (*Progress, error) {
	progress := new(Progress)

	key := progressKey(zone)
	exists, err := svc.cache.Get(key, progress)
	if err != nil {
		return nil, err
	}

	if !exists {
		events, err := svc.eventDb.Timezone(zone)
		if err != nil {
			return nil, err
		}
		es, err := b_essential.MakeEssential(events, progressStart, zone)
		if err != nil {
			return nil, err
		}
		progress = makeProgress(es, progressStart)
		progress.GcRunning, err = svc.GetYearGcRunning(2021, progressStart)
		if err != nil {
			return nil, err
		}
		err = svc.cache.Set(key, progress, 0)
		if err != nil {
			return nil, err
		}
	}

	return progress, nil
}

func progressKey(zone *mtime.TimeZone) (key string) {
	key = "WayneDataProgress" + fmt.Sprint(zone.Year) + zone.DateString()
	return
}

func (svc *ProgressService) GetYearGcRunning(year int, progressStart time.Time) (*GcRunning, error) {
	zone := mtime.NewTimeZone(mtime.TypeYear, year, 1)
	events, err := svc.eventDb.Timezone(zone)
	if err != nil {
		return nil, err
	}
	es, err := b_essential.MakeEssential(events, progressStart, zone)
	if err != nil {
		return nil, err
	}

	return gcRunningInfo(es), nil
}
