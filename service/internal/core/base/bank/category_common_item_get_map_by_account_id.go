package bank

import (
	"context"

	"github.com/ahmetb/go-linq/v3"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/cache"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
)

func newCategoryCommonItemGetMapByAccountID() *categoryCommonItemGetMapByAccountID {
	return &categoryCommonItemGetMapByAccountID{}
}

type categoryCommonItemGetMapByAccountID struct {
}

func (c *categoryCommonItemGetMapByAccountID) Handler(ctx context.Context) (map[int64]*model.BankDepositCategoryItem, error) {
	items, err := c.get(ctx)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, err
	}

	result := c.toMap(items)

	return result, err
}

func (c *categoryCommonItemGetMapByAccountID) get(ctx context.Context) ([]*model.BankDepositCategoryItem, error) {
	result, ok := c.checkCategoryCache(cache.KeyBankDepositCategoryItem)
	if !ok {
		items, err := dao.Category.GetItems(packet.DB.New())
		if err != nil {
			return nil, err
		}
		result = items
		cache.Cache().Set(cache.KeyBankDepositCategoryItem, result, 0)
	}

	return result, nil
}

func (c *categoryCommonItemGetMapByAccountID) checkCategoryCache(key string) ([]*model.BankDepositCategoryItem, bool) {
	data, ok := cache.Cache().Get(key)
	if !ok {
		return nil, ok
	}
	result := data.([]*model.BankDepositCategoryItem)
	return result, ok
}

func (c *categoryCommonItemGetMapByAccountID) toMap(items []*model.BankDepositCategoryItem) map[int64]*model.BankDepositCategoryItem {
	result := make(map[int64]*model.BankDepositCategoryItem)
	linq.From(items).
		Select(func(i interface{}) interface{} {
			return linq.KeyValue{
				Key:   i.(*model.BankDepositCategoryItem).AccountID,
				Value: i,
			}
		}).
		ToMap(&result)

	return result
}
