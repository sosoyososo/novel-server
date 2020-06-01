package NovelSpider

import (
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
	return defaultDB.Create(c).Error
}

func CatelogWithDetailURLHasLoaded(detailUrl string) (bool, error) {
	var c int
	return c > 0, defaultDB.Model(CatelogInfo{}).Where("detail_url = ?", detailUrl).
		Count(&c).Error
}
func CatelogNumOfSummary(summaryId string) (int, error) {
	var c int
	return c, defaultDB.Model(CatelogInfo{}).Where("novel_id = ?", summaryId).
		Count(&c).Error
}

func ListCatelog(page, size int) (*[]CatelogInfo, int, error) {
	var c int
	err := defaultDB.Model(CatelogInfo{}).Count(&c).Error
	if nil != err {
		return nil, 0, err
	}
	var list []CatelogInfo
	return &list, c, defaultDB.Model(CatelogInfo{}).Offset(page * size).
		Limit(size).Scan(&list).Error
}

func CatelogListOfNovel(novelID string, page, size int) (*[]CatelogInfo, int, error) {
	defer func() {
		utils.CallFuncInNewRecoveryRoutine(func() {
			throttleTask.ThrottleDurationTask("CatelogListOfNovel", time.Hour*24,
				func() error {
					// load novel chapters list task
					return nil
				})
		})
	}()

	var c int
	err := defaultDB.Model(CatelogInfo{}).Where("novel_id = ?", novelID).
		Count(&c).Error
	if nil != err {
		return nil, 0, err
	}

	var list []CatelogInfo
	return &list, c, defaultDB.Model(CatelogInfo{}).Where("novel_id = ?", novelID).
		Offset(page * size).Limit(size).Scan(&list).Error
}

func CatelogPageUrlListOfNovel(novelID string) ([]string, error) {
	type urlContainer struct {
		URL string
	}
	var list []urlContainer
	err := defaultDB.Model(CatelogInfo{}).
		Where("novel_id = ?", novelID).
		Select("detail_url as url").
		Scan(&list).Error
	if nil != err {
		return nil, err
	}

	var ret []string
	for _, v := range list {
		ret = append(ret, v.URL)
	}
	return ret, nil
}
