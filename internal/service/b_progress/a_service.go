package b_progress

import (
	"context"
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/service/b_essential"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"time"
)

type ProgressService struct {
	ctx     context.Context
	cache   *model.Cache
	eventDb *model.EventDbModel
}

func NewEssentialService(c context.Context) ProgressService {
	return ProgressService{
		ctx:     c,
		cache:   model.NewCache(global.CacheEngine),
		eventDb: model.NewEventModel(global.DBEngine),
	}
}

func (svc *ProgressService) GetProgress(zone *mtime.TimeZone, progressStart time.Time) (*Progress, error) {
	events, err := svc.eventDb.Timezone(zone)
	if err != nil {
		return nil, err
	}
	es, err := b_essential.MakeEssential(events, progressStart, zone)
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
