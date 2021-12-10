package c_anaerobic

import (
	"context"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"github.com/Sapomie/wayne-data/internal/model/resp"
	"github.com/Sapomie/wayne-data/pkg/convert"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	"time"
)

type AnaerobicService struct {
	ctx   context.Context
	cache *model.Cache
	db    *gorm.DB
}

func NewAnaerobicService(c context.Context, db *gorm.DB, cache *redis.Pool) AnaerobicService {
	return AnaerobicService{
		ctx:   c,
		cache: model.NewCache(cache),
		db:    db,
	}
}

func (svc AnaerobicService) ListAnaerobicS() ([]*resp.AnaerobicResp, *resp.AnaerobicSum, error) {
	anaerobicS, err := model.NewAnaerobicModel(svc.db).GetAll()
	if err != nil {
		return nil, nil, err
	}

	anaerobicResponses := make([]*resp.AnaerobicResp, 0)
	for _, anaerobic := range anaerobicS {
		anaerobicResp := toAnaerobicResponse(anaerobic)
		anaerobicResponses = append(anaerobicResponses, anaerobicResp)
	}
	sum := toAnaerobicSum(anaerobicS)
	sum.Protein1, err = svc.getProtein(mtime.NewTimeZone(mtime.TypeYear, 2021, 1).BeginAndEnd())
	if err != nil {
		return nil, nil, err
	}

	return anaerobicResponses, sum, nil
}

func (svc AnaerobicService) ListAnaerobicTimeZone(typ mtime.TimeType) ([]*resp.AnaerobicSum, *resp.AnaerobicSum, error) {

	numPresent := mtime.NewMTime(cons.Newest).TimeZoneNum(typ)
	zoneRuns := make([]*resp.AnaerobicSum, 0)
	for i := 1; i <= numPresent; i++ {
		zone := mtime.NewTimeZone(typ, 2021, i)
		as, err := model.NewAnaerobicModel(svc.db).Timezone(zone)
		if err != nil {
			return nil, nil, err
		}
		asZone := toAnaerobicSum(as)
		asZone.Date = zone.DateString()
		asZone.Id = zone.Num
		asZone.Protein1, err = svc.getProtein(zone.BeginAndEnd())
		zoneRuns = append(zoneRuns, asZone)
	}

	anaerobicYear, err := model.NewAnaerobicModel(svc.db).GetAll()
	if err != nil {
		return nil, nil, err
	}
	sum := toAnaerobicSum(anaerobicYear)
	sum.Protein1, err = svc.getProtein(mtime.NewTimeZone(mtime.TypeYear, 2021, 1).BeginAndEnd())
	if err != nil {
		return nil, nil, err
	}

	return zoneRuns, sum, nil
}

func toAnaerobicResponse(a *model.Anaerobic) *resp.AnaerobicResp {

	return &resp.AnaerobicResp{
		Date:     time.Unix(a.StartTime, 0).Format(mtime.TimeTemplate5),
		Name:     a.Name,
		Group:    a.Group,
		Times:    a.Times,
		Addition: a.Addition,
	}

}

func (svc AnaerobicService) getProtein(start, end time.Time) (int, error) {

	proteinEvents, err := model.NewEventModel(svc.db).ByStuffName(start, end, cons.StuProtein1)
	if err != nil {
		return 0, err
	}
	return len(proteinEvents), nil

}

func toAnaerobicSum(as model.AnaerobicS) *resp.AnaerobicSum {

	var (
		totalGroup              int
		sitUpGroup              int
		pushUpGroup             int
		dumbbellPressGroup      int
		sitUpTimesTotal         int
		pushUpTimesTotal        int
		dumbbellPressTimesTotal int
		dumbbellPressMassTotal  float64
	)

	for _, ab := range as {
		totalGroup += ab.Group
		switch ab.Name {
		case model.SitUp:
			sitUpGroup += ab.Group
			sitUpTimesTotal += ab.Times * ab.Group
		case model.PushUp:
			pushUpGroup += ab.Group
			pushUpTimesTotal += ab.Times * ab.Group
		case model.DumbbellPress:
			dumbbellPressGroup += ab.Group
			dumbbellPressTimesTotal += ab.Times * ab.Group
			dumbbellPressMassTotal += ab.Addition * float64(ab.Group)
		}
	}

	resp := &resp.AnaerobicSum{
		TotalGroup:          totalGroup,
		SitUpGroups:         sitUpGroup,
		PushUpGroups:        pushUpGroup,
		DumbbellPressGroups: dumbbellPressGroup,
	}
	if sitUpGroup > 0 {
		resp.SitUpPerGroup = sitUpTimesTotal / sitUpGroup
	}
	if pushUpGroup > 0 {
		resp.PushUpPerGroup = pushUpTimesTotal / pushUpGroup
	}
	if dumbbellPressGroup > 0 {
		resp.DumbbellPressPerGroup = dumbbellPressTimesTotal / dumbbellPressGroup
		resp.DumbbellPressMass = convert.FloatTo((dumbbellPressMassTotal / float64(dumbbellPressGroup))).Decimal(1)
	}

	return resp
}
