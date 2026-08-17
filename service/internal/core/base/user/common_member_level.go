package user

import (
	"context"

	"github.com/jinzhu/gorm"

	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/cache"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newMemberLevelCommon() IMemberLevelCommon {
	return &memberLevelCommonUseCase{
		get: newMemberLevelCommonGet(),
	}
}

type IMemberLevelCommon interface {
	// 查找分級
	First(dc *gorm.DB, cond *condition.MemberLevelQuery) (*model.MemberLevel, error)
	// 列出所有分級
	List(ctx context.Context, cond *condition.MemberLevelQuery) ([]*model.MemberLevel, error)
	// 更新註冊時的會員層級人數
	Update(dc *gorm.DB, cond *condition.MemberLevelUpdate) error
}
type memberLevelCommonUseCase struct {
	get *memberLevelCommonGet
}

// 查找分級
func (cm *memberLevelCommonUseCase) First(dc *gorm.DB, cond *condition.MemberLevelQuery) (*model.MemberLevel, error) {
	result, err := dao.MemberLevel.First(dc, condition.NewQuery(cond))
	if err != nil {
		return nil, errs.ConvertDB(err)
	}

	return result, nil
}

// 列出所有分級
func (cm *memberLevelCommonUseCase) List(ctx context.Context, cond *condition.MemberLevelQuery) ([]*model.MemberLevel, error) {
	return cm.get.Handler(ctx, cond)
}

// 更新註冊時的會員層級人數
func (cm *memberLevelCommonUseCase) Update(dc *gorm.DB, cond *condition.MemberLevelUpdate) error {
	updCond := condition.NewUpdate(
		cond,
		&condition.MemberLevelUpdate{
			ID: cond.ID,
		},
	)

	if err := dao.MemberLevel.Update(dc, updCond); err != nil {
		logger.ApLog().Error(err)
		return errs.DBInsertFailed
	}

	cm.cleanCache()

	return nil
}

func (cm *memberLevelCommonUseCase) cleanCache() {
	cache.Cache().Delete(cache.KeyMemberLevel)
}
