package binder

import (
	"go.uber.org/dig"

	"github.com/james730922/wallet/service/internal/controller/apictrl"
)

func provideController(binder *dig.Container) {
	for _, provider := range []interface{}{apictrl.NewController} {
		if err := binder.Provide(provider); err != nil {
			panic(err)
		}
	}
}
