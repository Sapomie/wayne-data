package model

import (
	"fmt"
	"github.com/Sapomie/wayne_data/global"
	"github.com/Sapomie/wayne_data/pkg/setting"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Model struct {
	Id        uint64 `gorm:"primary_key" json:"id"`
	CreatedOn int64  `json:"created_on"`
	UpdatedOn int64  `json:"updated_on"`
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

	return db, nil
}
