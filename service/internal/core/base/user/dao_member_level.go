package user

import (
	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newMemberLevelDAO() iMemberLevelDAO {
	return &memberLevelDAO{}
}

type iMemberLevelDAO interface {
	// 查找單一分級
	First(dc *gorm.DB, cond condition.IQuery) (*model.MemberLevel, error)
	// 列出所有分級
	List(dc *gorm.DB) ([]*model.MemberLevel, error)
	// 更新註冊時的會員層級人數
	Update(dc *gorm.DB, cond condition.IUpdate) error
}

type memberLevelDAO struct{}

// 查找單一分級
func (dao *memberLevelDAO) First(dc *gorm.DB, cond condition.IQuery) (*model.MemberLevel, error) {
	var result model.MemberLevel

	if err := dc.Where(cond.Where()).First(&result).Error; err != nil {
		return nil, errs.ConvertDB(err)
	}

	return &result, nil
}

// 列出所有分級
func (dao *memberLevelDAO) List(dc *gorm.DB) ([]*model.MemberLevel, error) {
	var result []*model.MemberLevel

	if err := dc.Order("sort").
		Find(&result).Error; err != nil {
		return nil, errs.ConvertDB(err)
	}

	return result, nil
}

// 更新註冊時的會員層級人數
func (dao *memberLevelDAO) Update(dc *gorm.DB, cond condition.IUpdate) error {
	if err := dc.Model(model.MemberLevel{}).
		Where(cond.Where()).
		Updates(cond.Update()).
		Error; err != nil {
		return err
	}

	return nil
}
