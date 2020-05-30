package NovelSpider

var (
	loadedSummary = []string{}
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
	return defaultDB.Create(s).Error
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

func ListSummary(page, size int) (*[]Summary, int, error) {
	var c int
	err := defaultDB.Model(Summary{}).Count(&c).Error
	if nil != err {
		return nil, 0, err
	}
	var list []Summary
	return &list, c, defaultDB.Model(Summary{}).
		Offset(page * size).
		Limit(size).
		Scan(&list).Error
}

func SummaryDetail(ID string) (*Summary, error) {
	var ret Summary
	return &ret, defaultDB.Model(Summary{}).Where("id = ?", ID).Scan(&ret).Error
}

func SearchSummary(key string) (*[]Summary, error) {
	var list []Summary
	like := "%" + key + "%"
	return &list, defaultDB.Model(Summary{}).
		Where("title like ? or summary like ? or author like ?", like, like, like).
		Scan(&list).Error
}
