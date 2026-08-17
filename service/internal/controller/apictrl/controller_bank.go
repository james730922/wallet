package apictrl

import (
	"github.com/gin-gonic/gin"

	"github.com/james730922/wallet/service/internal/controller/handler"
	"github.com/james730922/wallet/service/internal/core/base/bank"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/pb/zqbapis"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
)

func newBankController(core bank.IBankCommon) IBankController {
	return &bankController{
		core: core,
	}
}

type IBankController interface {
	List(ctx *gin.Context)
}

type bankController struct {
	core bank.IBankCommon
}

func (bk *bankController) List(ctx *gin.Context) {
	banks, err := bk.core.List(ctx, &condition.BankCodeQuery{})
	if err != nil {
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	resp := bk.arrangeResp(banks)

	if err := handler.Ctx.ResponseProtoBufStatusOK(ctx, resp); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
	}
}

func (bk *bankController) arrangeResp(banks []*model.BankCode) *zqbapis.BankListResp {
	data := []*zqbapis.BankList{}

	for _, v := range banks {
		data = append(data, &zqbapis.BankList{
			Code:  v.Code,
			Name:  v.Name,
			Image: v.Image,
		})
	}

	return &zqbapis.BankListResp{Data: data}

}
