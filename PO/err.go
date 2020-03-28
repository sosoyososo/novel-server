package PO

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type ServiceErr struct {
	Code int    `json:"code" comment:"错误码"`
	Msg  string `json:"msg" comment:"错误提示信息"`
}

func (s ServiceErr) RendFail(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		BasePO{
			Data:    nil,
			Succeed: false,
			Code:    s.Code,
			Message: s.Msg,
		},
	)
}

func (s *ServiceErr) Error() string {
	return s.Msg
}

func NormalErr(msg string) *ServiceErr {
	return &ServiceErr{Code: Code_Fail, Msg: msg}
}

func ServiceErrUnWrapper(err error, defaultErr ServiceErr) *ServiceErr {
	fmt.Println(err)
	if reflect.TypeOf(err) == reflect.TypeOf(&ServiceErr{}) {
		v := err.(*ServiceErr)
		return v
	} else {
		return &defaultErr
	}
}
