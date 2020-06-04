package NovelSpider

import (
	"sync"

	"../utils"
)

var (
	loadedSummary_check_l = &sync.Mutex{}
	loadedSummary         = []string{}
)

type Summary struct {
	BaseModel
	Title      string `json:"title"`
	Author     string `json:"author"`
	CoverURL   string `json:"coverURL"`
	Summary    string `json:"summary" gorm:"size:250"`
	CatelogURL string `json:"catelogURL"`
}

func initSummary() {
	type DBURL struct {
		AbsoluteURL string
	}

	var uList []DBURL
	err := defaultDB.Model(Summary{}).Select("absolute_url").Scan(&uList).Error
	if nil != err {
		panic(err)
	}

	for _, u := range uList {
		loadedSummary = append(loadedSummary, u.AbsoluteURL)
	}
}

func (s *Summary) Create() error {
	s.initBase()
	return defaultDB.SyncW(func(db *utils.DBTools) error {
		return db.Create(s).Error
	})
}

func summaryLoadedCheckAndMark(pageUrl string) bool {
	loadedSummary_check_l.Lock()
	defer loadedSummary_check_l.Unlock()

	for _, url := range loadedSummary {
		if url == pageUrl {
			return true
		}
	}
	loadedSummary = append(loadedSummary, pageUrl)
	return false
}

func ListSummary(page, size int) (*[]Summary, int, error) {
	var c int
	err := defaultDB.SyncW(func(db *utils.DBTools) error {
		return db.Model(Summary{}).Count(&c).Error
	})
	if nil != err {
		return nil, 0, err
	}
	var list []Summary
	return &list, c, defaultDB.SyncW(func(db *utils.DBTools) error {
		return db.Model(Summary{}).
			Offset(page * size).
			Limit(size).
			Scan(&list).Error
	})
}

func SummaryDetail(ID string) (*Summary, error) {
	var ret Summary
	return &ret, defaultDB.SyncW(func(db *utils.DBTools) error {
		return db.Model(Summary{}).Where("id = ?", ID).Scan(&ret).Error
	})
}

func SearchSummary(key string) (*[]Summary, error) {
	var list []Summary
	like := "%" + key + "%"
	return &list, defaultDB.SyncW(func(db *utils.DBTools) error {
		return db.Model(Summary{}).
			Where("title like ? or summary like ? or author like ?", like, like, like).
			Scan(&list).Error
	})
}
