package api

import (
	"hrserver/types"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (h *Handlers) HandlerJobList(c *gin.Context) {
	var (
		err error
		rsp *types.RspJoblist
	)
	req := new(types.ReqJoblist)

	//数据绑定并解析到req结构体
	if err = c.ShouldBindWith(req, binding.JSON); err != nil {
		goto errDeal
	}
	rsp = new(types.RspJoblist)
	HandlerSuccess(c, "HandlerPost", rsp)
	return
errDeal:
	HandlerFailed(c, "HandlerPost", err.Error())
	return
}
