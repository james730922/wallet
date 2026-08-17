package ctxs

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

type keyName int

const (
	selfContext = "zqb_self_context"
)

const (
	memberID keyName = iota
	sessionID
	merchantID
	tokenID
	keyReqClientIP
	keyReqHeaderUserAgent
	keyReqDeviceID
)

func New() context.Context {
	return context.Background()
}

func SetSelfContext(ctx context.Context, sc interface{}) interface{} {
	switch t := sc.(type) {
	//case *ydras.Context:
	//	t.SetValue(selfContext, ctx)
	case *gin.Context:
		t.Set(selfContext, ctx)
	default:
		logger.ApLog().Panicf("%s, type: %v", errs.FrameworkContextErrorType, t)
	}

	return sc
}

func GetSelfContext(ctx interface{}) context.Context {
	var iSelf interface{}

	switch t := ctx.(type) {
	//case *ydras.Context:
	//	if v, ok := t.GetValue(selfContext); ok {
	//		iSelf = v
	//	}
	case *gin.Context:
		if v, ok := t.Get(selfContext); ok {
			iSelf = v
		}
	default:
		logger.ApLog().Panicf("%s, type: %v", errs.FrameworkContextErrorType, t)
	}

	if c, ok := iSelf.(context.Context); ok {
		return c
	}

	return New()
}

func SetSessionID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, sessionID, id)
}

func GetSessionID(ctx context.Context) (string, bool) {
	val := ctx.Value(sessionID)

	v, ok := val.(string)
	if !ok {
		return "", false
	}

	return v, true
}

func SetMemberID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, memberID, id)
}

func GetMemberID(ctx context.Context) (int64, bool) {
	val := ctx.Value(memberID)

	v, ok := val.(int64)
	if !ok {
		return 0, false
	}

	return v, true
}

func SetMerchantID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, merchantID, id)
}

func GetMerchantID(ctx context.Context) (string, bool) {
	val := ctx.Value(merchantID)

	v, ok := val.(string)
	if !ok {
		return "", false
	}
	return v, true
}

func SetToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenID, token)
}

func GetToken(ctx context.Context) (string, bool) {
	val := ctx.Value(tokenID)

	v, ok := val.(string)
	if !ok {
		return "", false
	}

	return v, true
}

func SetClientIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, keyReqClientIP, ip)
}

func GetClientIP(ctx context.Context) (string, bool) {
	val := ctx.Value(keyReqClientIP)
	v, ok := val.(string)
	if !ok {
		return "", false
	}
	return v, true
}

func SetReqHeaderUserAgent(ctx context.Context, userAgent string) context.Context {
	return context.WithValue(ctx, keyReqHeaderUserAgent, userAgent)
}

func GetReqHeaderUserAgent(ctx context.Context) (string, bool) {
	val := ctx.Value(keyReqHeaderUserAgent)

	v, ok := val.(string)
	if !ok {
		return "", false
	}

	return v, true
}

func SetDeviceID(ctx context.Context, deviceID string) context.Context {
	return context.WithValue(ctx, keyReqDeviceID, deviceID)
}

func GetDeviceID(ctx context.Context) (string, bool) {
	val := ctx.Value(keyReqDeviceID)
	v, ok := val.(string)
	if !ok {
		return "", false
	}
	return v, true
}
