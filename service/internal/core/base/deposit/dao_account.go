package deposit

import (
	"github.com/jinzhu/gorm"

	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newAccountDAO() IAccountDAO {
	return &accountDAO{}
}

type IAccountDAO interface {
	First(dc *gorm.DB, cond condition.IQuery) (*model.BankAccount, error)
}

type accountDAO struct{}

func (dao *accountDAO) First(dc *gorm.DB, cond condition.IQuery) (*model.BankAccount, error) {
	var result model.BankAccount
	if err := dc.Where(cond.Where()).First(&result).Error; err != nil {
		return nil, errs.ConvertDB(err)
	}

	return &result, nil
}
