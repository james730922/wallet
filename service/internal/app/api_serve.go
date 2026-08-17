package app

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"

	"github.com/james730922/wallet/service/internal/controller/apictrl"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/thirdparty/observability"
	"github.com/james730922/wallet/service/internal/utils/conf"
)

func NewApiServe(set apiServeSet) IApiServe {
	return &apiServe{apiServeSet: set, middleware: newGinMiddleware()}
}

type IApiServe interface{ Run() }

type apiServe struct {
	apiServeSet
	middleware *ginMiddleware
	engine     *gin.Engine
}

type apiServeSet struct {
	dig.In
	ApiController *apictrl.Controller
	Metrics       *observability.Metrics
	Tracing       *observability.Tracing
}

func (a *apiServe) Run() { go a.start() }

func (a *apiServe) start() {
	a.engine = gin.New()
	if err := a.engine.SetTrustedProxies(conf.Service().GetTrustedProxies()); err != nil {
		logger.SysLog().Panicf("configure trusted proxies failed: %v", err)
	}
	a.globMiddleware()
	a.appRouter()
	address := conf.Service().GetApisHTTPAddress()
	logger.SysLog().Infof("serve start [api], at => %s", address)
	if err := a.engine.Run(address); err != nil {
		logger.SysLog().Errorf("api serve stopped: %v", err)
	}
}

func (a *apiServe) globMiddleware() {
	a.engine.Use(a.Tracing.Middleware())
	a.engine.Use(a.Metrics.Middleware())
	if conf.Log().GetGinLogEnable() {
		a.engine.Use(gin.Logger())
	}
	a.engine.Use(gin.Recovery())
	a.engine.Use(a.middleware.middlewareGenSelfContext)
	a.engine.Use(a.middleware.middlewareTraceInOut)
}

func (a *apiServe) appRouter() {
	a.engine.GET("/livez", a.Metrics.Livez)
	a.engine.GET("/readyz", a.Metrics.Readyz)
	a.engine.GET("/metrics", gin.WrapH(a.Metrics.Handler()))

	anonymous := a.engine.Group("/api")
	anonymous.Use(a.middleware.validateIsSimulator)
	anonymous.POST("/v1/registration", a.ApiController.AuthLogin.Registration)
	anonymous.POST("/v1/login", a.ApiController.AuthLogin.LoginWithPasswd)
	anonymous.POST("/v1/captcha/register", a.ApiController.AuthLogin.CaptchaGeeTestRegister)

	apis := a.engine.Group("/api")
	apis.Use(a.middleware.validateIsSimulator)
	apis.Use(a.ApiController.HttpAuthToken.MiddlewareAuth)
	apis.POST("/v1/logout", a.ApiController.AuthLogin.IdentifierLogout)
	apis.POST("/v1/update-pwd", a.ApiController.AuthLogin.UpdatePasswd)
	apis.GET("/v1/security-passwd/first", a.ApiController.AuthLogin.FirstUseSecurityPasswd)
	apis.POST("/v1/security-passwd/forget", a.ApiController.AuthLogin.ForgetSecurityPasswd)
	apis.POST("/v1/security-passwd/update", a.ApiController.AuthLogin.UpdateSecurityPasswd)
	apis.POST("/v1/security-passwd/identifier", a.ApiController.AuthLogin.SecurityPasswdIdentifier)
	apis.GET("/v1/banks", a.ApiController.Bank.List)
	apis.GET("/v1/deposit", a.ApiController.Deposit.List)
	apis.GET("/v1/deposit/methods", a.ApiController.Deposit.Methods)
	apis.POST("/v1/deposit/order", a.ApiController.Deposit.Order)
	apis.GET("/v1/transaction", a.ApiController.Transaction.List)
	apis.GET("/v1/transaction/detail/:id", a.ApiController.Transaction.Detail)
	apis.POST("/v1/scan-pay/scan", a.ApiController.ScanPay.Scan)
	apis.POST("/v1/scan-pay/pay", a.ApiController.ScanPay.Pay)
	apis.POST("/v1/scan-pay/qrcode", a.ApiController.ScanPay.QRCode)
}
