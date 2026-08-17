package binder

import "go.uber.org/dig"

var (
	binder *dig.Container
)

func New() *dig.Container {
	binder = dig.New()

	provideThirdParty(binder)
	provideCore(binder)
	provideUsecase(binder)
	provideController(binder)
	provideApp(binder)

	return binder
}
