package auth

import (
	"context"

	"github.com/james730922/wallet/service/internal/pb/zqbapis"
)

func newLoginMember() ILoginMember {
	return &loginMemberUseCase{
		logout:                   newLoginMemberLogout(),
		regist:                   newLoginMemberRegistration(),
		login:                    newLoginMemberWithPasswd(),
		updatePasswd:             newUpdatePasswd(),
		firstUseScanPay:          newSecurityPasswdFirstUse(),
		forgetScanPayPasswd:      newSecurityPasswdForget(),
		securityPasswdIdentifier: newSecurityPasswdIdentifier(),
		updateScanPayPasswd:      newSecurityPasswdUpdate(),
		captchaGeeTestRegister:   newCaptchaGeeTestRegister(),
	}
}

type ILoginMember interface {
	// 登出
	IdentifierLogout(ctx context.Context) (*zqbapis.LogoutResp, error)
	// 註冊
	Registration(ctx context.Context, req *zqbapis.LoginRegistrationReq) (*zqbapis.LoginRegistrationResp, error)
	// 以密碼登錄
	LoginWithPasswd(ctx context.Context, req *zqbapis.LoginIdentifierPasswdReq) (*zqbapis.LoginIdentifierPasswdResp, error)
	// 更新密碼
	UpdatePasswd(ctx context.Context, req *zqbapis.LoginUpdatePasswdReq) (*zqbapis.LoginUpdatePasswdResp, error)
	// 驗證使否首次安全密碼
	FirstUseSecurityPasswd(ctx context.Context) (*zqbapis.SecurityPasswdFirstUseResp, error)
	// 忘記安全密碼
	ForgetSecurityPasswd(ctx context.Context, req *zqbapis.SecurityPasswdForgetReq) (*zqbapis.SecurityPasswdForgetResp, error)
	// 更新安全密碼
	UpdateSecurityPasswd(ctx context.Context, req *zqbapis.SecurityPasswdUpdateReq) (*zqbapis.SecurityPasswdUpdateResp, error)
	// 安全密碼認證，取得安全碼 Token
	SecurityIdentifier(ctx context.Context, req *zqbapis.SecurityPasswdIdentifierReq) (*zqbapis.SecurityPasswdIdentifierResp, error)
	// 滑動驗證註冊
	CaptchaGeeTestRegister(ctx context.Context) (*zqbapis.CaptchaGeeTestRegisterResp, error)
}

// 會員登入
type loginMemberUseCase struct {
	logout                   *loginMemberLogout
	regist                   *loginMemberRegistration
	login                    *loginMemberWithPasswd
	updatePasswd             *updatePasswd
	firstUseScanPay          *securityPasswdFirstUse
	forgetScanPayPasswd      *securityPasswdForget
	securityPasswdIdentifier *securityPasswdIdentifier
	updateScanPayPasswd      *securityPasswdUpdate
	captchaGeeTestRegister   *captchaGeeTestRegister
}

// 登出
func (lg *loginMemberUseCase) IdentifierLogout(ctx context.Context) (*zqbapis.LogoutResp, error) {
	return lg.logout.Handler(ctx)
}

// 註冊
func (lg *loginMemberUseCase) Registration(ctx context.Context, req *zqbapis.LoginRegistrationReq) (*zqbapis.LoginRegistrationResp, error) {
	return lg.regist.Handler(ctx, req)
}

func (lg *loginMemberUseCase) LoginWithPasswd(ctx context.Context, req *zqbapis.LoginIdentifierPasswdReq) (*zqbapis.LoginIdentifierPasswdResp, error) {
	return lg.login.Handler(ctx, req)
}

func (lg *loginMemberUseCase) UpdatePasswd(ctx context.Context, req *zqbapis.LoginUpdatePasswdReq) (*zqbapis.LoginUpdatePasswdResp, error) {
	return lg.updatePasswd.Handler(ctx, req)
}

func (lg *loginMemberUseCase) FirstUseSecurityPasswd(ctx context.Context) (*zqbapis.SecurityPasswdFirstUseResp, error) {
	return lg.firstUseScanPay.Handler(ctx)
}

func (lg *loginMemberUseCase) ForgetSecurityPasswd(ctx context.Context, req *zqbapis.SecurityPasswdForgetReq) (*zqbapis.SecurityPasswdForgetResp, error) {
	return lg.forgetScanPayPasswd.Handler(ctx, req)
}

func (lg *loginMemberUseCase) UpdateSecurityPasswd(ctx context.Context, req *zqbapis.SecurityPasswdUpdateReq) (*zqbapis.SecurityPasswdUpdateResp, error) {
	return lg.updateScanPayPasswd.Handler(ctx, req)
}

func (lg *loginMemberUseCase) SecurityIdentifier(ctx context.Context, req *zqbapis.SecurityPasswdIdentifierReq) (*zqbapis.SecurityPasswdIdentifierResp, error) {
	return lg.securityPasswdIdentifier.Handler(ctx, req)
}

func (lg *loginMemberUseCase) CaptchaGeeTestRegister(ctx context.Context) (*zqbapis.CaptchaGeeTestRegisterResp, error) {
	return lg.captchaGeeTestRegister.Handler(ctx)
}
