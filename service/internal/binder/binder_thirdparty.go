package binder

import (
	"go.uber.org/dig"

	"github.com/james730922/wallet/service/internal/thirdparty/cache"
	"github.com/james730922/wallet/service/internal/thirdparty/db"
	"github.com/james730922/wallet/service/internal/thirdparty/event"
	"github.com/james730922/wallet/service/internal/thirdparty/fileserver"
	"github.com/james730922/wallet/service/internal/thirdparty/observability"
	"github.com/james730922/wallet/service/internal/thirdparty/redis"
	"github.com/james730922/wallet/service/internal/thirdparty/snowflake"
)

func provideThirdParty(binder *dig.Container) {
	if err := binder.Provide(snowflake.NewIDGenerator); err != nil {
		panic(err)
	}

	if err := binder.Provide(db.NewGORM); err != nil {
		panic(err)
	}

	if err := binder.Provide(db.NewGORMSlave, dig.Name("dbSlave")); err != nil {
		panic(err)
	}

	if err := binder.Provide(fileserver.New); err != nil {
		panic(err)
	}

	if err := binder.Provide(fileserver.NewZQBFileServer); err != nil {
		panic(err)
	}

	if err := binder.Provide(redis.NewRedis); err != nil {
		panic(err)
	}

	if err := binder.Invoke(cache.NewCache); err != nil {
		panic(err)
	}

	if err := binder.Invoke(event.New); err != nil {
		panic(err)
	}

	if err := binder.Provide(observability.NewMetrics); err != nil {
		panic(err)
	}
	if err := binder.Provide(observability.NewTracing); err != nil {
		panic(err)
	}

}
