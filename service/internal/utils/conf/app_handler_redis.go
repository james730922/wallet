package conf

type RedisHandler struct{}

func (RedisHandler) GetHost() string {
	return appConf.v.GetString("redis.host")
}

func (RedisHandler) GetMockHost() string {
	return appConf.v.GetString("redis.mock_host")
}

func (RedisHandler) GetPasswd() string {
	return appConf.v.GetString("redis.passwd")
}

func (RedisHandler) GetDB() int {
	return appConf.v.GetInt("redis.db")
}
