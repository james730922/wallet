package bank

import (
	"context"

	"github.com/ahmetb/go-linq/v3"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/cache"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
)

func newBankCommonList() *bankCommonList {
	return &bankCommonList{}
}

type bankCommonList struct {
}

func (b *bankCommonList) Handler(ctx context.Context, cond *condition.BankCodeQuery) ([]*model.BankCode, error) {
	data, ok := b.checkCache(cache.KeyBankCode)
	if !ok {
		bankCode, err := b.fetch(ctx)
		if err != nil {
			return nil, err
		}
		data = bankCode
	}
	result := b.filter(data, cond)

	return result, nil
}

func (b *bankCommonList) fetch(ctx context.Context) ([]*model.BankCode, error) {
	result, ok := b.checkCache(cache.KeyBankCode)
	if !ok {
		bankCode, err := b.getBankCode(ctx)
		if err != nil {
			logger.ApLog().Error(err)
			return nil, err
		}
		cache.Cache().Set(cache.KeyBankCode, bankCode, 0)
		result = bankCode
	}

	return result, nil
}

func (b *bankCommonList) checkCache(key string) ([]*model.BankCode, bool) {
	data, ok := cache.Cache().Get(key)
	if !ok {
		return nil, ok
	}
	result := data.([]*model.BankCode)
	return result, ok
}

func (b *bankCommonList) getBankCode(ctx context.Context) ([]*model.BankCode, error) {
	result, err := dao.Bank.List(packet.DB.New())
	if err != nil {
		logger.ApLog().Errorf("select bank code err: %v", err)
		return nil, err
	}

	return result, nil
}

func (b *bankCommonList) filter(banks []*model.BankCode, cond *condition.BankCodeQuery) []*model.BankCode {
	res := []*model.BankCode{}

	linq.From(banks).
		Where(func(i interface{}) bool {
			if cond.Code == nil {
				return true
			}
			return i.(*model.BankCode).Code == *cond.Code
		}).
		Where(func(i interface{}) bool {
			if cond.Status == nil {
				return true
			}
			return int(i.(*model.BankCode).Status) == *cond.Status
		}).
		Where(func(i interface{}) bool {
			if cond.BankName == nil {
				return true
			}
			return i.(*model.BankCode).Name == *cond.BankName
		}).
		ToSlice(&res)

	return res
}
