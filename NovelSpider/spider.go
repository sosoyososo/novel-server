package NovelSpider

import (
	"fmt"

	"../utils"
)

type spiderCtx struct {
	conf       SpiderConf
	loadedUrls []string
}

func StartSummarySpider() {
	for _, confRef := range confMap {
		conf := *confRef
		conf.LoadValidPage(conf.EntryURL)
	}
}

func StartCatelogSpider() {
	for i := 0; ; i++ {
		list, err := ListSummary(i, 20)
		if nil != err {
			utils.ErrorLogger.Logf("%v %v\n", utils.PrintFuncName(), err)
			continue
		}

		for j, _ := range *list {
			s := (*list)[j]
			fmt.Printf("load %v\n", s.Title)
			conf := LoadConf(s.ConfKey)
			if nil == conf {
				continue
			}
			conf.loadCatelog(s.AbsoluteURL, &s)
		}

		if len(*list) < 20 {
			break
		}
	}
}
