package conf

import (
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func DB() DBHandler {
	return appConf.properties.appYaml.DB
}

func FileServer() FileServerHandler {
	return appConf.properties.appYaml.FileServer
}

func Redis() RedisHandler {
	return appConf.properties.appYaml.Redis
}

func ScanPayCrypto() ScanPayCryptoHandler {
	return appConf.properties.appYaml.ScanPayCrypto
}

func newAppConf() *AppConf {
	return &AppConf{
		v:              viper.New(),
		lastChangeTime: time.Now(),
		properties:     &AppProperties{},
	}
}

type AppConf struct {
	v              *viper.Viper
	lastChangeTime time.Time
	properties     *AppProperties
}

func (c *AppConf) GetLastChangeTime() time.Time {
	return c.lastChangeTime
}

func (c *AppConf) Load() {
	c.loadYaml()
}

func (c *AppConf) loadYaml() {
	c.v.SetConfigName("app.conf")
	c.v.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
	c.v.AddConfigPath(configDir("app.conf"))
	if err := c.v.ReadInConfig(); err != nil {
		panic(err)
	}
	c.v.OnConfigChange(func(in fsnotify.Event) {
		c.lastChangeTime = time.Now()
	})
	c.v.WatchConfig()
}

type AppProperties struct {
	appYaml AppYaml
}

type AppYaml struct {
	DB            DBHandler            `yaml:"db"`
	FileServer    FileServerHandler    `yaml:"fileserver"`
	Redis         RedisHandler         `yaml:"redis"`
	ScanPayCrypto ScanPayCryptoHandler `yaml:"scanpay_crypto"`
}
