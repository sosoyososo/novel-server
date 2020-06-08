package NovelSpider

import (
	"errors"
	"fmt"

	"../utils"
	"github.com/jinzhu/gorm"
)

type DetailInfo struct {
	BaseModel
	Title         string `json:"title"`
	Content       string `json:"content"`
	NovelID       string `json:"novelID"`
	ChapterID     string `json:"chapterID"`
	UpdateTimeStr string `json:"updateTimeStr"`
}

func detailTableNameWithNovelID(novelID string) string {
	baseName := "detail_infos_"
	if len(novelID) > 4 {
		return baseName + novelID[len(novelID)-4:]
	}
	return ""
}

func detailTableWithNovelID(novelID string) *gorm.DB {
	tableName := detailTableNameWithNovelID(novelID)
	if len(tableName) == 0 {
		return nil
	}
	if !defaultDB.HasTable(tableName) {
		defaultDB.CreateTable(&DetailInfo{})
	}
	return defaultDB.Table(tableName)
}

func (c *DetailInfo) createTableIfNeeded() error {
	tableName := detailTableNameWithNovelID(c.NovelID)
	if len(tableName) <= 0 {
		return nil
	}
	if !defaultDB.HasTable(tableName) {
		return defaultDB.CreateTable(&DetailInfo{}).Error
	}
	return nil
}

func (d *DetailInfo) Create() error {
	d.initBase()
	err := d.createTableIfNeeded()
	if nil != err {
		return err
	}
	return defaultDB.SyncW(func(db *utils.DBTools) error {
		return detailTableWithNovelID(d.NovelID).Create(d).Error
	})
}

func ChapterDetail(novelId, id string) (*DetailInfo, error) {
	var detail DetailInfo
	err := defaultDB.SyncR(func(db *utils.DBTools) error {
		queryResult := detailTableWithNovelID(novelId).Where("chapter_id = ?", id).Scan(&detail)
		if !queryResult.RecordNotFound() && queryResult.Error != nil {
			return queryResult.Error
		}
		return nil
	})
	if nil != err {
		return nil, err
	}

	var catelog CatelogInfo
	err = defaultDB.SyncW(func(db *utils.DBTools) error {
		return cateTableWithNovelID(novelId).Where("id = ?", id).Scan(&catelog).Error
	})
	if nil != err {
		return nil, err
	}

	conf := LoadConf(catelog.ConfKey)
	if nil == conf {
		return nil, errors.New("没找到对应配置")
	}
	return conf.loadDetail(catelog.DetailURL, &catelog)
}

func DetailSplit() {
	var c int
	err := defaultDB.Model(DetailInfo{}).Count(&c).Error
	if nil != err {
		panic(err)
	}
	pageC := c / 1000
	if c%1000 > 0 {
		pageC += 1
	}
	for page := 0; page < pageC; page++ {
		var list []DetailInfo
		err := defaultDB.Model(DetailInfo{}).Offset(1000 * page).Limit(1000).Scan(&list).Error
		if nil != err {
			panic(err)
		}
		for i, v := range list {
			tableName := detailTableNameWithNovelID(v.NovelID)
			if !defaultDB.HasTable(tableName) {
				err := defaultDB.Table(tableName).CreateTable(&DetailInfo{}).Error
				if nil != err {
					panic(err)

				}
			}
			err := defaultDB.Table(tableName).FirstOrCreate(&v, "id = ?", v.ID).Error
			fmt.Printf("%v/%v %v\n", page*1000+i, c, v.ID)
			if nil != err {
				panic(err)
			}
		}
	}
}
