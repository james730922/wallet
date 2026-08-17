package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/go-redis/redis/v7"

	"github.com/james730922/wallet/service/internal/utils/conf"
	"github.com/james730922/wallet/service/internal/utils/ctxs"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

const reserveRegistrationAttemptScript = `
local ip_count = tonumber(redis.call("GET", KEYS[1]) or "0")
local device_count = tonumber(redis.call("GET", KEYS[2]) or "0")
if ip_count >= tonumber(ARGV[2]) or device_count >= tonumber(ARGV[3]) then
  return 0
end

local new_ip_count = redis.call("INCR", KEYS[1])
if new_ip_count == 1 or redis.call("TTL", KEYS[1]) < 0 then
  redis.call("EXPIRE", KEYS[1], ARGV[1])
end
local new_device_count = redis.call("INCR", KEYS[2])
if new_device_count == 1 or redis.call("TTL", KEYS[2]) < 0 then
  redis.call("EXPIRE", KEYS[2], ARGV[1])
end
return 1`

type registrationGuard struct {
	redis     *redis.Client
	window    time.Duration
	ipMax     int64
	deviceMax int64
}

func newRegistrationGuard(client *redis.Client) *registrationGuard {
	loginConf := conf.LoginMember()
	return &registrationGuard{
		redis:     client,
		window:    loginConf.GetRegistrationAttemptWindow(),
		ipMax:     loginConf.GetRegistrationIPMaxAttempts(),
		deviceMax: loginConf.GetRegistrationDeviceMaxAttempts(),
	}
}

// beforeAttempt 在 captcha 關閉時仍強制 IP/device 限流，且 Redis 異常時 fail closed。
func (g *registrationGuard) beforeAttempt(ctx context.Context) error {
	clientIP, _ := ctxs.GetClientIP(ctx)
	deviceID, _ := ctxs.GetDeviceID(ctx)
	if deviceID == "" {
		deviceID, _ = ctxs.GetReqHeaderUserAgent(ctx)
	}
	ipKey, deviceKey := registrationAttemptKeys(clientIP, deviceID)
	seconds := int64(g.window / time.Second)
	if seconds < 1 {
		seconds = 1
	}

	result, err := g.redis.Eval(
		reserveRegistrationAttemptScript,
		[]string{ipKey, deviceKey},
		seconds,
		g.ipMax,
		g.deviceMax,
	).Int64()
	if err != nil {
		return errs.CommonServiceUnavailable
	}
	if result != 1 {
		return errs.CommonFrequentOperationError
	}
	return nil
}

func registrationAttemptKeys(clientIP, deviceID string) (string, string) {
	if clientIP == "" {
		clientIP = "unknown"
	}
	if deviceID == "" {
		deviceID = "unknown"
	}
	ipDigest := sha256.Sum256([]byte(clientIP))
	deviceDigest := sha256.Sum256([]byte(clientIP + "\x00" + deviceID))
	return "auth:registration:ip:" + hex.EncodeToString(ipDigest[:]),
		"auth:registration:device:" + hex.EncodeToString(deviceDigest[:])
}
