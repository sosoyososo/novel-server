package utils

import (
	"fmt"
	"os"
)

var (
	conf GlobalConf
)

type Env int

const (
	EnvNone Env = iota
	EnvRelease
	EnvRaspiDebug
)

const (
	env = EnvRelease
)

type WxConfig struct {
	AppID      string `json:"appid"`
	Secret     string `json:"secret"`
	EncryptKey string `json:"encryptKey"`
	KeyPath    string `json:"keyPath"`
	CertPath   string `json:"certPath"`
}

type ReadPocketConf struct {
	CreateFee          float32 `json:"createFee"`
	DrawCashFee        float32 `json:"drawCashFee"`
	DrawCashCheckLimit int     `json:"drawCashCheckLimit"`
}

type ServiceConf struct {
	Port         string `json:"port"`
	CookieDomain string `json:"cookieDomain"`
}

type GlobalConf struct {
	Server           *ServiceConf         `json:"server"`
	RedPocketConf    *ReadPocketConf      `json:"redPocketConf"`
	UserDbConf       *DBConfig            `json:"userDB"`
	RedPocketDbConf  *DBConfig            `json:"redPocketDB"`
	CommonDbConf     *DBConfig            `json:"commonDB"`
	BidDbConf        *DBConfig            `json:"bidingDB"`
	RedisConf        *RedisConf           `json:"redis"`
	WxDefaultAppName string               `json:"wxDefaultAppName"`
	MiniAppList      *map[string]WxConfig `json:"wxList"`
	WxPay            *WxConfig            `json:"wxPay"`
}

func InitDB(conf *DBConfig, list []interface{}) (*DBTools, error) {
	db, err := conf.CreateDB()
	if nil != err {
		panic(err)
	}

	updateTables := func(db *DBTools, modelList []interface{}) {
		if ShouldUpdateDB() {
			for index, _ := range modelList {
				v := modelList[index]
				if !db.HasTable(v) {
					db.CreateTable(v)
				}
				//新建的表，同时在新增字段的时候保持数据库表的更新
				// 不会修改或者删除已有字段
				db.AutoMigrate(v)
			}
		}
	}
	updateTables(db, list)

	return db, nil
}

func InitConf() {
	l := os.Args
	fmt.Println(l)

	InfoLogger.Logln("start load conf .......")
	err := NewJsonConfig(configPath(), &conf)
	if nil != err {
		panic(err)
	}

	if nil == conf.Server || len(conf.Server.Port) == 0 {
		panic("没有配置服务端口")
	}
	InfoLogger.Logln("load conf finished !")
}

func GetPathRelativeToProjRoot(path string) string {
	p := os.Getenv("projRoot")
	if len(p) <= 0 {
		p = "."
	}
	return p + "/" + path
}

func configPath() string {
	path := "conf.json"
	if env == EnvRaspiDebug {
		path = "conf-raspi.json"
	}
	return GetPathRelativeToProjRoot(path)
}

func GetConf() GlobalConf {
	return conf
}
