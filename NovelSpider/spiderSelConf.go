package NovelSpider

import "github.com/PuerkitoBio/goquery"

type NodeSel struct {
	Sel  string `json:"sel"`
	Attr string `json:"attr"`
}

type SelectorConf struct {
	Sel          string             `json:"sel"`
	IsList       bool               `json:"isList"`
	SubSelKeyMap map[string]NodeSel `json:"subSelKeyMap"`
}

func (conf *SelectorConf) ParseConfList(
	sel *goquery.Selection,
) []map[string]interface{} {
	list := []map[string]interface{}{}
	sel.Each(func(_ int, sel *goquery.Selection) {
		ret := conf.ParseConf(sel)
		list = append(list, ret)
	})
	return list
}

func (conf *SelectorConf) ParseConf(
	sel *goquery.Selection,
) map[string]interface{} {
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
