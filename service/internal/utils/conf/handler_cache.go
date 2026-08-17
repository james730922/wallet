package conf

import (
	"time"
)

type CacheHandler struct{}

// 預設過期秒數
func (CacheHandler) GetDefaultExpiration() time.Duration {
	return zqbConf.v.GetDuration("cache.default_expiration") * time.Second
}
