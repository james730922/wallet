package geetestsdk

type FuncGeeTestCaptchaOpt func(o *geeTestCaptcha)

func (f FuncGeeTestCaptchaOpt) Set(o *geeTestCaptcha) {
	f(o)
}

type GeeTestCaptchaOpt func()

func (o GeeTestCaptchaOpt) Alert(alert func()) FuncGeeTestCaptchaOpt {
	return func(o *geeTestCaptcha) {
		o.alert = alert
	}
}
