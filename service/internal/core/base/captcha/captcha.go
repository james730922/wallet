package captcha

import "go.uber.org/dig"

var (
	self *captcha
)

func NewCaptcha() captchaOut {
	self = &captcha{
		geeTest: newGeeTestCaptcha(),
	}

	return captchaOut{
		GeeTest: self.geeTest,
	}
}

type captcha struct {
	geeTest *geeTestCaptcha
}

type captchaOut struct {
	dig.Out

	GeeTest IGeeTestCaptcha
}
