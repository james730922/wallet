package auth

import (
	"context"

	"github.com/james730922/wallet/service/internal/pb/zqbapis"
)

func newCaptchaGeeTestRegister() *captchaGeeTestRegister {
	return &captchaGeeTestRegister{}
}

type captchaGeeTestRegister struct{}

func (hd *captchaGeeTestRegister) Handler(ctx context.Context) (*zqbapis.CaptchaGeeTestRegisterResp, error) {

	registerResp := packet.GeeTestCaptcha.Register("")

	return &zqbapis.CaptchaGeeTestRegisterResp{
		Success:    registerResp.Success,
		Gt:         registerResp.GT,
		Challenge:  registerResp.Challenge,
		NewCaptcha: registerResp.NewCaptcha,
	}, nil
}
