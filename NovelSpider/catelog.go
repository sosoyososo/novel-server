package NovelSpider

import (
	"errors"
	"fmt"
	"time"

	"../throttleTask"
	"../utils"
	"github.com/jinzhu/gorm"
)

type CatelogInfo struct {
	BaseModel
	Title     string `json:"title"`
	NovelID   string `json:"novelID"`
	DetailURL string `json:"detailURL"`
}

func cateTableNameWithNovelID(novelID string) string {
	baseName := "catelog_infos_"
	if len(novelID) > 2 {
		return baseName + novelID[len(novelID)-2:]
	}
	return ""
}

func cateTableWithNovelID(novelID string) *gorm.DB {
	tableName := cateTableNameWithNovelID(novelID)
	if len(tableName) == 0 {
		return nil
	}
	return defaultDB.Table(tableName)
}

func (c *CatelogInfo) createTableIfNeeded() error {
	tableName := cateTableNameWithNovelID(c.NovelID)
	if len(tableName) <= 0 {
		return nil
	}
	if !defaultDB.HasTable(tableName) {
		return defaultDB.CreateTable(&CatelogInfo{}).Error
	}
	return nil
}

func (c *CatelogInfo) Create() error {
	c.initBase()
	err := c.createTableIfNeeded()
	if nil != err {
		return err
	}
	return defaultDB.SyncW(func(db *utils.DBTools) error {
		return cateTableWithNovelID(c.NovelID).Create(c).Error
	})
}

func CatelogNumOfSummary(summaryId string) (int, error) {
	var c int
	return c, defaultDB.SyncR(func(db *utils.DBTools) error {
		return cateTableWithNovelID(summaryId).Model(CatelogInfo{}).Where("novel_id = ?", summaryId).
			Count(&c).Error
	})
}

func CatelogListOfNovel(novelID string, page, size int) (*[]CatelogInfo, int, error) {
	defer func() {
		utils.CallFuncInNewRecoveryRoutine(func() {
			throttleTask.ThrottleDurationTask("CatelogListOfNovel", time.Hour*24,
				func() error {
					// load novel chapters list task
					s, err := SummaryDetail(novelID)
					if nil != err {
						return err
					}

					conf := LoadConf(s.ConfKey)
					if nil == conf {
						return errors.New("没找到对应配置")
					}

					conf.loadCatelog(s.AbsoluteURL, s)

					return nil
				})
		})
	}()

	var c int
	err := defaultDB.SyncR(func(db *utils.DBTools) error {
		return cateTableWithNovelID(novelID).Where("novel_id = ?", novelID).
			Count(&c).Error
	})
	if nil != err {
		return nil, 0, err
	}

	var list []CatelogInfo
	return &list, c, defaultDB.SyncR(func(db *utils.DBTools) error {
		return cateTableWithNovelID(novelID).Where("novel_id = ?", novelID).
			Offset(page * size).Limit(size).Scan(&list).Error
	})
}

func CatelogPageUrlListOfNovel(novelID string) ([]string, error) {
	type urlContainer struct {
		URL string
	}
	var list []urlContainer
	err := defaultDB.SyncR(func(db *utils.DBTools) error {
		return cateTableWithNovelID(novelID).
			Where("novel_id = ?", novelID).
			Select("detail_url as url").
			Scan(&list).Error
	})
	if nil != err {
		return nil, err
	}

	var ret []string
	for _, v := range list {
		ret = append(ret, v.URL)
	}
	return ret, nil
}

func CateLogSplit() {
	var c int
	err := defaultDB.Model(CatelogInfo{}).Count(&c).Error
	if nil != err {
		panic(err)
	}
	pageC := c / 1000
	if c%1000 > 0 {
		pageC += 1
	}
	for page := 0; page < pageC; page++ {
		var list []CatelogInfo
		err := defaultDB.Model(CatelogInfo{}).Offset(1000 * page).Limit(1000).Scan(&list).Error
		if nil != err {
			panic(err)
		}
		for i, v := range list {
			tableName := cateTableNameWithNovelID(v.NovelID)
			if !defaultDB.HasTable(tableName) {
				err := defaultDB.Table(tableName).CreateTable(&CatelogInfo{}).Error
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
