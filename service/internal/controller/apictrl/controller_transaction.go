package apictrl

import (
	"github.com/gin-gonic/gin"

	"github.com/james730922/wallet/service/internal/controller/handler"
	"github.com/james730922/wallet/service/internal/core/usecase/member/transactionmember"
	"github.com/james730922/wallet/service/internal/pb/zqbapis"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/ctxs"
)

func newTransactionController(core transactionmember.IMember) ITransactionController {
	return &transactionController{
		core: core,
	}
}

type ITransactionController interface {
	List(ctx *gin.Context)
	Detail(ctx *gin.Context)
}

type transactionController struct {
	core transactionmember.IMember
}

func (mc *transactionController) List(ctx *gin.Context) {
	var req zqbapis.TransactionReq
	if err := handler.Ctx.BindProtoBuf(ctx, &req); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	resp, err := mc.core.List(ctxs.GetSelfContext(ctx), &req)
	if err != nil {
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	if err := handler.Ctx.ResponseProtoBufStatusOK(ctx, resp); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
	}
}

func (mc *transactionController) Detail(ctx *gin.Context) {
	id, err := handler.Ctx.GetParamInt64(ctx, "id")
	if err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	resp, err := mc.core.Detail(ctxs.GetSelfContext(ctx), id)
	if err != nil {
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	if err := handler.Ctx.ResponseProtoBufStatusOK(ctx, resp); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
	}
}
