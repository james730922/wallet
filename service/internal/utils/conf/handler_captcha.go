package conf

type ICaptcha interface {
	GeeTest() CaptchaGeeTestHandler
}

type CaptchaHandler struct {
	geeTestHandler CaptchaGeeTestHandler
}

func (c CaptchaHandler) GeeTest() CaptchaGeeTestHandler {
	return c.geeTestHandler
}

type CaptchaGeeTestHandler struct{}

func (CaptchaGeeTestHandler) GetEnable() bool {
	return zqbConf.v.GetBool("captcha.geetest.enable")
}

func (CaptchaGeeTestHandler) GetID() string {
	return zqbConf.v.GetString("captcha.geetest.id")
}

func (CaptchaGeeTestHandler) GetKey() string {
	return zqbConf.v.GetString("captcha.geetest.key")
}

func (CaptchaGeeTestHandler) GetSalt() string {
	return zqbConf.v.GetString("captcha.geetest.salt")
}

func (CaptchaGeeTestHandler) GetByPassCycleTimeSec() int {
	return zqbConf.v.GetInt("captcha.geetest.by_pass_cycle_time_sec")
}
