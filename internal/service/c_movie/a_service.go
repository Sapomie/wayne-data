package c_movie

import (
	"context"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/model/resp"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	"time"
)

type MovieService struct {
	ctx   context.Context
	cache *model.Cache
	db    *gorm.DB
}

func NewMovieService(c context.Context, db *gorm.DB, cache *redis.Pool) MovieService {
	return MovieService{
		ctx:   c,
		cache: model.NewCache(cache),
		db:    db,
	}
}

func (svc MovieService) ListMovies() ([]*resp.MovieResp, *resp.MovieSum, error) {
	movies, err := model.NewMovieModel(svc.db).GetAll()
	if err != nil {
		return nil, nil, err
	}

	movieResponses := make([]*resp.MovieResp, 0)
	for _, movie := range movies {
		movieResp := toMovieResponse(movie)
		movieResponses = append(movieResponses, movieResp)
	}

	return movieResponses, toMovieSum(movies), nil
}

func toMovieResponse(m *model.Movie) *resp.MovieResp {
	place := ""
	if m.Place == 2 {
		place = "Cinema"
	}

	return &resp.MovieResp{
		Date:   time.Unix(m.StartTime, 0).Format(mtime.TimeTemplate4),
		Name:   m.Name,
		NameEn: m.EnName,
		Rate:   m.Rate,
		Year:   m.Year,
		Place:  place,
	}

}

func toMovieSum(movies model.Movies) *resp.MovieSum {
	var rateSum, countPlace int

	for _, movie := range movies {
		rateSum += movie.Rate
		if movie.Place == 2 {
			countPlace++
		}
	}

	sum := &resp.MovieSum{
		MovieNumber:  len(movies),
		Rate:         rateSum / len(movies),
		CinemaNumber: countPlace,
	}
	return sum
}
