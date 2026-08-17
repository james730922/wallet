package conf

import (
	"flag"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Env() string {
	return zqbConf.properties.env
}

func Service() ServiceHandler {
	return zqbConf.properties.zqbYaml.Service
}

func Sign() SignHandler {
	return zqbConf.properties.zqbYaml.Sign
}

func Log() LogHandler {
	return zqbConf.properties.zqbYaml.Log
}

func Cache() CacheHandler {
	return zqbConf.properties.zqbYaml.Cache
}

func LoginMember() LoginMemberHandler {
	return zqbConf.properties.zqbYaml.Auth.LoginMember
}

func Pyroscope() PyroscopeHandler {
	return zqbConf.properties.zqbYaml.Pyroscope
}

func Blacklist() BlacklistHandler {
	return zqbConf.properties.zqbYaml.Blacklist
}

func Captcha() CaptchaHandler {
	return zqbConf.properties.zqbYaml.Captcha
}

func DataBase() DataBaseHandler {
	return zqbConf.properties.zqbYaml.DataBase
}

func Observability() ObservabilityHandler {
	return zqbConf.properties.zqbYaml.Observability
}

func GetZQBConf() *ZQBConf {
	return zqbConf
}

func newZQBConf() *ZQBConf {
	return &ZQBConf{
		v:              viper.New(),
		lastChangeTime: time.Now(),
		properties:     &ZQBProperties{},
	}
}

type ZQBConf struct {
	v              *viper.Viper
	lastChangeTime time.Time
	properties     *ZQBProperties
}

func (c *ZQBConf) GetLastChangeTime() time.Time {
	return c.lastChangeTime
}

func (c *ZQBConf) Load() {
	c.loadEnv()
	c.loadYaml()
}

func (c *ZQBConf) loadEnv() {
	envArg := flag.String("env", defaultEnv, "the server run which environment")
	flag.Parse()

	env := ""
	switch e := *envArg; e {
	case EnvTypeLocal, EnvTypeDev, EnvTypeProd:
		env = e
	default:
		env = defaultEnv
	}

	c.properties.env = env
	log.Printf("Env: %s\n", env)
}

func (c *ZQBConf) loadYaml() {
	c.v.SetConfigName("zqb.yaml")
	c.v.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
	c.v.AddConfigPath(configDir("zqb.yaml"))
	if err := c.v.ReadInConfig(); err != nil {
		panic(err)
	}
	c.v.OnConfigChange(func(in fsnotify.Event) {
		c.lastChangeTime = time.Now()
	})
	c.v.WatchConfig()
}

type ZQBProperties struct {
	zqbYaml ZqbYaml
	env     string
}

type ZqbYaml struct {
	Service       ServiceHandler       `yaml:"service"`
	Sign          SignHandler          `yaml:"sign"`
	Log           LogHandler           `yaml:"log"`
	Cache         CacheHandler         `yaml:"cache"`
	Auth          AuthHandler          `yaml:"auth"`
	Pyroscope     PyroscopeHandler     `yaml:"pyroscope"`
	Blacklist     BlacklistHandler     `yaml:"blacklist"`
	Captcha       CaptchaHandler       `yaml:"captcha"`
	DataBase      DataBaseHandler      `yaml:"database"`
	Observability ObservabilityHandler `yaml:"observability"`
}
