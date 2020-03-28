package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func NewJsonConfig(confPath string, v interface{}) error {
	f, err := os.Open(confPath)
	if nil != err {
		return err
	}

	buf, err := ioutil.ReadAll(f)
	if nil != err {
		panic(err)
	}

	err = json.Unmarshal(buf, &v)
	if nil != err {
		return err
	}
	return nil
}
