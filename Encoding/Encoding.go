package Encoding

import (
	"bytes"
	"io/ioutil"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

/*
封装一个工具，用来将一个网页转换成goquery可以使用的utf8格式
*/

var (
	Encoders = map[string]func([]byte) ([]byte, error){
		"gbk":    GbkToUtf8,
		"gb2312": Gb2312ToUtf8}
)

/*
HtmlContentEncoding 是一个接口
*/
type HtmlContentEncoding interface {
	ConvertToUtf8(s []byte) ([]byte, error)
}

/*
GbkToUtf8 从 gbk 转码为 utf8
*/
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func Gb2312ToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.HZGB2312.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
