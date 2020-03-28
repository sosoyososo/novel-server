package service

import (
	"io/ioutil"
	"reflect"

	"../PO"
	"../utils"
	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

var (
	jsonParser = jsoniter.ConfigCompatibleWithStandardLibrary
)

func init() {
	extra.RegisterFuzzyDecoders()
}

type ServiceCtx struct {
	Ctx *gin.Context
}

func CreateServiceCtx(c *gin.Context) ServiceCtx {
	return ServiceCtx{Ctx: c}
}

func zeroStructV(v1 interface{}) {
	vv := reflect.ValueOf(v1)
	if reflect.TypeOf(v1).Kind() == reflect.Ptr && reflect.Indirect(reflect.ValueOf(v1)).Kind() == reflect.Struct {
		vv = reflect.Indirect(reflect.ValueOf(v1))
	}
	if vv.Kind() != reflect.Struct {
		return
	}
	n := vv.NumField()
	for i := 0; i < n; i++ {
		f := vv.Field(i)
		f.Set(reflect.Zero(f.Type()))
	}

}

func (s ServiceCtx) ParseJsonBody(v interface{}) *PO.ServiceErr {
	zeroStructV(v)
	buf, err := ioutil.ReadAll(s.Ctx.Request.Body)
	if nil != err {
		utils.InfoLogger.Logf("read input fail : %v", err)
		return &PO.Error_WrongParameter
	}
	if len(buf) > 0 {
		err = jsonParser.Unmarshal(buf, v)
		if nil != err {
			utils.InfoLogger.Logf("json input fail : %v \n\t buf: %v\n\t data type: %v",
				err, string(buf), reflect.TypeOf(v))
			return &PO.Error_WrongParameter
		}
	}
	return nil
}
