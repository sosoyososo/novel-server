package NovelSpider

import (
	"../Encoding"
	"../Html"
	"../utils"
	"github.com/PuerkitoBio/goquery"
)

//TODO: all these func should be in handle queue

func (conf *SpiderConf) loadSummaryPage(pageUrl string) {
	if isSummaryLoaded(pageUrl) {
		return
	}
	markSummaryLoaded(pageUrl)

	utils.InfoLogger.Logf("hit summary url %v", pageUrl)
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
		c, err := CatelogNumOfSummary(s.ID)
		if nil != err {
			utils.ErrorLogger.Logf("%v %v\n", utils.PrintFuncName(), err)
			return
		}
		if c >= len(ret) {
			return
		}
		needCheck := true
		if c == 0 {
			needCheck = false
		}

		var list []CatelogInfo
		err = utils.MapListToJsonObjList(ret, &list)
		if nil != err {
			utils.ErrorLogger.Logf("%v %v\n", utils.PrintFuncName(), err)
			return
		}

		for _, c := range list {
			if needCheck {
				hasLoaded, err := CatelogWithDetailURLHasLoaded(c.DetailURL)
				if nil != err {
					utils.ErrorLogger.Logf("%v %v\n", utils.PrintFuncName(), err)
					continue
				}
				if hasLoaded {
					continue
				}
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
				return
			}

			if conf.HasLoaded(url) {
				return
			}
			conf.MarkLoaded(url)

			if conf.IsSummaryPage(url) {
				url = conf.ToAbsolutePath(url)
				if !conf.InSameSite(url) {
					utils.InfoLogger.Logf("skip other site url %v", pageUrl)
					return
				}
				conf.loadSummaryPage(url)
			} else if conf.IsValid(url) {
				url = conf.ToAbsolutePath(url)
				if !conf.InSameSite(url) {
					utils.InfoLogger.Logf("skip other site url %v", pageUrl)
					return
				}
				conf.LoadValidPage(url)
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
	w.Run()
}
