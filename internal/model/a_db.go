package model

import (
	"fmt"
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/pkg/setting"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Base struct {
	CreatedOn int64 `json:"created_on"`
	UpdatedOn int64 `json:"updated_on"`
}

type BaseDbModel struct {
	*gorm.DB
}

func NewBaseModel(model interface{}, db *gorm.DB) *BaseDbModel {
	return &BaseDbModel{DB: db.Model(model)}
}

func NewDBEngine(setting *setting.DatabaseSettingS) (*gorm.DB, error) {
	dns := fmt.Sprintf(`%v:%v@(%v)/%v?charset=%v`,
		setting.UserName,
		setting.Password,
		setting.Host,
		setting.DBName,
		setting.Charset,
	)
	db, err := gorm.Open("mysql", dns)
	if err != nil {
		return nil, err
	}

	if global.ServerSetting.RunMode == "debug" {
		db.LogMode(true)
	}
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(setting.MaxIdleConns)
	db.DB().SetMaxOpenConns(setting.MaxOpenConns)
	db.AutoMigrate(
		new(Event),
		new(Task),
		new(Parent),
		new(Stuff),
		new(Tag),
		new(Project),
		new(Book),
		new(Series),
		new(Run),
		new(Abbr),
		new(PTag),
	)

	return db, nil
}
