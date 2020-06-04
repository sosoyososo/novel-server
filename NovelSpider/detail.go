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
	queryResult := defaultDB.Model(DetailInfo{}).Where("chapter_id = ?", id).Scan(&detail)
	if !queryResult.RecordNotFound() && queryResult.Error != nil {
		return nil, queryResult.Error
	}

	var catelog CatelogInfo
	err := defaultDB.Model(CatelogInfo{}).Where("id = ?", id).Scan(&catelog).Error
	if nil != err {
		return nil, err
	}

	conf := LoadConf(catelog.ConfKey)
	if nil == conf {
		return nil, errors.New("没找到对应配置")
	}
	return conf.loadDetail(catelog.DetailURL, &catelog)
}
