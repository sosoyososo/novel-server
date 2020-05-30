package NovelSpider

type Tags struct {
	BaseModel
	NovelID string `json:"novelID"`
	Tag     string `json:"tag"`
}
