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

func (svc *ServiceProgress) GetProgress(typ mtime.TimeType, year, num int, progressStart time.Time) (*Progress, error) {
	events, err := svc.eventDb.Timezone(typ, year, num)
	if err != nil {
		return nil, err
	}
	es, err := essential.MakeEssential(events, progressStart, typ, year, num)
	if err != nil {
		return nil, err
	}

	return MakeProgress(es, progressStart), nil
}
