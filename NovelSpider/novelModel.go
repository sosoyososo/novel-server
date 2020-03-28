package NovelSpider

import (
	"errors"
	"time"

	"../utils"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	defaultDB     *gorm.DB
	loadedSummary = []string{}
)

func init() {
	db, err := gorm.Open("sqlite3", utils.GetPathRelativeToProjRoot("./gorm.db"))
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

func isSummaryLoaded(pageUrl string) bool {
	for _, url := range loadedSummary {
		if url == pageUrl {
			return true
		}
	}
	return false
}

func markSummaryLoaded(pageUrl string) {
	loadedSummary = append(loadedSummary, pageUrl)
}

type BaseModel struct {
	ID              string `json:"id" gorm:"size:40;primary_key"`
	CreateTimeStamp int64  `json:"create"`
	UpdateTimeStamp int64  `json:"update"`
	Closed          bool   `json:"-"`
	AbsoluteURL     string `json:"absoluteURL"`
	MD5             string `json:"-"`
	ConfKey         string `json:"confKey" comment:"配置key(对应配置文件的设置)"`
}

type Summary struct {
	BaseModel
	Title      string `json:"title"`
	Author     string `json:"author"`
	CoverURL   string `json:"coverURL"`
	Summary    string `json:"summary" gorm:"size:250"`
	CatelogURL string `json:"catelogURL"`
}

type Tags struct {
	BaseModel
	NovelID string `json:"novelID"`
	Tag     string `json:"tag"`
}

type CatelogInfo struct {
	BaseModel
	Title     string `json:"title"`
	NovelID   string `json:"novelID"`
	DetailURL string `json:"detailURL"`
}

type DetailInfo struct {
	BaseModel
	Title         string `json:"title"`
	Content       string `json:"content"`
	NovelID       string `json:"novelID"`
	ChapterID     string `json:"chapterID"`
	UpdateTimeStr string `json:"updateTimeStr"`
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

func (s *Summary) Create() error {
	s.initBase()
	return defaultDB.Create(s).Error
}

func (c *CatelogInfo) Create() error {
	c.initBase()
	return defaultDB.Create(c).Error
}

func (d *DetailInfo) Create() error {
	d.initBase()
	return defaultDB.Create(d).Error
}

func CatelogWithDetailURLHasLoaded(detailUrl string) (bool, error) {
	var c int
	return c > 0, defaultDB.Model(CatelogInfo{}).Where("detail_url = ?", detailUrl).Count(&c).Error
}
func CatelogNumOfSummary(summaryId string) (int, error) {
	var c int
	return c, defaultDB.Model(CatelogInfo{}).Where("novel_id = ?", summaryId).Count(&c).Error
}

func ListSummary(page, size int) (*[]Summary, error) {
	var list []Summary
	return &list, defaultDB.Model(Summary{}).Offset(page * size).Limit(size).Scan(&list).Error
}

func ListCatelog(page, size int) (*[]CatelogInfo, error) {
	var list []CatelogInfo
	return &list, defaultDB.Model(CatelogInfo{}).Offset(page * size).Limit(size).Scan(&list).Error
}

func SummaryDetail(ID string) (*Summary, error) {
	var ret Summary
	return &ret, defaultDB.Model(Summary{}).Where("id = ?", ID).Scan(&ret).Error
}

func ChapterListOfNovel(novelID string) (*[]CatelogInfo, error) {
	var list []CatelogInfo
	return &list, defaultDB.Model(CatelogInfo{}).Where("novel_id = ?", novelID).Scan(&list).Error
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

func InitCatelog() {
	for page := 0; ; page++ {
		list, err := ListCatelog(page, 20)
		if nil != err {
			panic(err)
		}
		if len(*list) == 0 {
			break
		}
		for i, _ := range *list {
			s := (*list)[i]
			if len(s.DetailURL) <= 0 {
				continue
			}
			s.initBase()
			defaultDB.Model(CatelogInfo{}).Where("detail_url = ?", s.DetailURL).Update(s)
		}
	}
}

func UpdateCatelogNovelID() {
	for page := 0; ; page++ {
		list, err := ListSummary(page, 20)
		if nil != err {
			panic(err)
		}
		if len(*list) == 0 {
			break
		}
		for i, _ := range *list {
			summary := (*list)[i]
			defaultDB.Model(CatelogInfo{}).
				Where("absolute_url = ?", summary.AbsoluteURL).
				Update("novel_id", summary.ID)
		}
	}
}
