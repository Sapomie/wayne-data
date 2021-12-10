package c_movie

import (
	"errors"
	"fmt"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"strconv"
	"strings"
)

func (svc MovieService) ProcessMovie() ([]string, error) {
	movies, infos, err := svc.makeMovies()
	if err != nil {
		return nil, err
	}

	err = svc.storeMovies(movies)
	if err != nil {
		return nil, err
	}

	return infos, nil
}

func (svc MovieService) makeMovies() (movies model.Movies, infos []string, err error) {

	start, end := mtime.NewTimeZone(mtime.TypeYear, 2021, 1).BeginAndEnd()
	events, err := model.NewEventModel(svc.db).ByTaskName(start, end, cons.Movie)
	if err != nil {
		return nil, nil, err
	}

	for _, event := range events {
		movie, err := makeMovie(event)
		if err != nil {
			info := fmt.Sprintf("make movie error,event start time: %v,coment: %v", event.Start(), event.Comment)
			infos = append(infos, info)
			continue
		}
		movies = append(movies, movie)
	}

	return
}

func makeMovie(event *model.Event) (*model.Movie, error) {
	fields := strings.Split(event.Comment, "，")
	if len(fields) < 5 {
		return nil, errors.New("not enough fields")
	}
	rate, err := strconv.Atoi(fields[2])
	if err != nil {
		return nil, errors.New("movie rate parse error")
	}
	if rate < 0 || rate > 100 {
		return nil, errors.New("movie rate out of range")
	}
	place, err := strconv.Atoi(fields[3])
	if err != nil {
		return nil, errors.New("movie place parse error")
	}
	if place < 1 || place > 3 {
		return nil, errors.New("movie place out of range")
	}
	year, err := strconv.Atoi(fields[4])
	if err != nil {
		return nil, errors.New("movie year parse error")
	}
	if year < 1900 || year > 2100 {
		return nil, errors.New("movie year out of range")
	}
	movie := &model.Movie{
		Date:      event.Date,
		Name:      fields[0],
		EnName:    fields[1],
		Rate:      rate,
		Place:     place,
		Year:      year,
		StartTime: event.StartTime,
	}
	if len(fields) == 6 {
		movie.Comment = fields[5]
	}

	return movie, nil
}

func (svc MovieService) storeMovies(movies model.Movies) error {
	for _, movie := range movies {
		mm := model.NewMovieModel(svc.db)
		exist, err := mm.Exists(movie.StartTime)
		if err != nil {
			return err
		}
		if exist {
			err := mm.Base.Where("start_time = ?", movie.StartTime).Update(movie).Error
			if err != nil {
				return err
			}
		} else {
			err := mm.Base.Create(movie).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}
