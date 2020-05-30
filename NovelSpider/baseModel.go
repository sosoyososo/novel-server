package NovelSpider

import (
	"reflect"
	"time"

	"../PO"
	"../utils"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	defaultDB *utils.DBTools
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
	initSummary()
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
	defaultDB = utils.CreateCustomDBFromGorm(db)

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

func ListModel(
	db *utils.DBTools, page, size int,
	model interface{}, preload, condition func(db *gorm.DB) *gorm.DB,
) (interface{}, int, error) {
	if nil == db {
		db = defaultDB
	}
	modelType := reflect.TypeOf(model)
	query := db.Model(model)

	var total int
	queryResult := query.Where("closed <> 1")
	if nil != condition {
		queryResult = condition(queryResult)
	}
	queryResult = queryResult.Count(&total)
	if queryResult.Error != nil {
		if queryResult.RecordNotFound() {
			return nil, 0, nil
		} else {
			return nil, 0, queryResult.Error
		}
	}

	list := reflect.New(
		reflect.SliceOf(modelType),
	).Interface()
	if nil != preload {
		queryResult = preload(query)
	}
	queryResult = queryResult.Model(model)
	if nil != condition {
		queryResult = condition(queryResult)
	}
	queryResult.
		Offset(size * page).
		Limit(size).
		Find(list)
	if queryResult.Error != nil {
		if queryResult.RecordNotFound() {
			return nil, 0, nil
		} else {
			return nil, 0, queryResult.Error
		}
	}
	return list, total, nil
}

func ModelDetail(
	id string, model interface{},
	preload, condition func(db *gorm.DB) *gorm.DB, db *utils.DBTools,
) (interface{}, error) {
	if nil == model {
		return nil, PO.NormalErr("数据类型错误")
	}
	if nil == db {
		db = defaultDB
	}

	queryResult := db.Model(model)
	if nil != preload {
		queryResult = preload(queryResult)
	}
	queryResult = queryResult.Model(model).Where("closed <> 1 and id = ?", id)
	if nil != condition {
		queryResult = condition(queryResult)
	}
	err := queryResult.Find(model).Error
	return model, err
}
