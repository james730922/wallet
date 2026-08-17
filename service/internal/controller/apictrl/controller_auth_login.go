package apictrl

import (
	"github.com/gin-gonic/gin"

	"github.com/james730922/wallet/service/internal/controller/handler"
	"github.com/james730922/wallet/service/internal/core/base/auth"
	"github.com/james730922/wallet/service/internal/pb/zqbapis"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/ctxs"
)

func newAuthLoginController(core auth.ILoginMember) IAuthLogin {
	return &authLoginController{
		core: core,
	}
}

type IAuthLogin interface {
	// 登出
	IdentifierLogout(ctx *gin.Context)
	// 註冊會員
	Registration(ctx *gin.Context)
	// 密碼登錄
	LoginWithPasswd(ctx *gin.Context)
	// 更新密碼
	UpdatePasswd(ctx *gin.Context)
	// 是否首次使用安全密碼
	FirstUseSecurityPasswd(ctx *gin.Context)
	// 忘記安全密碼
	ForgetSecurityPasswd(ctx *gin.Context)
	// 更新安全密碼
	UpdateSecurityPasswd(ctx *gin.Context)
	// 安全密碼認證，取得安全碼 Token
	SecurityPasswdIdentifier(ctx *gin.Context)
	// 滑動驗證初始化接口
	CaptchaGeeTestRegister(ctx *gin.Context)
}

type authLoginController struct {
	core auth.ILoginMember
}

func (ly *authLoginController) IdentifierLogout(ctx *gin.Context) {
	resp, err := ly.core.IdentifierLogout(ctxs.GetSelfContext(ctx))
	if err != nil {
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	if err := handler.Ctx.ResponseProtoBufStatusOK(ctx, resp); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
	}
}

func (ly *authLoginController) Registration(ctx *gin.Context) {
	var req zqbapis.LoginRegistrationReq

	if err := handler.Ctx.BindProtoBuf(ctx, &req); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	resp, err := ly.core.Registration(ctxs.GetSelfContext(ctx), &req)
	if err != nil {
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	if err := handler.Ctx.ResponseProtoBufStatusOK(ctx, resp); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
	}
}

func (ly *authLoginController) LoginWithPasswd(ctx *gin.Context) {
	var req zqbapis.LoginIdentifierPasswdReq

	if err := handler.Ctx.BindProtoBuf(ctx, &req); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	resp, err := ly.core.LoginWithPasswd(ctxs.GetSelfContext(ctx), &req)
	if err != nil {
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	if err := handler.Ctx.ResponseProtoBufStatusOK(ctx, resp); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
	}
}

func (ly *authLoginController) UpdatePasswd(ctx *gin.Context) {
	var req zqbapis.LoginUpdatePasswdReq

	if err := handler.Ctx.BindProtoBuf(ctx, &req); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	resp, err := ly.core.UpdatePasswd(ctxs.GetSelfContext(ctx), &req)
	if err != nil {
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	if err := handler.Ctx.ResponseProtoBufStatusOK(ctx, resp); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
	}
}

func (ly *authLoginController) FirstUseSecurityPasswd(ctx *gin.Context) {
	resp, err := ly.core.FirstUseSecurityPasswd(ctxs.GetSelfContext(ctx))
	if err != nil {
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	if err := handler.Ctx.ResponseProtoBufStatusOK(ctx, resp); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
	}
}

func (ly *authLoginController) ForgetSecurityPasswd(ctx *gin.Context) {
	var req zqbapis.SecurityPasswdForgetReq

	if err := handler.Ctx.BindProtoBuf(ctx, &req); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	resp, err := ly.core.ForgetSecurityPasswd(ctxs.GetSelfContext(ctx), &req)
	if err != nil {
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	if err := handler.Ctx.ResponseProtoBufStatusOK(ctx, resp); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
	}
}

func (ly *authLoginController) UpdateSecurityPasswd(ctx *gin.Context) {
	var req zqbapis.SecurityPasswdUpdateReq

	if err := handler.Ctx.BindProtoBuf(ctx, &req); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	resp, err := ly.core.UpdateSecurityPasswd(ctxs.GetSelfContext(ctx), &req)
	if err != nil {
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	if err := handler.Ctx.ResponseProtoBufStatusOK(ctx, resp); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
	}
}

func (ly *authLoginController) SecurityPasswdIdentifier(ctx *gin.Context) {
	var req zqbapis.SecurityPasswdIdentifierReq

	if err := handler.Ctx.BindProtoBuf(ctx, &req); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	resp, err := ly.core.SecurityIdentifier(ctxs.GetSelfContext(ctx), &req)
	if err != nil {
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	if err := handler.Ctx.ResponseProtoBufStatusOK(ctx, resp); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
	}
}

func (ly *authLoginController) CaptchaGeeTestRegister(ctx *gin.Context) {
	resp, err := ly.core.CaptchaGeeTestRegister(ctxs.GetSelfContext(ctx))
	if err != nil {
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
		return
	}

	if err := handler.Ctx.ResponseProtoBufStatusOK(ctx, resp); err != nil {
		logger.ApLog().Error(err)
		handler.Ctx.ResponseProtoBufBadRequest(ctx, err)
	}
}
