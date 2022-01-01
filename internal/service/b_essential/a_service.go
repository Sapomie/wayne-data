package b_essential

import (
	"context"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
)

type EssentialService struct {
	ctx   context.Context
	cache *model.Cache
	db    *gorm.DB
}

func NewEssentialService(c context.Context, db *gorm.DB, cache *redis.Pool) EssentialService {
	return EssentialService{
		ctx:   c,
		cache: model.NewCache(cache),
		db:    db,
	}
}

func (svc *EssentialService) GetEssentialList(typ mtime.TimeType) (Essentials, int, error) {

	key := getEssentialRedisKey(typ)

	var ess Essentials
	exists, err := svc.cache.Get(key, &ess)
	if err != nil {
		return nil, 0, err
	}
	if !exists {
		ess, err = svc.getEssentialsFromDB(typ)
		if err != nil {
			return nil, 0, err
		}
		err = svc.cache.Set(key, ess, 0)
		if err != nil {
			return nil, 0, err
		}
	}

	ess.Response()
	return ess, len(ess), nil
}

func getEssentialRedisKey(typ mtime.TimeType) string {
	var key string
	switch typ {
	case mtime.TypeDay:
		key = cons.RedisKeyDayEss
	case mtime.TypeTen:
		key = cons.RedisKeyTenEss
	case mtime.TypeMonth:
		key = cons.RedisKeyMonthEss
	case mtime.TypeQuarter:
		key = cons.RedisKeyQuarterEss
	case mtime.TypeHalf:
		key = cons.RedisKeyHalfEss
	case mtime.TypeYear:
		key = cons.RedisKeyYearEss
	}
	return key
}

func (svc *EssentialService) getEssentialsFromDB(typ mtime.TimeType) (Essentials, error) {
	start, _ := mtime.NewTimeZone(mtime.TypeYear, 2022, 1).BeginAndEnd()
	if typ == mtime.TypeYear {
		start = cons.DbOldest
	}
	events, _, err := model.NewEventModel(svc.db).GetAll()
	ess, err := MakeEssentials(events, start, cons.DbNewest, typ)
	if err != nil {
		return nil, err
	}
	return ess, nil
}
