package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Book struct {
	Id               int64  `gorm:"primary_key"`
	Name             string `gorm:"not null"`
	Category         string
	Author           string
	Year             int
	WordNumber       float64
	Rate             int
	Duration         float64
	FirstReadingTime int64
	LastReadingTime  int64
	Finish           int8
	CreatedTime      int64 `gorm:"not null" json:"created_time"`
	UpdatedTime      int64 `gorm:"not null" json:"updated_time"`
}

func (e *Book) TableName() string {
	return "c_book"
}

type Books []*Book

type BookModel struct {
	Base *BaseDbModel
}

func NewBookModel(db *gorm.DB) *BookModel {
	return &BookModel{NewBaseModel(new(Book), db)}
}

func (em *BookModel) Exists(name string) (bool, error) {
	var count int
	err := em.Base.Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	exists := count > 0
	return exists, nil
}

func (em *BookModel) GetAll() (Books, error) {
	var books Books
	err := em.Base.Order("last_time desc").Scan(&books).Error
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (em *BookModel) ByName(name string) (*Book, error) {
	book := new(Book)
	err := em.Base.Where("name = ?", name).Scan(book).Error
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (em *BookModel) ListBooks(limit, offset int) (Books, int, error) {
	var (
		books Books
		count int
	)
	db := em.Base.DB
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(limit).Offset(offset).Scan(&books).Error
	if err != nil {
		return nil, 0, err
	}
	return books, count, nil
}

func (em *BookModel) InsertAndGetBook(name string) (book *Book, info string, err error) {
	exists, err := em.Exists(name)
	if err != nil {
		return nil, "", err
	}
	if !exists {
		err = em.Base.Create(&Book{Name: name}).Error
		if err != nil {
			return nil, "", err
		}
		info = fmt.Sprintf("Add Book %v ", name)
	}
	book, err = em.ByName(name)
	if err != nil {
		return nil, "", err
	}
	return
}
