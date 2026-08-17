package apictrl

import (
	"github.com/gin-gonic/gin"

	"github.com/james730922/wallet/service/internal/controller/handler"
	"github.com/james730922/wallet/service/internal/core/usecase/member/scanpaymember"
	"github.com/james730922/wallet/service/internal/pb/zqbapis"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/ctxs"
)

func newScanPayController(core scanpaymember.IMember) IScanPayController {
	return &scanPayController{
		core: core,
	}
}

type IScanPayController interface {
	Scan(ctx *gin.Context)
	Pay(ctx *gin.Context)
	QRCode(ctx *gin.Context)
}

type scanPayController struct {
	core scanpaymember.IMember
}

func (ac *scanPayController) Scan(ctx *gin.Context) {
	var req zqbapis.ScanPayScanReq
	if err := handler.Ctx.BindProtoBuf(ctx, &req); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	resp, err := ac.core.Scan(ctxs.GetSelfContext(ctx), &req)
	if err != nil {
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	if err := handler.Ctx.ResponseProtoBufStatusOK(ctx, resp); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
	}
}

func (ac *scanPayController) Pay(ctx *gin.Context) {
	var req zqbapis.ScanPayAddPayReq
	if err := handler.Ctx.BindProtoBuf(ctx, &req); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	resp, err := ac.core.Pay(ctxs.GetSelfContext(ctx), &req)
	if err != nil {
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	if err := handler.Ctx.ResponseProtoBufStatusOK(ctx, resp); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
	}
}

func (ac *scanPayController) QRCode(ctx *gin.Context) {
	var req zqbapis.ScanPayScanReq
	if err := handler.Ctx.BindProtoBuf(ctx, &req); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	resp, err := ac.core.UploadQRCode(ctxs.GetSelfContext(ctx), &req)
	if err != nil {
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	if err := handler.Ctx.ResponseProtoBufStatusOK(ctx, resp); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
	}
}
