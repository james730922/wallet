package binder

import (
	"go.uber.org/dig"

	"github.com/james730922/wallet/service/internal/core/usecase/member/scanpaymember"
	"github.com/james730922/wallet/service/internal/core/usecase/member/transactionmember"
)

func provideUsecase(binder *dig.Container) {
	providers := []interface{}{
		scanpaymember.NewScanPay,
		transactionmember.NewTransactionMember,
	}
	for _, provider := range providers {
		if err := binder.Provide(provider); err != nil {
			panic(err)
		}
	}
}
