package api

import (
	"fmt"

	"hrserver/types"
	"hrserver/util/log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Handlers struct{}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) OnClose() {}

func (h *Handlers) UserAuthrize(c *gin.Context) {
	var (
		err error
	)
	//权限过滤器，如果验证成功继续，失败则goto errDeal
	if false {
		goto errDeal
	}
	c.Next()
	return
errDeal:
	HandlerFailed(c, "UserAuthrize", err.Error())
	c.Abort()
	return
}

func (h *Handlers) HandlerGet(c *gin.Context) {
	var (
		err error
	)
	if err != nil {
		goto errDeal
	}
	HandlerSuccess(c, "HandlerGet", "hello world")
	return
errDeal:
	HandlerFailed(c, "HandlerGet", err.Error())
	return
}

func (h *Handlers) HandlerPost(c *gin.Context) {
	var (
		err error
		rsp *types.RspPost
	)
	req := new(types.RepPost)

	//数据绑定并解析到req结构体
	if err = c.ShouldBindWith(req, binding.JSON); err != nil {
		goto errDeal
	}
	rsp = new(types.RspPost)
	rsp.Name = req.Name
	HandlerSuccess(c, "HandlerPost", rsp)
	return
errDeal:
	HandlerFailed(c, "HandlerPost", err.Error())
	return
}

func HandlerSuccess(c *gin.Context, requestType, data interface{}) {
	c.JSON(200, gin.H{
		"isSuccess": true,
		"message":   "success",
		"data":      data,
	})
	logMsg := fmt.Sprintf("From [%s] result success", c.Request.RemoteAddr)
	log.GetLog().LogInfo(requestType, logMsg)
}

func HandlerFailed(c *gin.Context, requestType, errMsg string) {
	c.JSON(200, gin.H{
		"isSuccess": false,
		"message":   errMsg,
	})
	logMsg := fmt.Sprintf("From [%s] result error [%s]", c.Request.RemoteAddr, errMsg)
	log.GetLog().LogError(requestType, logMsg)
}
