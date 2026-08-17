package app

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"

	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/conf"
	"github.com/james730922/wallet/service/internal/utils/ctxs"
)

func newGinMiddleware() *ginMiddleware {
	return new(ginMiddleware)
}

type ginMiddleware struct {
}

func (ys *ginMiddleware) middlewareGenSelfContext(ctx *gin.Context) {
	self := ctx.Request.Context()
	self = ctxs.SetSessionID(self, uuid.New().String())
	self = ctxs.SetClientIP(self, ctx.ClientIP())
	self = ctxs.SetReqHeaderUserAgent(self, ctx.GetHeader("User-Agent"))
	deviceID := strings.TrimSpace(ctx.GetHeader("X-Device-ID"))
	if len(deviceID) > 256 {
		deviceID = deviceID[:256]
	}
	self = ctxs.SetDeviceID(self, deviceID)

	ctxs.SetSelfContext(self, ctx)

	ctx.Next()
}

func (ys *ginMiddleware) middlewareTraceInOut(ctx *gin.Context) {
	self := ctxs.GetSelfContext(ctx)
	sessionID, _ := ctxs.GetSessionID(self)
	traceID := trace.SpanContextFromContext(self).TraceID().String()
	routing := ctx.Request.Method + " " + ctx.Request.URL.Path
	beginTime := time.Now().UTC()

	logger.AccessLog().Infof("Req: [sessionID: %s], traceID: %s, routing: %s, authenticated: %t, beginTime: %s",
		sessionID,
		traceID,
		routing,
		ctx.GetHeader("Authorization") != "",
		beginTime.String())

	ctx.Next()

	execTime := time.Now().UTC().Sub(beginTime)
	var logFunc func(format string, args ...interface{})
	if execTime > time.Second {
		logFunc = logger.AccessLog().Warnf
	} else {
		logFunc = logger.AccessLog().Infof
	}

	logFunc("Resp: [sessionID: %s], traceID: %s, routing: %s, status: %d, execTime: %s",
		sessionID,
		traceID,
		routing,
		ctx.Writer.Status(),
		execTime.String())
}

func (ys *ginMiddleware) validateIsSimulator(ctx *gin.Context) {
	self := ctxs.GetSelfContext(ctx)
	sessionID, _ := ctxs.GetSessionID(self)
	routing := ctx.Request.Method + " " + ctx.Request.URL.Path

	ua, _ := ctxs.GetReqHeaderUserAgent(self)

	if conf.Blacklist().GetSimulatorEnable() && strings.Contains(ua, "isSimulator") {
		logger.AccessLog().Warnf("It is Simulator: [sessionID: %s], routing: %s, ua: %s",
			sessionID,
			routing,
			ua,
		)
		ctx.AbortWithStatus(http.StatusLocked)
	}

	ctx.Next()
}
