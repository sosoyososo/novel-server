package service

import (
	"fmt"
	"net/http"

	"../PO"
	"../utils"
	"github.com/gin-gonic/gin"
)

var (
	ginEngine *gin.Engine = nil
)

func init() {
	ginEngine = gin.New()
	configEngin()
}

func configEngin() {
	ginEngine.Use(gin.Recovery())
	ginEngine.Use(gin.Logger())
	ginEngine.Use(corsMiddleWare())
	ginEngine.Use(authMiddleWare())

	RegisterAuthNoNeedPath("/")
	ginEngine.Any("/", func(ctx *gin.Context) {
		PO.RendSucceedData(ctx, "Succeed!")
	})

	RegisterAuthNoNeedPath("/doc")
	ginEngine.GET("/doc", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, apiDoc())
	})
}

func CreateService() {
	port := utils.GetConf().Server.Port
	fmt.Println("start server on " + port)
	err := ginEngine.Run(port)
	if nil != err {
		fmt.Println(err)
	}
}

func RegisterRawService(path, apiDesc string, handle func(*gin.Context)) {
	registerApiInfo(path, nil, apiDesc)
	ginEngine.Any(path, handle)
}

func RegisterJSONServiceV2(
	path string,
	po interface{},
	handle func(ctx ServiceCtx) (interface{}, error),
	apiDesc string,
) {
	registerApiInfo(path, po, apiDesc)
	ginEngine.Any(path, func(c *gin.Context) {
		ctx := CreateServiceCtx(c)
		err := ctx.ParseJsonBody(po)
		if nil != err {
			err.RendFail(c)
			return
		}

		err = nil
		data, normalErr := handle(ctx)
		if nil != normalErr {
			if poErr, ok := normalErr.(*PO.ServiceErr); ok {
				err = poErr
			} else {
				err = PO.NormalErr(normalErr.Error())
			}
		}

		if nil != err {
			err.RendFail(c)
		} else {
			PO.RendSucceedData(c, data)
		}
	})
}

func RegisterJSONService(
	path string,
	po interface{},
	handle func(ctx ServiceCtx) (interface{}, *PO.ServiceErr),
	apiDesc string,
) {
	registerApiInfo(path, po, apiDesc)
	ginEngine.Any(path, func(c *gin.Context) {
		ctx := CreateServiceCtx(c)
		err := ctx.ParseJsonBody(po)
		if nil != err {
			err.RendFail(c)
			return
		}
		data, err := handle(ctx)
		if nil != err {
			err.RendFail(c)
		} else {
			PO.RendSucceedData(c, data)
		}
	})
}

func RegisterListJSONServiceV2(
	path string, po interface{},
	handle func(ctx ServiceCtx) (interface{}, int, error),
	apiDesc string,
) {
	registerApiInfo(path, po, apiDesc)
	ginEngine.Any(path, func(c *gin.Context) {
		ctx := CreateServiceCtx(c)
		err := ctx.ParseJsonBody(po)
		if nil != err {
			err.RendFail(c)
			return
		}

		err = nil
		list, total, normalErr := handle(CreateServiceCtx(c))
		if nil != normalErr {
			if poErr, ok := normalErr.(*PO.ServiceErr); ok {
				err = poErr
			} else {
				err = PO.NormalErr(normalErr.Error())
			}
		}

		if nil != err {
			err.RendFail(c)
		} else {
			PO.RendSucceedList(c, list, total)
		}
	})
}

func RegisterListJSONService(path string, po interface{}, handle func(ctx ServiceCtx) (interface{}, int, *PO.ServiceErr), apiDesc string) {
	registerApiInfo(path, po, apiDesc)
	ginEngine.Any(path, func(c *gin.Context) {
		ctx := CreateServiceCtx(c)
		err := ctx.ParseJsonBody(po)
		if nil != err {
			err.RendFail(c)
			return
		}
		list, total, err := handle(CreateServiceCtx(c))
		if nil != err {
			err.RendFail(c)
		} else {
			PO.RendSucceedList(c, list, total)
		}
	})
}
