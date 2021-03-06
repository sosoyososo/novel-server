package NovelSpider

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"../Encoding"
	"../Html"
	"../utils"
	"github.com/PuerkitoBio/goquery"
)

const (
	mxaDownloadRoutine = 5
)

var (
	confMap            = map[string]SpiderConfRef{}
	downloaderThrottle = make(chan int, mxaDownloadRoutine)
)

type SpiderConfRef *SpiderConf

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

	checkLocker *sync.Mutex
	loadedUrls  []string
}

func initConf() {
	var confs map[string]SpiderConfRef
	err := utils.NewJsonConfig(
		utils.GetPathRelativeToProjRoot("./selConf.json"), &confs,
	)
	if nil != err {
		panic(err)
	}
	confMap = confs
	for key, _ := range confMap {
		confMap[key].checkLocker = &sync.Mutex{}
	}

	for i := 0; i < mxaDownloadRoutine; i++ {
		downloaderThrottle <- 1
	}
}

func downloaderCountChange(add bool) {
	if add {
		<-downloaderThrottle
	} else {
		downloaderThrottle <- 1
	}
}

func LoadConf(key string) *SpiderConf {
	conf := confMap[key]
	return conf
}

func (conf *SpiderConf) hasLoadCheckAndMark(path string) bool {
	conf.checkLocker.Lock()
	defer conf.checkLocker.Unlock()

	for _, url := range conf.loadedUrls {
		if url == path {
			return true
		}
	}
	conf.loadedUrls = append(conf.loadedUrls, path)
	return false
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

//TODO: all these func should be in handle queue
func (conf *SpiderConf) loadSummaryPage(pageUrl string) {
	if summaryLoadedCheckAndMark(pageUrl) {
		return
	}

	novelId := ""
	summarySelAction := Html.NewAction(conf.SummarySelectorConf.Sel, func(sel *goquery.Selection) {
		ret := conf.SummarySelectorConf.ParseConf(sel)
		ret["absoluteURL"] = pageUrl
		ret["confKey"] = conf.ConfKey

		var s Summary
		err := utils.MapToJsonObj(ret, &s)
		if nil != err {
			utils.ErrorLogger.Logf("%v %v\n", utils.PrintFuncName(), err)
		}
		err = s.Create()
		if nil != err {
			utils.ErrorLogger.Logf("%v %v\n", utils.PrintFuncName(), err)
		}
		novelId = s.ID
	})
	catelogSelAction := Html.NewAction(conf.CatelogSelectorConf.Sel, func(sel *goquery.Selection) {
		ret := conf.CatelogSelectorConf.ParseConfList(sel)
		for i, _ := range ret {
			ret[i]["absoluteURL"] = pageUrl
			ret[i]["confKey"] = conf.ConfKey
			ret[i]["novelID"] = novelId
			detailUrl := ret[i]["detailURL"]
			if detailUrl != nil {
				if pageUrl, ok := detailUrl.(string); ok {
					if len(pageUrl) > 0 {
						ret[i]["detailURL"] = conf.ToAbsolutePath(pageUrl)
					}
				}
			}
		}

		var list []CatelogInfo
		err := utils.MapListToJsonObjList(ret, &list)
		if nil != err {
			utils.ErrorLogger.Logf("%v %v\n", utils.PrintFuncName(), err)
			return
		}
		for _, c := range list {
			err = c.Create()
			if nil != err {
				utils.ErrorLogger.Logf("%v %v\n", utils.PrintFuncName(), err)
			}
		}
	})
	conf.loadPage(pageUrl, []Html.WorkerAction{summarySelAction, catelogSelAction})
}

func (conf *SpiderConf) loadCatelog(pageUrl string, s *Summary) {
	catelogSelAction := Html.NewAction(conf.CatelogSelectorConf.Sel, func(sel *goquery.Selection) {
		//parse chapters from html doc
		ret := conf.CatelogSelectorConf.ParseConfList(sel)
		for i, _ := range ret {
			ret[i]["absoluteURL"] = pageUrl
			ret[i]["confKey"] = conf.ConfKey
			ret[i]["novelID"] = s.ID
			detailUrl := ret[i]["detailURL"]
			if detailUrl != nil {
				if pageUrl, ok := detailUrl.(string); ok {
					if len(pageUrl) > 0 {
						ret[i]["detailURL"] = conf.ToAbsolutePath(pageUrl)
					}
				}
			}
		}

		//load list from db
		urlList, err := CatelogPageUrlListOfNovel(s.ID)
		if nil != err {
			utils.ErrorLogger.Logf("%v %v\n", utils.PrintFuncName(), err)
			return
		}

		//only insert not existed chapters
		if len(urlList) >= len(ret) {
			return
		}
		needCheck := true
		if len(urlList) == 0 {
			needCheck = false
		}

		var list []CatelogInfo
		err = utils.MapListToJsonObjList(ret, &list)
		if nil != err {
			utils.ErrorLogger.Logf("%v %v\n", utils.PrintFuncName(), err)
			return
		}

		for _, c := range list {
			shouldIgnoreLoaded := false
			if needCheck {
				for _, v := range urlList {
					if v == c.DetailURL {
						shouldIgnoreLoaded = true
						break
					}
				}
			}

			if shouldIgnoreLoaded {
				continue
			}

			err = c.Create()
			if nil != err {
				utils.ErrorLogger.Logf("%v %v\n", utils.PrintFuncName(), err)
			}
		}
	})
	conf.loadPage(pageUrl, []Html.WorkerAction{catelogSelAction})
}

func (conf *SpiderConf) loadDetail(pageUrl string, c *CatelogInfo) (*DetailInfo, error) {
	var detail DetailInfo
	var err error
	detailSelAction := Html.NewAction(conf.DetailSelectorConf.Sel, func(sel *goquery.Selection) {
		ret := conf.DetailSelectorConf.ParseConf(sel)
		ret["absoluteURL"] = pageUrl
		ret["confKey"] = conf.ConfKey
		ret["novelID"] = c.NovelID
		ret["chapterID"] = c.ID

		err = utils.MapToJsonObj(ret, &detail)
		if nil != err {
			utils.ErrorLogger.Logf("%v %v\n", utils.PrintFuncName(), err)
			return
		}
		err = detail.Create()
		if nil != err {
			utils.ErrorLogger.Logf("%v %v\n", utils.PrintFuncName(), err)
			return
		}
	})
	conf.loadPage(pageUrl, []Html.WorkerAction{detailSelAction})
	return &detail, err
}

func (conf *SpiderConf) LoadValidPage(pageUrl string) {
	pageUrl = conf.ToAbsolutePath(pageUrl)
	utils.InfoLogger.Logf("start load page %v", pageUrl)
	hrefSelAction := Html.NewAction("a", func(sel *goquery.Selection) {
		sel.Each(func(_ int, sel *goquery.Selection) {
			url, _ := sel.Attr("href")
			if len(url) <= 0 || url == "/" {
				return
			}

			if !conf.IsValid(url) {
				// utils.DebugLogger.Logf("invalid url %v", url)
				return
			}

			if conf.hasLoadCheckAndMark(url) {
				// utils.DebugLogger.Logf("preloaded url %v", url)
				return
			}

			if conf.IsSummaryPage(url) {
				url = conf.ToAbsolutePath(url)
				if !conf.InSameSite(url) {
					// utils.DebugLogger.Logf("skip other site url %v", pageUrl)
					return
				}
				utils.DebugLogger.Logf("find summary url %v", url)
				go conf.loadSummaryPage(url)
			} else if conf.IsValid(url) {
				url = conf.ToAbsolutePath(url)
				if !conf.InSameSite(url) {
					// utils.DebugLogger.Logf("skip other site url %v", pageUrl)
					return
				}
				go conf.LoadValidPage(url)
			}
		})
	})
	conf.loadPage(pageUrl, []Html.WorkerAction{hrefSelAction})
}

func (conf *SpiderConf) loadPage(url string, actions []Html.WorkerAction) {
	w := Html.New(url, actions)
	if len(conf.Charset) > 0 {
		w.Encoder = Encoding.Encoders[conf.Charset]
	}
	w.OnFail = func(err error) {
		utils.ErrorLogger.Logf("load page %v err %v", url, err)
	}
	w.OnFinish = func() {
		utils.InfoLogger.Logf("load page %v", url)
	}

	downloaderCountChange(true)
	fmt.Println("start : " + url)
	w.Run(func() {
		downloaderCountChange(false)
		fmt.Println("end : " + url)
	})
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
