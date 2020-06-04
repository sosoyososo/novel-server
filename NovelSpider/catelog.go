package NovelSpider

import (
	"errors"
	"time"

	"../throttleTask"
	"../utils"
)

type CatelogInfo struct {
	BaseModel
	Title     string `json:"title"`
	NovelID   string `json:"novelID"`
	DetailURL string `json:"detailURL"`
}

func (c *CatelogInfo) Create() error {
	c.initBase()
	return defaultDB.SyncW(func(db *utils.DBTools) error {
		return db.Create(c).Error
	})
}

func CatelogWithDetailURLHasLoaded(detailUrl string) (bool, error) {
	var c int

	return c > 0, defaultDB.SyncR(func(db *utils.DBTools) error {
		return db.Model(CatelogInfo{}).Where("detail_url = ?", detailUrl).
			Count(&c).Error
	})
}
func CatelogNumOfSummary(summaryId string) (int, error) {
	var c int
	return c, defaultDB.SyncR(func(db *utils.DBTools) error {
		return db.Model(CatelogInfo{}).Where("novel_id = ?", summaryId).
			Count(&c).Error
	})
}

func ListCatelog(page, size int) (*[]CatelogInfo, int, error) {
	var c int
	err := defaultDB.SyncR(func(db *utils.DBTools) error {
		return db.Model(CatelogInfo{}).Count(&c).Error
	})
	if nil != err {
		return nil, 0, err
	}
	var list []CatelogInfo
	return &list, c, defaultDB.SyncR(func(db *utils.DBTools) error {
		return db.Model(CatelogInfo{}).Offset(page * size).
			Limit(size).Scan(&list).Error
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
		return db.Model(CatelogInfo{}).Where("novel_id = ?", novelID).
			Count(&c).Error
	})
	if nil != err {
		return nil, 0, err
	}

	var list []CatelogInfo
	return &list, c, defaultDB.SyncR(func(db *utils.DBTools) error {
		return db.Model(CatelogInfo{}).Where("novel_id = ?", novelID).
			Offset(page * size).Limit(size).Scan(&list).Error
	})
}

func CatelogPageUrlListOfNovel(novelID string) ([]string, error) {
	type urlContainer struct {
		URL string
	}
	var list []urlContainer
	err := defaultDB.SyncR(func(db *utils.DBTools) error {
		return db.Model(CatelogInfo{}).
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
