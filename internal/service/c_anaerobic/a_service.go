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

func (svc AnaerobicService) ListAnaerobic() (*resp.Anaerobic, error) {
	anaerobic := new(resp.Anaerobic)
	key := cons.RedisKeyAnaerobic
	exists, err := svc.cache.Get(key, &anaerobic)
	if err != nil {
		return nil, err
	}
	if !exists {
		anaerobic, err = svc.GetAnaerobicFromDB()
		if err != nil {
			return nil, err
		}
		err = svc.cache.Set(key, anaerobic, 0)
		if err != nil {
			return nil, err
		}
	}

	return anaerobic, nil
}

func (svc AnaerobicService) ListAnaerobicZone(typ mtime.TimeType) (*resp.AnaerobicZone, error) {
	anaerobicZone := new(resp.AnaerobicZone)
	key := getAnaerobicZoneKey(typ)
	exists, err := svc.cache.Get(key, &anaerobicZone)
	if err != nil {
		return nil, err
	}
	if !exists {
		anaerobicZone, err = svc.GerAnaerobicZoneFromDB(typ)
		if err != nil {
			return nil, err
		}
		err = svc.cache.Set(key, anaerobicZone, 0)
		if err != nil {
			return nil, err
		}
	}

	return anaerobicZone, nil
}

func (svc AnaerobicService) GetAnaerobicFromDB() (*resp.Anaerobic, error) {
	anaerobicS, err := model.NewAnaerobicModel(svc.db).GetAll()
	if err != nil {
		return nil, err
	}

	anaerobicResponses := make([]*resp.AnaerobicItem, 0)
	for _, anaerobic := range anaerobicS {
		anaerobicResp := toAnaerobicResponse(anaerobic)
		anaerobicResponses = append(anaerobicResponses, anaerobicResp)
	}
	sum := toAnaerobicSum(anaerobicS)
	sum.Protein1, err = svc.getProtein(mtime.NewTimeZone(mtime.TypeYear, 2022, 1).BeginAndEnd())
	if err != nil {
		return nil, err
	}

	return &resp.Anaerobic{
		Items: anaerobicResponses,
		Sum:   sum,
	}, nil
}

func (svc AnaerobicService) GerAnaerobicZoneFromDB(typ mtime.TimeType) (*resp.AnaerobicZone, error) {

	numPresent := mtime.NewMTime(cons.DbNewest).TimeZoneNum(typ)
	zoneRuns := make([]*resp.AnaerobicSum, 0)
	for i := 1; i <= numPresent; i++ {
		zone := mtime.NewTimeZone(typ, 2022, i)
		as, err := model.NewAnaerobicModel(svc.db).Timezone(zone)
		if err != nil {
			return nil, err
		}
		asZone := toAnaerobicSum(as)
		asZone.Date = zone.DateString()
		asZone.Id = zone.Num
		asZone.Protein1, err = svc.getProtein(zone.BeginAndEnd())
		zoneRuns = append(zoneRuns, asZone)
	}

	anaerobicYear, err := model.NewAnaerobicModel(svc.db).GetAll()
	if err != nil {
		return nil, err
	}
	sum := toAnaerobicSum(anaerobicYear)
	sum.Protein1, err = svc.getProtein(mtime.NewTimeZone(mtime.TypeYear, 2022, 1).BeginAndEnd())
	if err != nil {
		return nil, err
	}

	return &resp.AnaerobicZone{
		Items: zoneRuns,
		Sum:   sum,
	}, nil
}

func toAnaerobicResponse(a *model.Anaerobic) *resp.AnaerobicItem {

	return &resp.AnaerobicItem{
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

func getAnaerobicZoneKey(typ mtime.TimeType) string {
	str := ""
	switch typ {
	case mtime.TypeTen:
		str = "Ten"
	case mtime.TypeMonth:
		str = "Month"
	case mtime.TypeQuarter:
		str = "Quarter"
	case mtime.TypeYear:
		str = "Year"
	}

	return cons.RedisKeyAnaerobicZonePrefix + str

}
