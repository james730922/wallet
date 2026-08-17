package binder

import (
	"go.uber.org/dig"

	"github.com/james730922/wallet/service/internal/core/base/auth"
	"github.com/james730922/wallet/service/internal/core/base/bank"
	"github.com/james730922/wallet/service/internal/core/base/captcha"
	"github.com/james730922/wallet/service/internal/core/base/deposit"
	"github.com/james730922/wallet/service/internal/core/base/paycode"
	"github.com/james730922/wallet/service/internal/core/base/scanpay"
	"github.com/james730922/wallet/service/internal/core/base/simplifymemberid"
	"github.com/james730922/wallet/service/internal/core/base/transaction"
	"github.com/james730922/wallet/service/internal/core/base/user"
	"github.com/james730922/wallet/service/internal/core/base/wallet"
)

func provideCore(binder *dig.Container) {
	providers := []interface{}{
		bank.NewBank,
		auth.NewAuth,
		wallet.NewWallet,
		user.NewUser,
		transaction.NewTransaction,
		deposit.NewDeposit,
		simplifymemberid.NewSimplifyMemberID,
		paycode.NewPayment,
		scanpay.NewScanPay,
		captcha.NewCaptcha,
	}
	for _, provider := range providers {
		if err := binder.Provide(provider); err != nil {
			panic(err)
		}
	}
}
