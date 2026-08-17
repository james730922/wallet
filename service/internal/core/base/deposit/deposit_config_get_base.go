package deposit

import (
	"context"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/cache"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

type depositConfigGetBase struct {
}

func (ad *depositConfigGetBase) Handler(ctx context.Context) (*model.DepositConfig, error) {
	result, ok := ad.checkCache(cache.KeyDepositConfig)
	if !ok {
		config, err := ad.fetch(ctx)
		if err != nil {
			return nil, err
		}
		result = config
	}

	return result, nil
}

func (ad *depositConfigGetBase) fetch(ctx context.Context) (*model.DepositConfig, error) {
	result, ok := ad.checkCache(cache.KeyDepositConfig)
	if !ok {
		config, err := ad.getDepositConfig(ctx)
		if err != nil {
			logger.ApLog().Error(err)
			return nil, err
		}
		cache.Cache().Set(cache.KeyDepositConfig, config, 0)
		result = config
	}

	return result, nil
}

func (ad *depositConfigGetBase) checkCache(key string) (*model.DepositConfig, bool) {
	data, ok := cache.Cache().Get(key)
	if !ok {
		return nil, ok
	}
	result := data.(*model.DepositConfig)
	return result, ok
}

func (ad *depositConfigGetBase) getDepositConfig(ctx context.Context) (*model.DepositConfig, error) {
	result, err := dao.Config.First(packet.DB.New())
	if err != nil {
		logger.ApLog().Error(err)
		return nil, errs.CommonNoData
	}

	return result, nil
}
