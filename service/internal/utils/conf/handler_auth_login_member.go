package conf

import (
	"time"
)

type LoginMemberHandler struct{}

// token expiration min
func (LoginMemberHandler) GetMemberTokenExpiration() time.Duration {
	return zqbConf.v.GetDuration("auth.member_token_expiration") * time.Minute
}

func (LoginMemberHandler) GetMemberSecureTokenExpiration() time.Duration {
	return zqbConf.v.GetDuration("auth.member_security_passwd_token_expiration") * time.Minute
}

func (LoginMemberHandler) GetMemberForgetSecureTokenExpiration() time.Duration {
	return zqbConf.v.GetDuration("auth.member_forget_security_passwd_token_expiration") * time.Minute
}

func (LoginMemberHandler) GetLoginAttemptWindow() time.Duration {
	seconds := zqbConf.v.GetInt("auth.member_login_attempt_window_sec")
	if seconds <= 0 {
		seconds = 300
	}
	return time.Duration(seconds) * time.Second
}

func (LoginMemberHandler) GetLoginAccountMaxAttempts() int64 {
	value := zqbConf.v.GetInt64("auth.member_login_account_max_attempts")
	if value <= 0 {
		return 5
	}
	return value
}

func (LoginMemberHandler) GetLoginIPMaxAttempts() int64 {
	value := zqbConf.v.GetInt64("auth.member_login_ip_max_attempts")
	if value <= 0 {
		return 30
	}
	return value
}

func (LoginMemberHandler) GetLoginCaptchaAfterAttempts() int64 {
	value := zqbConf.v.GetInt64("auth.member_login_captcha_after_attempts")
	if value <= 0 {
		return 3
	}
	return value
}

func (LoginMemberHandler) GetLoginMaxConcurrentHashes() int {
	value := zqbConf.v.GetInt("auth.member_login_max_concurrent_hashes")
	if value <= 0 {
		return 4
	}
	return value
}

func (LoginMemberHandler) GetRegistrationAttemptWindow() time.Duration {
	seconds := zqbConf.v.GetInt("auth.member_registration_attempt_window_sec")
	if seconds <= 0 {
		seconds = 3600
	}
	return time.Duration(seconds) * time.Second
}

func (LoginMemberHandler) GetRegistrationIPMaxAttempts() int64 {
	value := zqbConf.v.GetInt64("auth.member_registration_ip_max_attempts")
	if value <= 0 {
		return 20
	}
	return value
}

func (LoginMemberHandler) GetRegistrationDeviceMaxAttempts() int64 {
	value := zqbConf.v.GetInt64("auth.member_registration_device_max_attempts")
	if value <= 0 {
		return 5
	}
	return value
}
