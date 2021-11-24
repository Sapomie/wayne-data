package c_book

import (
	"context"
	"errors"
	"fmt"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

type BookService struct {
	ctx     context.Context
	cache   *model.Cache
	bookDb  *model.BookModel
	eventDb *model.EventModel
}

func NewBookService(c context.Context, db *gorm.DB, cache *redis.Pool) BookService {
	return BookService{
		ctx:     c,
		cache:   model.NewCache(cache),
		bookDb:  model.NewBookModel(db),
		eventDb: model.NewEventModel(db),
	}
}

func (svc BookService) ProcessBook() ([]string, error) {
	books, infos, err := svc.makeBooks()
	if err != nil {
		return nil, err
	}

	err = svc.storeBooks(books)
	if err != nil {
		return nil, err
	}

	return infos, nil
}

func (svc BookService) makeBooks() (books model.Books, infos []string, err error) {

	start, end := mtime.NewTimeZone(mtime.TypeYear, 2021, 1).BeginAndEnd()
	events, err := svc.eventDb.ByTaskName(start, end, cons.Nonfiction)
	if err != nil {
		return nil, nil, err
	}

	bookMap := make(map[string]model.Events, 0)
	for _, event := range events {
		strs := strings.Split(event.Comment, "，")
		name := strs[0]
		bookMap[name] = append(bookMap[name], event)
	}

	for name, bookEvents := range bookMap {
		book := &model.Book{Name: name}
		for _, event := range bookEvents {

			if isBookFirstTime(event) {
				book.Category, book.Author, book.Year, book.WordNumber, err = bookInfo(event)
				book.FirstTime = event.StartTime
				if err != nil {
					infos = append(infos, fmt.Sprintf("make book error,event start time: %v,coment: %v", event.Start(), event.Comment))
					continue
				}
			}

			if isBookLastTime(event) {
				book.Rate, err = bookRate(event)
				if err != nil {
					infos = append(infos, fmt.Sprintf("make book error,event start time: %v,coment: %v", event.Start(), event.Comment))
					continue
				}
				book.Finish = model.ProjectFinish
			}

			if event.StartTime > book.LastTime {
				book.LastTime = event.StartTime
			}

			book.Duration += event.Duration
		}
		books = append(books, book)
	}

	return
}

func (svc BookService) storeBooks(books model.Books) error {
	mm := svc.bookDb
	for _, book := range books {
		exist, err := mm.Exists(book.Name)
		if err != nil {
			return err
		}
		if exist {
			err := mm.Base.Where("name = ?", book.Name).Update(book).Error
			if err != nil {
				return err
			}
		} else {
			err := mm.Base.Create(book).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isBookFirstTime(event *model.Event) bool {
	if strings.Contains(event.Comment, "、s") {
		return true
	}
	return false
}

func isBookLastTime(event *model.Event) bool {
	if strings.Contains(event.Comment, "、e") {
		return true
	}
	return false
}

func bookInfo(event *model.Event) (category, author string, year int, word float64, err error) {
	strs := strings.Split(event.Comment, "，")
	if len(strs) < 6 {
		return "", "", 0, 0, errors.New("wrong length of book comment")
	}
	category = strs[1]
	author = strs[2]
	year, err = strconv.Atoi(strs[3])
	if err != nil {
		return "", "", 0, 0, err
	}

	word, err = strconv.ParseFloat(strs[4], 64)
	if err != nil {
		return "", "", 0, 0, err
	}

	return
}

func bookRate(event *model.Event) (rate int, err error) {
	strs := strings.Split(event.Comment, "，")
	if len(strs) < 3 {
		return 0, errors.New("wrong length of book comment")
	}
	rate, err = strconv.Atoi(strs[1])
	if err != nil {
		return 0, err
	}
	return
}
