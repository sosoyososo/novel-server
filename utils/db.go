package utils

import (
	"fmt"
	"log"
	"sync"

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

func CreateCustomDBFromGorm(db *gorm.DB) *DBTools {
	db = db.LogMode(true)
	db.SetLogger(log.New(DBLogger, "\r\n", 0))
	return &DBTools{DB: db, InTransaction: false}
}

type DBTools struct {
	*gorm.DB
	InTransaction bool
	l             *sync.RWMutex
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

func (d *DBTools) SyncR(f func(d *DBTools)) {
	if nil == d.l {
		d.l = &sync.RWMutex{}
	}
	d.l.RLock()
	defer d.l.RUnlock()

	f(d)
}

func (d *DBTools) SyncW(f func(d *DBTools)) {
	if nil == d.l {
		d.l = &sync.RWMutex{}
	}
	d.l.Lock()
	defer d.l.Unlock()

	f(d)
}
