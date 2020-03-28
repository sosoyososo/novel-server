package NovelSpider

import (
	"fmt"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"../utils"
)

func TestSummarySpider(t *testing.T) {
}

func TestSummaryLoadPage(t *testing.T) {
	confRef := confMap["biquge"]
	conf := *confRef
	conf.loadSummaryPage("https://www.biqubu.com/book_1064/")
}

func TestGOQuery(t *testing.T) {
	f, err := os.Open("/Users/karsa/Desktop/1.htm")
	if nil != err {
		t.Error(err)
	}
	doc, err := goquery.NewDocumentFromReader(f)
	if nil != err {
		t.Error(err)
	}

	var confs map[string]SpiderConf
	err = utils.NewJsonConfig(utils.GetPathRelativeToProjRoot("./selConf.json"), &confs)
	if nil != err {
		t.Error(err)
	}

	conf := confs["biquge"]

	sel := doc.Find(conf.SummarySelectorConf.Sel)
	for key, node := range conf.SummarySelectorConf.SubSelKeyMap {
		fmt.Printf("%v :  %v\n", key, sel.Find(node.Sel).First().Text())
	}
}

func TestDB(t *testing.T) {
	// defaultDB.Model(&CatelogInfo{}).ModifyColumn("description", "text")
	// fmt.Println(len(loadedSummary))
	// sL, err := ListSummary(0, 100)
	// if nil != err {
	// 	t.Error(err)
	// }
	// fmt.Println(sL)

	// list, err := ListCatelog(0, 10)
	// if nil != err {
	// 	t.Error(err)
	// }
	// fmt.Println(list)
}

func TestRegexp(t *testing.T) {
	match, err := regexpMatch("^/book_[^/]*/.+", "/book_30856/111.html")
	if nil != err {
		t.Error(err)
	}
	fmt.Println(match)
}

func TestLoadAddtionalCatelog(t *testing.T) {
	confRef := confMap["biquge"]
	conf := *confRef
	var s Summary
	err := defaultDB.Model(Summary{}).Where("absolute_url = ?", "https://www.biqubu.com/book_22078/").Scan(&s).Error
	if nil != err {
		t.Error(err)
	}
	conf.loadCatelog(s.AbsoluteURL, &s)
}

func TestUpdateBaseModel(t *testing.T) {
	for page := 0; ; page++ {
		list, err := ListSummary(page, 20)
		if nil != err {
			t.Error(err)
		}
		if len(*list) == 0 {
			break
		}
		for i, _ := range *list {
			s := (*list)[i]
			s.initBase()
			if len(s.AbsoluteURL) == 0 {
				continue
			}
			defaultDB.Model(Summary{}).Where("absolute_url = ?", s.AbsoluteURL).Update(s)
		}
	}
}

func TestLoadDetail(t *testing.T) {
	var confs map[string]SpiderConf
	err := utils.NewJsonConfig(utils.GetPathRelativeToProjRoot("./selConf.json"), &confs)
	if nil != err {
		panic(err)
	}

	conf := confs["biquge"]
	conf.loadSummaryPage("https://www.biqubu.com/book_656/")
}
