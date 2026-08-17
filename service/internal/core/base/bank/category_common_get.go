package bank

import (
	"context"

	"github.com/ahmetb/go-linq/v3"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/cache"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newCategoryCommonGet() *categoryCommonGet {
	return &categoryCommonGet{}
}

type categoryCommonGet struct {
}

func (c *categoryCommonGet) Handler(ctx context.Context, cond *condition.DepositCategoryQuery) ([]*model.DepositCategory, error) {
	data, ok := c.getCache(cache.KeyBankDepositCategory)
	if !ok {
		newData, err := c.get()
		if err != nil {
			logger.ApLog().Error(err)
			return nil, err
		}

		data = newData
	}

	results := c.find(data, cond)

	return results, nil
}

func (c *categoryCommonGet) get() ([]*model.DepositCategory, error) {
	if result, ok := c.getCache(cache.KeyBankDepositCategory); ok {
		return result, nil
	}

	result, err := dao.Category.List(packet.DB.New())
	if err != nil {
		logger.ApLog().Error(err)
		return nil, errs.CommonNoData
	}
	cache.Cache().Set(cache.KeyBankDepositCategory, result, 0)

	return result, nil
}

func (c *categoryCommonGet) getCache(key string) ([]*model.DepositCategory, bool) {
	data, ok := cache.Cache().Get(key)
	if !ok {
		return nil, ok
	}
	result := data.([]*model.DepositCategory)
	return result, ok
}

func (c *categoryCommonGet) find(items []*model.DepositCategory, cond *condition.DepositCategoryQuery) []*model.DepositCategory {
	var res []*model.DepositCategory
	linq.From(items).
		Where(func(i interface{}) bool {
			if cond.ID == nil {
				return true
			}
			return *cond.ID == i.(*model.DepositCategory).ID
		}).
		Where(func(i interface{}) bool {
			if cond.Name == nil {
				return true
			}
			return *cond.Name == i.(*model.DepositCategory).Name
		}).
		OrderByT(func(dc *model.DepositCategory) int {
			return dc.Sort
		}).ToSlice(&res)

	return res
}
