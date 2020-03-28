package utils

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type DBConfig struct {
	Host   string `json:"host"`
	User   string `json:"user"`
	Pswd   string `json:"pswd"`
	DBName string `json:"dbName"`
}

func (conf *DBConfig) CreateDB() (*DBTools, error) {
	urlstr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.User,
		conf.Pswd,
		conf.Host,
		conf.DBName,
	)
	db, err := gorm.Open("mysql", urlstr)
	if nil == err {
		db = db.Set("gorm:table_options", "ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=1;")
	}
	if nil != err {
		return nil, err
	}
	db = db.LogMode(true)
	db.SetLogger(log.New(InfoLogger, "\r\n", 0))
	return &DBTools{DB: db, InTransaction: false}, nil
}

type DBTools struct {
	*gorm.DB
	InTransaction bool
}

func (d *DBTools) Begin() *DBTools {
	return &DBTools{DB: d.DB.Begin(), InTransaction: true}
}

func (d *DBTools) Transaction(ac func(tx *DBTools) error) error {
	db := d.Begin()
	err := ac(db)
	if nil != err {
		db.Rollback()
		return err
	}
	return db.Commit().Error
}
