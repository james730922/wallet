package binder

import (
	"go.uber.org/dig"

	"github.com/james730922/wallet/service/internal/app"
)

func provideApp(binder *dig.Container) {
	for _, provider := range []interface{}{app.NewApiServe, app.NewPyroscopeServe} {
		if err := binder.Provide(provider); err != nil {
			panic(err)
		}
	}
}
