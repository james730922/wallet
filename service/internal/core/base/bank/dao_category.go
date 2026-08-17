package bank

import (
	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newCategoryDAO() ICategoryDAO {
	return &categoryDAO{}

}

type ICategoryDAO interface {
	List(dc *gorm.DB) ([]*model.DepositCategory, error)
	GetItems(dc *gorm.DB) ([]*model.BankDepositCategoryItem, error)
	ListType(dc *gorm.DB) ([]*model.BankDepositCategoryType, error)
}

type categoryDAO struct{}

func (dao *categoryDAO) List(dc *gorm.DB) ([]*model.DepositCategory, error) {
	var result []*model.DepositCategory

	if err := dc.Order("sort").
		Find(&result).Error; err != nil {
		return nil, errs.ConvertDB(err)
	}

	return result, nil
}

func (dao *categoryDAO) GetItems(dc *gorm.DB) ([]*model.BankDepositCategoryItem, error) {
	var result []*model.BankDepositCategoryItem

	if err := dc.Find(&result).Error; err != nil {
		return nil, errs.ConvertDB(err)
	}

	return result, nil
}

func (dao *categoryDAO) ListType(dc *gorm.DB) ([]*model.BankDepositCategoryType, error) {
	var result []*model.BankDepositCategoryType

	if err := dc.Find(&result).Error; err != nil {
		return nil, errs.ConvertDB(err)
	}

	return result, nil
}
