package apictrl

import (
	"github.com/gin-gonic/gin"

	"github.com/james730922/wallet/service/internal/controller/handler"
	"github.com/james730922/wallet/service/internal/core/base/deposit"
	"github.com/james730922/wallet/service/internal/pb/zqbapis"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/ctxs"
)

func newDepositController(core deposit.IDepositMember) IDepositController {
	return &depositController{
		core: core,
	}
}

type IDepositController interface {
	List(ctx *gin.Context)
	Methods(ctx *gin.Context)
	Order(ctx *gin.Context)
}

type depositController struct {
	core deposit.IDepositMember
}

func (dy *depositController) List(ctx *gin.Context) {
	var req zqbapis.DepositListReq
	if err := handler.Ctx.BindProtoBuf(ctx, &req); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}
	resp, err := dy.core.List(ctxs.GetSelfContext(ctx), &req)
	if err != nil {
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	if err := handler.Ctx.ResponseProtoBufStatusOK(ctx, resp); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
	}
}

func (dy *depositController) Methods(ctx *gin.Context) {
	resp, err := dy.core.Methods(ctxs.GetSelfContext(ctx))
	if err != nil {
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	if err := handler.Ctx.ResponseProtoBufStatusOK(ctx, resp); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
	}
}

func (dy *depositController) Order(ctx *gin.Context) {
	var req zqbapis.DepositOrderReq
	if err := handler.Ctx.BindProtoBuf(ctx, &req); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	resp, err := dy.core.Order(ctxs.GetSelfContext(ctx), &req)
	if err != nil {
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	if err := handler.Ctx.ResponseProtoBufStatusOK(ctx, resp); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
	}
}
