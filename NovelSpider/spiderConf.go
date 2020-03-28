package NovelSpider

import (
	"net/url"
	"regexp"
	"strings"

	"../utils"
	"github.com/PuerkitoBio/goquery"
)

var (
	confMap = map[string]SpiderConfRef{}
)

type SpiderConfRef *SpiderConf

type NodeSel struct {
	Sel  string `json:"sel"`
	Attr string `json:"attr"`
}

type SelectorConf struct {
	Sel          string             `json:"sel"`
	IsList       bool               `json:"isList"`
	SubSelKeyMap map[string]NodeSel `json:"subSelKeyMap"`
}

func (conf *SelectorConf) ParseConfList(sel *goquery.Selection) []map[string]interface{} {
	list := []map[string]interface{}{}
	sel.Each(func(_ int, sel *goquery.Selection) {
		ret := conf.ParseConf(sel)
		list = append(list, ret)
	})
	return list
}

func (conf *SelectorConf) ParseConf(sel *goquery.Selection) map[string]interface{} {
	handleNode := func(node NodeSel, sel *goquery.Selection) string {
		if len(node.Attr) > 0 {
			c, _ := sel.Attr(node.Attr)
			return c
		} else {
			return sel.First().Text()
		}
	}
	ret := map[string]interface{}{}
	for key, node := range conf.SubSelKeyMap {
		ret[key] = handleNode(node, sel.Find(node.Sel))
	}
	return ret
}

type SpiderConf struct {
	ConfKey             string       `json:"confKey"`
	Charset             string       `json:"charset"`
	BaseURL             string       `json:"baseUrl"`
	EntryURL            string       `json:"entryUrl"`
	SummaryURLRegExp    string       `json:"summaryURLRegExp"`
	InvalidURLRegExp    []string     `json:"invalidURLRegExp"`
	SummarySelectorConf SelectorConf `json:"summarySelectorConf"`
	CatelogSelectorConf SelectorConf `json:"catelogSelectorConf"`
	DetailSelectorConf  SelectorConf `json:"detailSelectorConf"`

	loadedUrls []string
}

func init() {
	var confs map[string]SpiderConfRef
	err := utils.NewJsonConfig(utils.GetPathRelativeToProjRoot("./selConf.json"), &confs)
	if nil != err {
		panic(err)
	}
	confMap = confs
}

func LoadConf(key string) *SpiderConf {
	conf := confMap[key]
	return conf
}

func (conf *SpiderConf) HasLoaded(path string) bool {
	for _, url := range conf.loadedUrls {
		if url == path {
			return true
		}
	}
	return false
}

func (conf *SpiderConf) MarkLoaded(path string) {
	conf.loadedUrls = append(conf.loadedUrls, path)
}

func (conf *SpiderConf) IsValid(path string) bool {
	for _, regexStr := range conf.InvalidURLRegExp {
		ismatch, err := regexpMatch(regexStr, path)
		if nil != err {
			utils.ErrorLogger.Logf("%v %v\n", utils.PrintFuncName(), err)
			continue
		}
		if ismatch {
			return false
		}
	}
	return true
}

func (conf *SpiderConf) IsSummaryPage(path string) bool {
	isMatch, err := regexpMatch(conf.SummaryURLRegExp, path)
	if nil != err {
		utils.ErrorLogger.Logf("%v %v\n", utils.PrintFuncName(), err)
	}
	return isMatch
}

func (conf *SpiderConf) ToAbsolutePath(path string) string {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	if strings.HasPrefix(path, "www.") {
		return "https://" + path
	}

	baseUrl := conf.BaseURL
	if strings.HasSuffix(baseUrl, "/") {
		baseUrl = baseUrl[:len(baseUrl)-1]
	}
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	return baseUrl + "/" + path
}

func (conf *SpiderConf) InSameSite(fullUrl string) bool {
	u1, err := url.Parse(fullUrl)
	if nil != err {
		return false
	}
	u2, err := url.Parse(conf.BaseURL)
	if nil != err {
		utils.ErrorLogger.Logf("wong conf base url %v", conf.BaseURL)
		panic(err)
	}
	return u1.Hostname() == u2.Hostname()
}

/*******************************************************/
func regexpMatch(regexpStr, str string) (bool, error) {
	reg, err := regexp.Compile(regexpStr)
	if nil != err {
		utils.ErrorLogger.Logf("%v %v\n", utils.PrintFuncName(), err)
		return false, err
	}
	return reg.MatchString(str), nil
}
