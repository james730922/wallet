package bank

import (
	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newBankDAO() IBankDAO {
	return &bankDAO{}
}

type IBankDAO interface {
	List(dc *gorm.DB) ([]*model.BankCode, error)
}

type bankDAO struct {
	model.BankCode
}

func (dao *bankDAO) List(dc *gorm.DB) ([]*model.BankCode, error) {
	var result []*model.BankCode

	if err := dc.Order("code").
		Find(&result).Error; err != nil {
		return nil, errs.ConvertDB(err)
	}

	return result, nil
}
