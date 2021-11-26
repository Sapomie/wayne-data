package b_essential

import (
	"context"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
)

const (
	RedisDayEssKey     = "WayneDataEssDay"
	RedisTenEssKey     = "WayneDataEssTen"
	RedisMonthEssKey   = "WayneDataEssMonth"
	RedisQuarterEssKey = "WayneDataEssQuarter"
	RedisHalfEssKey    = "WayneDataEssHalf"
	RedisYearEssKey    = "WayneDataEssYear"
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
	var ess Essentials

	var key string
	switch typ {
	case mtime.TypeDay:
		key = RedisDayEssKey
	case mtime.TypeTen:
		key = RedisTenEssKey
	case mtime.TypeMonth:
		key = RedisMonthEssKey
	case mtime.TypeQuarter:
		key = RedisQuarterEssKey
	case mtime.TypeHalf:
		key = RedisHalfEssKey
	case mtime.TypeYear:
		key = RedisYearEssKey
	}

	exists, err := svc.cache.Get(key, &ess)
	if err != nil {
		return nil, 0, err
	}
	if !exists {
		start, _ := mtime.NewTimeZone(mtime.TypeYear, 2021, 1).BeginAndEnd()
		events, _, err := model.NewEventModel(svc.db).GetAll()
		ess, err = MakeEssentials(events, start, cons.Newest, typ)
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
