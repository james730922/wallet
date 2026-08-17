package deposit

import (
	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newDepositConfigMemberDAO() IDepositConfigMemberDAO {
	return &depositConfigMemberDAO{}
}

type IDepositConfigMemberDAO interface {
	First(dc *gorm.DB, memberID int64) (*model.DepositConfigMember, error)
}

type depositConfigMemberDAO struct {
}

func (dao *depositConfigMemberDAO) First(dc *gorm.DB, memberID int64) (*model.DepositConfigMember, error) {
	result := model.DepositConfigMember{}

	err := dc.
		Where("member_id = ?", memberID).
		First(&result).Error
	if err != nil {
		return nil, errs.ConvertDB(err)
	}

	return &result, nil
}
