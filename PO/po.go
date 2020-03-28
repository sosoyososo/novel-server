package PO

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BasePO struct {
	Data    interface{} `json:"data,omitempty" comment:"数据"`
	Succeed bool        `json:"succeed" comennt:"是否成功"`
	Code    int         `json:"code" comment:"错误码"`
	Message string      `json:"message" comment:"错误提示信息"`
}

type PageItem struct {
	List  interface{} `json:"list" comment:"列表"`
	Total int         `json:"total" comment:"总的数据条数"`
}

func RendSucceedData(c *gin.Context, data interface{}) {
	c.JSON(
		http.StatusOK,
		BasePO{
			Data:    data,
			Succeed: true,
			Code:    Succeed_Succeed.Code,
			Message: Succeed_Succeed.Msg,
		},
	)
}

func RendSucceedList(c *gin.Context, list interface{}, total int) {
	c.JSON(
		http.StatusOK,
		BasePO{
			Succeed: true,
			Code:    Succeed_Succeed.Code,
			Message: Succeed_Succeed.Msg,
			Data:    PageItem{list, total},
		},
	)
}
