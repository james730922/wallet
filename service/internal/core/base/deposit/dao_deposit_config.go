package deposit

import (
	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newDepositConfigDAO() IDepositConfigDAO {
	return &depositConfigDAO{}
}

type IDepositConfigDAO interface {
	First(dc *gorm.DB) (*model.DepositConfig, error)
}

type depositConfigDAO struct{}

func (dao *depositConfigDAO) First(dc *gorm.DB) (*model.DepositConfig, error) {
	result := model.DepositConfig{}

	if err := dc.First(&result).Error; err != nil {
		return nil, errs.ConvertDB(err)
	}

	return &result, nil
}
