package bank

import (
	"context"

	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/cache"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newCategoryCommonTypeGet() *categoryCommonTypeGet {
	return &categoryCommonTypeGet{}
}

type categoryCommonTypeGet struct {
}

func (c *categoryCommonTypeGet) Handler(ctx context.Context) ([]*model.BankDepositCategoryType, error) {
	result, ok := c.getCache(cache.KeyBankDepositCategoryType)
	if !ok {
		newData, err := c.get()
		if err != nil {
			logger.ApLog().Error(err)
			return nil, err
		}

		result = newData
	}

	return result, nil
}

func (c *categoryCommonTypeGet) get() ([]*model.BankDepositCategoryType, error) {
	if result, ok := c.getCache(cache.KeyBankDepositCategoryType); ok {
		return result, nil
	}

	result, err := dao.Category.ListType(packet.DB.New())
	if err != nil {
		logger.ApLog().Error(err)
		return nil, errs.CommonNoData
	}
	cache.Cache().Set(cache.KeyBankDepositCategoryType, result, 0)

	return result, nil
}

func (c *categoryCommonTypeGet) getCache(key string) ([]*model.BankDepositCategoryType, bool) {
	data, ok := cache.Cache().Get(key)
	if !ok {
		return nil, ok
	}
	result := data.([]*model.BankDepositCategoryType)
	return result, ok
}
