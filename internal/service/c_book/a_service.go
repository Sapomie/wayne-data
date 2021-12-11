package c_book

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

type BookService struct {
	ctx   context.Context
	cache *model.Cache
	db    *gorm.DB
}

func NewBookService(c context.Context, db *gorm.DB, cache *redis.Pool) BookService {
	return BookService{
		ctx:   c,
		cache: model.NewCache(cache),
		db:    db,
	}
}

func (svc BookService) ListBooks() (*resp.BookResp, error) {
	resp := new(resp.BookResp)
	key := cons.RedisKeyBook
	exists, err := svc.cache.Get(key, &resp)
	if err != nil {
		return nil, err
	}
	if !exists {
		resp, err = svc.GetBooksFromDB()
		if err != nil {
			return nil, err
		}
		err = svc.cache.Set(key, resp, 0)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func (svc BookService) GetBooksFromDB() (*resp.BookResp, error) {
	books, err := model.NewBookModel(svc.db).GetAll()
	if err != nil {
		return nil, err
	}

	bookResponses := make([]*resp.BookItemResp, 0)
	for _, book := range books {
		bookResp := toBookResponse(book)
		bookResponses = append(bookResponses, bookResp)
	}

	return &resp.BookResp{
		Items: bookResponses,
		Sum:   toBookSum(books),
	}, nil
}

func toBookResponse(b *model.Book) *resp.BookItemResp {

	var finishMark string
	switch b.Finish {
	case model.BookFinish:
		finishMark = "Finish"
	case model.BookAbandon:
		finishMark = "Abandon"
	}

	return &resp.BookItemResp{
		Name:       b.Name,
		Category:   b.Category,
		Author:     b.Author,
		Year:       b.Year,
		WordNumber: b.WordNumber,
		Duration:   b.Duration,
		Rate:       b.Rate,
		Finish:     finishMark,
		FirstTime:  time.Unix(b.FirstTime, 0).Format(mtime.TimeTemplate4),
		LastTime:   time.Unix(b.LastTime, 0).Format(mtime.TimeTemplate4),
	}

}

func toBookSum(books model.Books) *resp.BookSumResp {
	var (
		finishNum         int
		wordNumberSum     float64
		durationFinishSum float64
		durationSum       float64
		rateSum           int
		categories        = make(map[string]int)
	)

	for _, book := range books {
		if book.Finish == model.BookFinish {
			finishNum++
			wordNumberSum += book.WordNumber
			durationFinishSum += book.Duration
			rateSum += book.Rate
		}
		durationSum += book.Duration
		categories[book.Category]++
	}

	var (
		categoryNum   = len(categories)
		wordNumberAvg = wordNumberSum / float64(finishNum)
		durationAvg   = durationSum / float64(finishNum)
		rateAvg       = rateSum / finishNum
	)

	sum := &resp.BookSumResp{
		BookNumber:     len(books),
		CategoryNumber: categoryNum,
		WordNumberAvg:  convert.FloatTo(wordNumberAvg).Decimal(2),
		DurationAvg:    convert.FloatTo(durationAvg).Decimal(2),
		RateAvg:        rateAvg,
		Finish:         finishNum,
	}

	return sum
}
