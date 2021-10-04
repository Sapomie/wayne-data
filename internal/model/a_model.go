package model

import (
	"fmt"
	"github.com/Sapomie/wayne_data/pkg/setting"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Model struct {
	Id        uint64 `gorm:"primary_key" json:"id"`
	CreatedOn int64  `json:"created_on"`
	UpdatedOn int64  `json:"updated_on"`
}

func NewDBEngine(ds *setting.DatabaseSettingS) (*gorm.DB, error) {
	dns := fmt.Sprintf(`%v:%v@(%v)/%v?charset=%v`,
		ds.UserName, ds.Password, ds.Host, ds.DBName, ds.Charset)
	db, err := gorm.Open("mysql", dns)
	if err != nil {
		return nil, err
	}
	db.SingularTable(true)

	return db, nil
}
