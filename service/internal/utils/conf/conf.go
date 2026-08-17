package conf

import (
	"time"
)

type EnvType = string

const (
	EnvTypeLocal EnvType = "local"
	EnvTypeDev   EnvType = "dev"
	EnvTypeProd  EnvType = "prod"

	defaultEnv EnvType = EnvTypeLocal
)

var (
	zqbConf *ZQBConf
	appConf *AppConf
)

func Mock() {
	zqbConf = newZQBConf()
	zqbConf.loadEnv()

	appConf = newAppConf()
	appConf.Load()

}

func Start() {
	zqbConf = newZQBConf()
	zqbConf.Load()

	appConf = newAppConf()
	appConf.Load()

}

type IConf interface {
	GetLastChangeTime() time.Time
}

func GetZqbConf() IConf {
	return zqbConf
}

func GetAppConf() IConf {
	return appConf
}
