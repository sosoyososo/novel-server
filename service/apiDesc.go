package service

import (
	"fmt"
	"strings"

	"../PO"
	"../utils"
)

var (
	apiItemMap  = map[string]apiItem{}
	baseApiDesc = []string{}
	baseApiPos  = []apiItem{}
	apiDocRet   = ""
)

type apiItem struct {
	Path      string      `json:"request path"`
	Desc      string      `json:"api desc"`
	Do        interface{} `json:"input parameter"`
	Po        interface{} `json:"output parameter"`
	Curl      string      `json:"curl excutor"`
	ExtraInfo string      `json:"extra desc"`
}

func init() {
	baseApiPos = append(baseApiPos, apiItem{
		Desc: "统一返回结构",
		Po:   PO.BasePO{},
	})
	baseApiPos = append(baseApiPos, apiItem{
		Desc:      "分页返回数据字段信息",
		Po:        PO.PageItem{},
		ExtraInfo: "list的内容在具体接口中列出",
	})
	baseApiDesc = []string{
		"两个upload为表单接口，其余默认是JSON API",
	}
}

func (api *apiItem) DocDesc() string {
	return fmt.Sprintf(
		"api desc : \n\t %v\nrequest path : \n\t %v\ninput parameter : \n%v\ncurl action: \n\t%v\noutput parameter : \n%v\n extra :\n\t%v\n",
		api.Desc, api.Path, jsonDesc(api.Do), api.Curl, jsonDesc(api.Po), api.ExtraInfo,
	)
}

func RegisterBaseApiDesc(desc string) {
	baseApiDesc = append(baseApiDesc, desc)
}

func RegisterApiDescPO(path string, po interface{}) {
	item := apiItemMap[path]
	if len(item.Path) > 0 {
		item.Po = po
		apiItemMap[path] = item
	}
}

func RegisterApiDescCURL(path, curl string) {
	item := apiItemMap[path]
	if len(item.Path) > 0 {
		item.Curl = curl
		apiItemMap[path] = item
	}
}

func RegisterApiDescExtra(path, extra string) {
	item := apiItemMap[path]
	if len(item.Path) > 0 {
		item.ExtraInfo = extra
		apiItemMap[path] = item
	}
}

func registerApiInfo(path string, do interface{}, apiDesc string) {
	apiDocRet = ""
	apiItemMap[path] = apiItem{path, apiDesc, do, nil, "", ""}
}

func apiDoc() string {
	tmp := apiDocRet
	if len(tmp) > 0 {
		return tmp
	}
	list := []string{}
	for _, v := range baseApiDesc {
		list = append(list, v)
	}
	for _, v := range baseApiPos {
		list = append(list, v.DocDesc())
	}
	for _, v := range apiItemMap {
		list = append(list, v.DocDesc())
	}
	apiDocRet = strings.Join(
		list,
		"\n\n************************************************************************************\n\n",
	)
	return apiDocRet
}

func jsonDesc(do interface{}) string {
	if nil != do {
		return utils.ModelStructToDocMap(do, "\t")
	}
	return ""
}
