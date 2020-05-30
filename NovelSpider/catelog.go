package NovelSpider

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
	return c > 0, defaultDB.Model(CatelogInfo{}).Where("detail_url = ?", detailUrl).Count(&c).Error
}
func CatelogNumOfSummary(summaryId string) (int, error) {
	var c int
	return c, defaultDB.Model(CatelogInfo{}).Where("novel_id = ?", summaryId).Count(&c).Error
}

func ListCatelog(page, size int) (*[]CatelogInfo, error) {
	var list []CatelogInfo
	return &list, defaultDB.Model(CatelogInfo{}).Offset(page * size).Limit(size).Scan(&list).Error
}

func ChapterListOfNovel(novelID string) (*[]CatelogInfo, error) {
	var list []CatelogInfo
	return &list, defaultDB.Model(CatelogInfo{}).Where("novel_id = ?", novelID).Scan(&list).Error
}

func ChapterPageUrlListOfNovel(novelID string) ([]string, error) {
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
