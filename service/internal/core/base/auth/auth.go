package auth

import (
	"sync"

	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"go.uber.org/dig"

	"github.com/james730922/wallet/service/internal/core/base/captcha"
	"github.com/james730922/wallet/service/internal/core/base/user"
	"github.com/james730922/wallet/service/internal/utils/conf"
	"github.com/james730922/wallet/service/internal/utils/password"
)

var (
	packet authSet
)

var (
	once sync.Once
	self *auth
)

func NewAuth(set authSet) authOut {
	once.Do(func() {
		packet = set
		password.SetMaxConcurrentArgon2(conf.LoginMember().GetLoginMaxConcurrentHashes())
		guard := newLoginGuard(set.Redis)
		registrationGuard := newRegistrationGuard(set.Redis)

		self = &auth{
			authOut: authOut{
				Token:       newToken(),
				LoginMember: newLoginMember(),
			},
			memberRepository:  newMemberRepository(),
			loginGuard:        guard,
			registrationGuard: registrationGuard,
		}

	})

	return self.authOut
}

type authSet struct {
	dig.In

	DB             *gorm.DB
	Member         user.IMemberCommon
	Redis          *redis.Client
	GeeTestCaptcha captcha.IGeeTestCaptcha
}

type auth struct {
	authOut
	memberRepository  iMemberRepository
	loginGuard        *loginGuard
	registrationGuard *registrationGuard
}

type authOut struct {
	dig.Out

	Token       IToken
	LoginMember ILoginMember
}
