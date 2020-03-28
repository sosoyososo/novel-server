package Config

import (
	"encoding/json"
	"io/ioutil"
)

type Action struct {
	Sel         string  `json:"Sel"`
	Name        string  `json:"Name"`
	Attr        string  `json:"Attr"`
	Actions     Actions `json:"Actions"`
	ListActions Actions `json:"ListActions"`
}
type Actions []Action

type PageAction struct {
	CollectionName string   `json:"CollectionName"`
	URLREs         []string `json:"URLREs"`
	Actions        Actions  `json:"Actions"`
}
type PageActions []PageAction

type Config struct {
	CookiesString string      `json:"cookiesString"`
	EntryURL      string      `json:"entryURL"`
	Encoding      string      `json:"Encoding"`
	SkipUrls      []string    `json:"SkipUrls"`
	Actions       PageActions `json:"PageActions"`
	DbAddr        string      `json:"DbAddr"`
	DbName        string      `json:"DbName"`
}

// FIXME: 深度解析，现在只是解析了一层，后续再解决
func LoadConfig() (*Config, error) {
	raw, err := ioutil.ReadFile("./Config.json")
	if nil != err {
		return nil, err
	}

	var c Config
	err = json.Unmarshal(raw, &c)
	if nil != err {
		return nil, err
	}
	return &c, nil

}
