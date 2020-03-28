package utils

import "encoding/json"

func MapToJsonObj(m map[string]interface{}, v interface{}) error {
	buf, err := json.Marshal(m)
	if nil != err {
		return err
	}
	err = json.Unmarshal(buf, v)
	return err
}

func MapListToJsonObjList(m []map[string]interface{}, v interface{}) error {
	buf, err := json.Marshal(m)
	if nil != err {
		return err
	}
	err = json.Unmarshal(buf, v)
	return err
}
