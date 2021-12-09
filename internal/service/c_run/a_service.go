package c_run

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
)

type RunService struct {
	ctx   context.Context
	cache *model.Cache
	db    *gorm.DB
}

func NewRunService(c context.Context, db *gorm.DB, cache *redis.Pool) RunService {
	return RunService{
		ctx:   c,
		cache: model.NewCache(cache),
		db:    db,
	}
}

func (svc RunService) ListRuns() ([]*resp.RunResp, *resp.RunSum, error) {
	runs, err := model.NewRunModel(svc.db).GetAll()
	if err != nil {
		return nil, nil, err
	}

	runResponses := make([]*resp.RunResp, 0)
	for _, run := range runs {
		runResp := toRunResponse(run)
		runResponses = append(runResponses, runResp)
	}

	return runResponses, toRunSum(runs), nil
}

func (svc RunService) ListRunTimeZone(typ mtime.TimeType) ([]*resp.RunSum, *resp.RunSum, error) {

	numPresent := mtime.NewMTime(cons.Newest).TimeZoneNum(typ)
	zoneRuns := make([]*resp.RunSum, 0)
	for i := 1; i <= numPresent; i++ {
		zone := mtime.NewTimeZone(typ, 2021, i)
		runs, err := model.NewRunModel(svc.db).Timezone(zone)
		if err != nil {
			return nil, nil, err
		}
		runZone := toRunSum(runs)
		runZone.Date = zone.DateString()
		runZone.Id = zone.Num
		zoneRuns = append(zoneRuns, runZone)
	}

	runsYear, err := model.NewRunModel(svc.db).GetAll()
	if err != nil {
		return nil, nil, err
	}
	sum := toRunSum(runsYear)

	return zoneRuns, sum, nil
}

func toRunResponse(r *model.Run) *resp.RunResp {
	distance := convert.FloatTo(float64(r.Distance) / 1000).Decimal(2)
	TimeCostMinute := convert.FloatTo(float64(r.TimeCost) / 60).Decimal(2)
	var speed float64
	if r.TimeCost > 0 {
		speed = convert.FloatTo((float64(r.Distance) / 1000) / (float64(r.TimeCost) / 3600)).Decimal(2)
	} else {
		speed = 0
	}

	return &resp.RunResp{
		Id:          r.Id,
		Date:        r.Date,
		Distance:    distance,
		TimeCost:    TimeCostMinute,
		Pace:        fmt.Sprintf(`%v'%v"`, r.Pace/60, r.Pace%60),
		Speed:       speed,
		Rate:        r.Rate,
		Temperature: r.Temperature,
		Altitude:    r.Altitude,
	}

}

func toRunSum(runs model.Runs) *resp.RunSum {

	var totalDistance, totalDistanceWithRate, totalTimeCost int
	var totalRateMulDistance, totalTemperature, totalAltitude int
	var countTemp, countRate, countAlt int

	for _, run := range runs {
		totalDistance += run.Distance
		totalTimeCost += run.TimeCost

		if run.Temperature > model.RunMinTemperature {
			totalTemperature += run.Temperature
			countTemp++
		}

		if run.Altitude > model.RunMinAltitude {
			totalAltitude += run.Altitude
			countAlt++
		}

		if run.Rate > model.RunMinRate {
			totalRateMulDistance += run.Rate * run.Distance
			totalDistanceWithRate += run.Distance
			countRate++
		}
	}
	var rateAvg, paceAvg, tempAvg, altAvg, distanceAvg, speed = 60, 0, -273, -450, 0.0, 0.0
	if totalDistanceWithRate > 0 {
		rateAvg = totalRateMulDistance / totalDistanceWithRate
	}
	if totalDistance > 0 {
		paceAvg = int(float64(totalTimeCost*1000) / float64(totalDistance))
	}
	if countTemp > 0 {
		tempAvg = totalTemperature / countTemp
	}
	if countAlt > 0 {
		altAvg = totalAltitude / countAlt
	}
	if len(runs) > 0 {
		distanceAvg = convert.FloatTo(float64(totalDistance) / 1000 / float64(len(runs))).Decimal(2)
	}
	if totalTimeCost > 0 {
		speed = convert.FloatTo((float64(totalDistance) / 1000) / (float64(totalTimeCost) / 3600)).Decimal(2)
	}

	return &resp.RunSum{
		Times:              len(runs),
		Distance:           float64(totalDistance) / 1000,
		DistanceAverage:    distanceAvg,
		Pace:               fmt.Sprintf(`%v'%v"`, paceAvg/60, paceAvg%60),
		Speed:              speed,
		RateAverage:        rateAvg,
		TemperatureAverage: tempAvg,
		AltitudeAverage:    altAvg,
	}

}
