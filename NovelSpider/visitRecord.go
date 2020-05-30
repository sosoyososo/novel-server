package NovelSpider

type VisitOP struct {
	NovelID string `json:"novelId"`
	Tag     string `json:"tag"`
}
type VisitInfo struct {
	BaseModel
	VisitOP
}

func (op VisitOP) Visit() error {
	info := VisitInfo{
		VisitOP: op,
	}
	info.initBase()
	return defaultDB.Create(info).Error
}
