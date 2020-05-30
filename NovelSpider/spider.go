package NovelSpider

import (
	"../timerTask"
)

type spiderCtx struct {
	conf       SpiderConf
	loadedUrls []string
}

func initTask() {
	registerSummarySpider()
}

func StartSummarySpider() {
	for _, confRef := range confMap {
		conf := *confRef
		conf.LoadValidPage(conf.EntryURL)
	}
}

func registerSummarySpider() {
	timerTask.RegisterSlow2Repeat("registerSummarySpider", func(d interface{}) timerTask.TaskHandleResultType {
		StartSummarySpider()
		return timerTask.TaskHandleResultTypeNone
	}, "")
}
