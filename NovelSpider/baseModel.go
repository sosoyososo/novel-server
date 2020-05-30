package NovelSpider

import (
	"log"
	"time"

	"../utils"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	defaultDB     *gorm.DB
	loadedSummary = []string{}
)

type BaseModel struct {
	ID              string `json:"id" gorm:"size:40;primary_key"`
	CreateTimeStamp int64  `json:"create"`
	UpdateTimeStamp int64  `json:"update"`
	Closed          bool   `json:"-"`
	AbsoluteURL     string `json:"absoluteURL"`
	MD5             string `json:"-"`
	ConfKey         string `json:"confKey" comment:"配置key(对应配置文件的设置)"`
}

func init() {
	initConf()
	initDB()
}

func initDB() {
	db, err := gorm.Open(
		"sqlite3",
		utils.GetPathRelativeToProjRoot("./gorm.db"),
	)
	if nil != err {
		panic(err)
	}
	list := []interface{}{
		Summary{},
		CatelogInfo{},
		DetailInfo{},
	}
	for _, m := range list {
		if !db.HasTable(m) {
			db.CreateTable(m)
		}
		db.AutoMigrate(m)
	}
	defaultDB = db

	db = db.LogMode(true)
	db.SetLogger(log.New(utils.InfoLogger, "\r\n", 0))

	type DBURL struct {
		AbsoluteURL string
	}

	var uList []DBURL
	err = defaultDB.Model(Summary{}).Select("absolute_url").Scan(&uList).Error
	if nil != err {
		panic(err)
	}

	for _, u := range uList {
		loadedSummary = append(loadedSummary, u.AbsoluteURL)
	}
}

func (b *BaseModel) initBase() error {
	id, err := utils.GenerateID()
	if err != nil {
		return err
	}
	b.ID = id
	n := time.Now()
	b.CreateTimeStamp = n.Unix()
	b.UpdateTimeStamp = n.Unix()
	return nil
}
