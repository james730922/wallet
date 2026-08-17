package deposit

import (
	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/db"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newDepositConfigMemberLevelDAO() IDepositConfigMemberLevelDAO {
	return &depositConfigMemberLevelDAO{}
}

type IDepositConfigMemberLevelDAO interface {
	First(dc *gorm.DB, memberLevel int64) (*model.DepositConfigMemberLevel, error)
	List(dc *gorm.DB, cond *condition.DepositConfigMemberLevelQuery) ([]*model.DepositConfigMemberLevel, error)
}

type depositConfigMemberLevelDAO struct {
}

func (dao *depositConfigMemberLevelDAO) First(dc *gorm.DB, memberLevel int64) (*model.DepositConfigMemberLevel, error) {
	result := model.DepositConfigMemberLevel{}

	err := dc.
		Where("member_level = ?", memberLevel).
		First(&result).Error
	if err != nil {
		return nil, errs.ConvertDB(err)
	}

	return &result, nil
}

func (dao *depositConfigMemberLevelDAO) List(dc *gorm.DB, cond *condition.DepositConfigMemberLevelQuery) ([]*model.DepositConfigMemberLevel, error) {
	iQuery := condition.NewQuery(cond)
	result := make([]*model.DepositConfigMemberLevel, 0, iQuery.Paging().Size)

	if err := dc.
		Scopes(db.ParseWhere(iQuery.Where())).
		Scopes(db.ParsePaging(iQuery.Paging())).Find(&result).
		Error; err != nil {
		return nil, errs.ConvertDB(err)
	}
	return result, nil
}
