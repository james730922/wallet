package apictrl

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/james730922/wallet/service/internal/controller/handler"
	"github.com/james730922/wallet/service/internal/core/base/auth"
	"github.com/james730922/wallet/service/internal/utils/ctxs"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newAuthTokenHttpController(token auth.IToken) IAuthTokenHttpController {
	return &authTokenHttpController{
		token: token,
	}
}

type IAuthTokenHttpController interface {
	MiddlewareAuth(ctx *gin.Context)
	MiddlewareAuthMemberSecurityPasswd(ctx *gin.Context)
}

type authTokenHttpController struct {
	token auth.IToken
}

func (at *authTokenHttpController) MiddlewareAuth(ctx *gin.Context) {
	sc := ctxs.GetSelfContext(ctx)

	token := ctx.GetHeader("Authorization")

	memberRepository, err := at.token.AuthMember(sc, token)
	if err != nil {
		handler.Ctx.Response(
			ctx,
			http.StatusUnauthorized,
			handler.ContentTypeProtoBuf,
			handler.ErrorHandler.Marshal(errs.AuthTokenUnauthorized),
		)
		ctx.Abort()
		return
	}

	sc = ctxs.SetMemberID(sc, memberRepository.ID)
	sc = ctxs.SetToken(sc, memberRepository.Token)
	ctx = ctxs.SetSelfContext(sc, ctx).(*gin.Context)

	ctx.Next()
}

func (at *authTokenHttpController) MiddlewareAuthMemberSecurityPasswd(ctx *gin.Context) {
	sc := ctxs.GetSelfContext(ctx)

	token := ctx.GetHeader("X-Security-Authorization")

	_, err := at.token.AuthMemberSecurityPasswd(sc, token)
	if err != nil {
		handler.Ctx.Response(
			ctx,
			http.StatusForbidden,
			handler.ContentTypeProtoBuf,
			handler.ErrorHandler.Marshal(errs.AuthMemberSecurityPasswdTokenUnauthorized),
		)
		ctx.Abort()
		return
	}

	ctx.Next()
}
