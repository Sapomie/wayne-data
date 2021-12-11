package c_movie

import (
	"context"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/model/cons"
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

func (svc MovieService) ListMovies() (*resp.Movie, error) {
	movies := new(resp.Movie)
	key := cons.RedisKeyMovie
	exists, err := svc.cache.Get(key, &movies)
	if err != nil {
		return nil, err
	}
	if !exists {
		movies, err = svc.GetMovieFromDB()
		if err != nil {
			return nil, err
		}
		err = svc.cache.Set(key, movies, 0)
		if err != nil {
			return nil, err
		}
	}

	return movies, nil
}

func (svc MovieService) GetMovieFromDB() (*resp.Movie, error) {
	movies, err := model.NewMovieModel(svc.db).GetAll()
	if err != nil {
		return nil, err
	}

	movieResponses := make([]*resp.MovieItem, 0)
	for _, movie := range movies {
		movieResp := toMovieResponse(movie)
		movieResponses = append(movieResponses, movieResp)
	}

	return &resp.Movie{
		Items: movieResponses,
		Sum:   toMovieSum(movies),
	}, nil
}

func toMovieResponse(m *model.Movie) *resp.MovieItem {
	place := ""
	if m.Place == 2 {
		place = "Cinema"
	}

	return &resp.MovieItem{
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
