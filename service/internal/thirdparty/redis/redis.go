package redis

import (
	"github.com/go-redis/redis/v7"

	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/conf"
)

var defaultPoolSize = 1000

func NewRedis() *redis.Client {
	logger.ApLog().Debug("init redis")
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis().GetHost(),
		Password: conf.Redis().GetPasswd(), // no password set
		DB:       conf.Redis().GetDB(),     // use default DB
		PoolSize: defaultPoolSize,
	})

	pong, err := rdb.Ping().Result()
	if err != nil {
		logger.ApLog().Panic(err)
	}

	logger.ApLog().Debugf("redis reply %v", pong)

	return rdb
}

func NewRedisMock() *redis.Client {
	logger.ApLog().Debugf("init mock redis : %s", conf.Redis().GetMockHost())
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis().GetMockHost(),
		Password: conf.Redis().GetPasswd(), // no password set
		DB:       conf.Redis().GetDB(),     // use default DB
		PoolSize: defaultPoolSize,
	})

	pong, err := rdb.Ping().Result()
	if err != nil {
		logger.ApLog().Panic(err)
	}

	logger.ApLog().Debugf("redis reply %v", pong)
	return rdb
}
