package deposit

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"github.com/ahmetb/go-linq/v3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/cache"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newDepositCommonCategoryItem() *depositCommonCategoryItem {
	return &depositCommonCategoryItem{}
}

type depositCommonCategoryItem struct {
	mx sync.Mutex
}

func (de *depositCommonCategoryItem) Handler(ctx context.Context, cond *condition.BankDepositCategoryItemViewCond) ([]*model.BankDepositCategoryItemView, error) {
	data, ok := de.checkCache(cache.KeyBankDepositCategoryItemForMember)
	if !ok {
		category, err := de.fetch(ctx)
		if err != nil {
			logger.ApLog().Error(err)
			return nil, err
		}
		data = category
	}
	result := de.filter(data, cond)

	return result, nil
}

func (de *depositCommonCategoryItem) fetch(ctx context.Context) ([]*model.BankDepositCategoryItemView, error) {
	categorieItems, err := dao.Member.SelectCategoryItem(packet.DB.New())
	if err != nil {
		logger.ApLog().Error(err)
		return nil, errs.CommonNoData
	}
	cache.Cache().Set(cache.KeyBankDepositCategoryItemForMember, categorieItems, 0)

	return categorieItems, nil
}

func (de *depositCommonCategoryItem) checkCache(key string) ([]*model.BankDepositCategoryItemView, bool) {
	data, ok := cache.Cache().Get(key)
	if !ok {
		return nil, ok
	}
	result := data.([]*model.BankDepositCategoryItemView)
	return result, ok
}

func (de *depositCommonCategoryItem) filter(i interface{}, cond *condition.BankDepositCategoryItemViewCond) []*model.BankDepositCategoryItemView {
	var res []*model.BankDepositCategoryItemView

	linq.From(i.([]*model.BankDepositCategoryItemView)).
		Where(func(i interface{}) bool {
			if cond.AccountStatus == nil {
				return true
			}
			return i.(*model.BankDepositCategoryItemView).AccountStatus == aws.IntValue(cond.AccountStatus)
		}).
		Where(func(i interface{}) bool {
			if cond.AccountVisible == nil {
				return true
			}
			return i.(*model.BankDepositCategoryItemView).AccountVisible == aws.IntValue(cond.AccountVisible)
		}).
		Where(func(i interface{}) bool {
			if cond.CategoryStatus == nil {
				return true
			}
			return i.(*model.BankDepositCategoryItemView).CategoryStatus == aws.IntValue(cond.CategoryStatus)
		}).
		Where(func(i interface{}) bool {
			if cond.MemberLevel == nil || *cond.MemberLevel == model.BankAccountAllLevels.Int() {
				return true
			}
			levels := strings.Split(strings.Trim(i.(*model.BankDepositCategoryItemView).Levels, ","), ",")
			for _, v := range levels {
				if v == strconv.Itoa(*cond.MemberLevel) || v == strconv.Itoa(model.BankAccountAllLevels.Int()) {
					return true
				}
			}
			return false
		}).ToSlice(&res)

	return res
}
