package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"

	"github.com/james730922/wallet/service/internal/utils/conf"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

const reserveLoginAttemptScript = `
local account_count = tonumber(redis.call("GET", KEYS[1]) or "0")
local ip_count = tonumber(redis.call("GET", KEYS[2]) or "0")
if account_count >= tonumber(ARGV[2]) or ip_count >= tonumber(ARGV[3]) then
  return {0, account_count, ip_count}
end

local new_account_count = redis.call("INCR", KEYS[1])
if new_account_count == 1 or redis.call("TTL", KEYS[1]) < 0 then
  redis.call("EXPIRE", KEYS[1], ARGV[1])
end
local new_ip_count = redis.call("INCR", KEYS[2])
if new_ip_count == 1 or redis.call("TTL", KEYS[2]) < 0 then
  redis.call("EXPIRE", KEYS[2], ARGV[1])
end
return {1, account_count, ip_count}`

type loginGuard struct {
	redis         *redis.Client
	attemptWindow time.Duration
	accountMax    int64
	ipMax         int64
	captchaAfter  int64
}

func newLoginGuard(client *redis.Client) *loginGuard {
	loginConf := conf.LoginMember()
	return &loginGuard{
		redis:         client,
		attemptWindow: loginConf.GetLoginAttemptWindow(),
		accountMax:    loginConf.GetLoginAccountMaxAttempts(),
		ipMax:         loginConf.GetLoginIPMaxAttempts(),
		captchaAfter:  loginConf.GetLoginCaptchaAfterAttempts(),
	}
}

func (g *loginGuard) beforeAttempt(countryCode, mobile, clientIP string) (bool, error) {
	accountKey, ipKey := loginAttemptKeys(countryCode, mobile, clientIP)
	seconds := int64(g.attemptWindow / time.Second)
	if seconds < 1 {
		seconds = 1
	}
	result, err := g.redis.Eval(
		reserveLoginAttemptScript,
		[]string{accountKey, ipKey},
		seconds,
		g.accountMax,
		g.ipMax,
	).Result()
	if err != nil {
		return false, errs.CommonServiceUnavailable
	}
	values, ok := result.([]interface{})
	if !ok || len(values) != 3 {
		return false, errs.CommonServiceUnavailable
	}
	if loginAttemptCount(values, 0) != 1 {
		return false, errs.CommonFrequentOperationError
	}

	accountCount := loginAttemptCount(values, 1)
	ipCount := loginAttemptCount(values, 2)
	return g.evaluateAttempts(accountCount, ipCount)
}

func (g *loginGuard) evaluateAttempts(accountCount, ipCount int64) (bool, error) {
	if accountCount >= g.accountMax || ipCount >= g.ipMax {
		return false, errs.CommonFrequentOperationError
	}

	return accountCount >= g.captchaAfter, nil
}

func (g *loginGuard) resetAccount(countryCode, mobile string) error {
	accountKey, _ := loginAttemptKeys(countryCode, mobile, "")
	return g.redis.Del(accountKey).Err()
}

func loginAttemptKeys(countryCode, mobile, clientIP string) (string, string) {
	accountDigest := sha256.Sum256([]byte(countryCode + "\x00" + mobile))
	if clientIP == "" {
		clientIP = "unknown"
	}
	ipDigest := sha256.Sum256([]byte(clientIP))
	return "auth:login:account:" + hex.EncodeToString(accountDigest[:]),
		"auth:login:ip:" + hex.EncodeToString(ipDigest[:])
}

func loginAttemptCount(values []interface{}, index int) int64 {
	if index >= len(values) || values[index] == nil {
		return 0
	}
	value, err := strconv.ParseInt(fmt.Sprint(values[index]), 10, 64)
	if err != nil || value < 0 {
		return 0
	}
	return value
}
