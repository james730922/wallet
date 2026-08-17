package service

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"

	"github.com/james730922/wallet/service/internal/app"
	"github.com/james730922/wallet/service/internal/binder"
	"github.com/james730922/wallet/service/internal/core/base/scanpay"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/conf"
)

func newServer() IServer {
	return &server{}
}

type server struct {
}

func (s server) Run() {
	s.preStart()
	s.ginMode()

	logger.SysLog().Infof("[Build Info] %s", getBuildInfo())

	binder := binder.New()
	if err := binder.Invoke(s.gen); err != nil {
		panic(err)
	}
}

func (s server) preStart() {
	conf.Start()
	logger.Start()
}

func (s server) ginMode() {
	gin.DisableConsoleColor()
	switch conf.Env() {
	case conf.EnvTypeProd:
		gin.SetMode(gin.ReleaseMode)
	case conf.EnvTypeDev:
		gin.SetMode(gin.TestMode)
	}
}

func (s server) gen(set serveSet) {
	go set.PyroscopeServe.Run()
	go set.ApiServe.Run()
	go set.ScanPayReconciler.Run(context.Background())
}

type serveSet struct {
	dig.In

	ApiServe          app.IApiServe
	PyroscopeServe    app.IPyroscopeServe
	ScanPayReconciler scanpay.IReconciler
}
