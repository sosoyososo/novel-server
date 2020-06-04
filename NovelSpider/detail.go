package NovelSpider

import (
	"errors"

	"../utils"
)

type DetailInfo struct {
	BaseModel
	Title         string `json:"title"`
	Content       string `json:"content"`
	NovelID       string `json:"novelID"`
	ChapterID     string `json:"chapterID"`
	UpdateTimeStr string `json:"updateTimeStr"`
}

func (d *DetailInfo) Create() error {
	d.initBase()
	return defaultDB.SyncW(func(db *utils.DBTools) error {
		return db.Create(d).Error
	})
}

func ChapterDetail(id string) (*DetailInfo, error) {
	var detail DetailInfo
	err := defaultDB.SyncR(func(db *utils.DBTools) error {
		queryResult := db.Model(DetailInfo{}).Where("chapter_id = ?", id).Scan(&detail)
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
		return db.Model(CatelogInfo{}).Where("id = ?", id).Scan(&catelog).Error
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
