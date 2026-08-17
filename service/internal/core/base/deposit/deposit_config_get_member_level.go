package deposit

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/cache"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/errs"
	"github.com/james730922/wallet/service/internal/utils/tools"
)

type depositConfigGetMemberLevel struct {
}

func (hd *depositConfigGetMemberLevel) Handler(cond *condition.DepositConfigMemberLevelQuery) ([]*model.DepositConfigMemberLevel, error) {
	cache, ok := hd.getCache()
	if !ok {
		tmp, err := hd.getFromDB(packet.DB.New())
		if err != nil {
			logger.ApLog().Error(err)
			return nil, err
		}
		cache = tmp
	}

	return hd.filter(cache, cond), nil
}

func (hd *depositConfigGetMemberLevel) getCache() ([]*model.DepositConfigMemberLevel, bool) {
	tmp, ok := cache.Cache().Get(cache.KeyDepositConfigMemberLevel)
	if !ok {
		return nil, false
	}

	cache, ok := tmp.([]*model.DepositConfigMemberLevel)
	if !ok {
		return nil, false
	}

	return cache, true
}
func (hd *depositConfigGetMemberLevel) getFromDB(dc *gorm.DB) ([]*model.DepositConfigMemberLevel, error) {
	config, err := dao.ConfigMemberLevel.List(dc, &condition.DepositConfigMemberLevelQuery{})
	if err != nil {
		logger.ApLog().Error(err)
		return nil, errs.CommonNoData
	}

	cache.Cache().Set(cache.KeyDepositConfigMemberLevel, config, 0)
	return config, nil
}

func (hd *depositConfigGetMemberLevel) filter(cache []*model.DepositConfigMemberLevel, cond *condition.DepositConfigMemberLevelQuery) []*model.DepositConfigMemberLevel {
	filteredCache := make([]*model.DepositConfigMemberLevel, 0, len(cache))

	linq.From(cache).
		Where(func(i interface{}) bool {
			if cond.MemberLevels != nil && len(*cond.MemberLevels) > 0 {
				return tools.Int64InSlice(*cond.MemberLevels, i.(*model.DepositConfigMemberLevel).MemberLevel)
			}
			return true
		}).
		Where(func(i interface{}) bool {
			if cond.Statuses != nil && len(*cond.Statuses) > 0 {
				return tools.IntInSlice(*cond.Statuses, int(i.(*model.DepositConfigMemberLevel).Status))
			}
			return true
		}).
		Where(func(i interface{}) bool {
			if cond.Status != nil {
				return *cond.Status == int(i.(*model.DepositConfigMemberLevel).Status)
			}
			return true
		}).
		Where(func(i interface{}) bool {
			if cond.StartAtAddedTime != nil {
				return i.(*model.DepositConfigMemberLevel).AddedTime.After(aws.TimeValue(cond.StartAtAddedTime))
			}
			return true
		}).
		Where(func(i interface{}) bool {
			if cond.EndAtAddedTime != nil {
				return i.(*model.DepositConfigMemberLevel).AddedTime.Before(aws.TimeValue(cond.EndAtAddedTime))
			}
			return true
		}).
		ToSlice(&filteredCache)

	return filteredCache
}
