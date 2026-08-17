package captcha

import (
	"time"

	"github.com/james730922/wallet/service/internal/utils/errs"

	"gitlab.paradise-soft.com.tw/glob/geetestsdk"

	"github.com/james730922/wallet/service/internal/utils/conf"
)

func newGeeTestCaptcha() *geeTestCaptcha {
	obj := &geeTestCaptcha{}
	obj.Init()

	return obj
}

type IGeeTestCaptcha interface {
	Register(account string) *geetestsdk.RegisterResp
	Validate(challenge, validate, secCode string) error
	Enabled() bool
}

func (c *geeTestCaptcha) Enabled() bool {
	return conf.Captcha().GeeTest().GetEnable()
}

type geeTestCaptcha struct {
	geeTest geetestsdk.ICaptcha

	lastChangeTime time.Time
}

func (c *geeTestCaptcha) Init() {
	c.lastChangeTime = conf.GetZqbConf().GetLastChangeTime()
	c.geeTest = geetestsdk.NewGeeTestCaptcha(c.GenConfig())
	c.OnConfigChanged()
}

func (c *geeTestCaptcha) OnConfigChanged() {
	go func() {
		for {
			confLastChangeTime := conf.GetZqbConf().GetLastChangeTime()
			if c.lastChangeTime.Before(confLastChangeTime) {
				c.lastChangeTime = confLastChangeTime
				c.geeTest.RefreshConfig(c.GenConfig())
			}

			<-time.After(10 * time.Second)
		}
	}()
}

func (c *geeTestCaptcha) Register(account string) *geetestsdk.RegisterResp {
	if !conf.Captcha().GeeTest().GetEnable() {
		return &geetestsdk.RegisterResp{
			Success:    true,
			NewCaptcha: true,
		}
	}

	return c.geeTest.Register(account)
}
func (c *geeTestCaptcha) Validate(challenge, validate, secCode string) error {
	if !conf.Captcha().GeeTest().GetEnable() {
		return nil
	}

	result := c.geeTest.Validate(challenge, validate, secCode)
	if result.Status != 1 {
		return errs.LoginMemberCaptchaValidateFailed
	}

	return nil
}

func (c *geeTestCaptcha) GenConfig() geetestsdk.CaptchaConfig {
	return geetestsdk.CaptchaConfig{
		GeeTestID:                 conf.Captcha().GeeTest().GetID(),
		GeeTestKey:                conf.Captcha().GeeTest().GetKey(),
		GeeTestUserIDSalt:         conf.Captcha().GeeTest().GetSalt(),
		GeeTestByPassCycleTimeSec: conf.Captcha().GeeTest().GetByPassCycleTimeSec(),
	}
}
