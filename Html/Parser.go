package Html

import (
	"fmt"
	"regexp"

	"../Config"
	"github.com/PuerkitoBio/goquery"
)

type PageParser struct {
	config Config.Config
}

func NewParser(config Config.Config) PageParser {
	p := PageParser{}
	p.config = config
	return p
}

func (p *PageParser) ParseDocument(worker *HtmlWorker) {
	page := handlePageContent(worker.URL, worker.document, p.config)
	fmt.Println("+++========================+++========================+++========================")
	fmt.Printf("%v\n", worker.URL)
	if len(page) > 0 {
		fmt.Println(page)
	}
}

func handlePageContent(url string, document *goquery.Document, config Config.Config) map[string]interface{} {
	pageResult := map[string]interface{}{}

	i := 0
	for i < len(config.Actions) {
		pageAction := config.Actions[i]
		valid := false

		j := 0
		for j < len(pageAction.URLREs) {
			re := pageAction.URLREs[j]
			validID := regexp.MustCompile(re)
			if validID.MatchString(url) {
				valid = true
				break
			}
			j++
		}

		if valid {
			k := 0
			for k < len(pageAction.Actions) {
				action := pageAction.Actions[k]
				result := handleActionOnSelection(action, document.Selection)
				k++
				for k, v := range result {
					pageResult[k] = v
				}
			}
		}

		i++
	}

	return pageResult
}

func handleActionOnSelection(action Config.Action, sel *goquery.Selection) map[string]interface{} {
	pageResult := map[string]interface{}{}

	if len(action.Sel) > 0 {
		sel = sel.Find(action.Sel)
	}
	if len(action.Name) > 0 {
		if len(action.Attr) > 0 {
			content, isExist := sel.Attr(action.Attr)
			if isExist {
				pageResult[action.Name] = content
			}
		} else {
			t := sel.Text()
			pageResult[action.Name] = t
		}
	} else {
		if len(action.Actions) > 0 {
			i := 0
			results := []map[string]interface{}{}
			for i < len(action.Actions) {
				result := handleActionOnSelection(action.Actions[i], sel)
				results = append(results, result)
				i++
			}
			pageResult["arrays"] = results
		}

		if len(action.ListActions) > 0 {
			list := [][]map[string]interface{}{}
			sel.Each(func(index int, s *goquery.Selection) {
				results := []map[string]interface{}{}
				i := 0
				for i < len(action.ListActions) {
					result := handleActionOnSelection(action.ListActions[i], s)
					results = append(results, result)
					i++
				}
				list = append(list, results)
			})
			pageResult["list"] = list
		}
	}

	return pageResult
}
