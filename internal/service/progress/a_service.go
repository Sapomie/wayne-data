package progress

import (
	"context"
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/service/essential"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"time"
)

type ServiceProgress struct {
	ctx     context.Context
	cache   *model.Cache
	eventDb *model.EventDbModel
}

func NewEssentialService(c context.Context) ServiceProgress {
	return ServiceProgress{
		ctx:     c,
		cache:   model.NewCache(global.CacheEngine),
		eventDb: model.NewEventModel(global.DBEngine),
	}
}

func (svc *ServiceProgress) GetProgress(zone *mtime.TimeZone, progressStart time.Time) (*Progress, error) {
	events, err := svc.eventDb.Timezone(zone)
	if err != nil {
		return nil, err
	}
	es, err := essential.MakeEssential(events, progressStart, zone)
	if err != nil {
		return nil, err
	}
	progress := makeProgress(es, progressStart)

	progress.GcRunning, err = svc.GetYearGcRunning(2021, progressStart)
	if err != nil {
		return nil, err
	}
	return progress, nil
}

func (svc *ServiceProgress) GetYearGcRunning(year int, progressStart time.Time) (*GcRunning, error) {
	zone := mtime.NewTimeZone(mtime.TypeYear, year, 1)
	events, err := svc.eventDb.Timezone(zone)
	if err != nil {
		return nil, err
	}
	es, err := essential.MakeEssential(events, progressStart, zone)
	if err != nil {
		return nil, err
	}

	return gcRunningInfo(es), nil
}
