package wallet

import (
	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/db"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newWalletDAO() IWalletDAO {
	return &walletDAO{}
}

type IWalletDAO interface {
	First(dc *gorm.DB, memberID int64) (*model.WalletMember, error)
	FirstForUpdate(dc *gorm.DB, memberID int64) (*model.WalletMember, error)
	Insert(dc *gorm.DB, cond *model.WalletMember) error
	Update(dc *gorm.DB, cond condition.IUpdate) error
}

type walletDAO struct{}

func (dao *walletDAO) First(dc *gorm.DB, memberID int64) (*model.WalletMember, error) {
	var result model.WalletMember

	if err := dc.Where("member_id = ?", memberID).
		First(&result).Error; err != nil {
		return nil, errs.ConvertDB(err)
	}

	return &result, nil
}

func (dao *walletDAO) FirstForUpdate(dc *gorm.DB, memberID int64) (*model.WalletMember, error) {
	var result model.WalletMember

	if err := dc.
		Set(db.GormSetSelectForUpdate()).
		Where("member_id = ?", memberID).
		First(&result).Error; err != nil {
		return nil, errs.ConvertDB(err)
	}

	return &result, nil
}

func (dao *walletDAO) Insert(dc *gorm.DB, cond *model.WalletMember) error {
	if err := dc.Create(cond).Error; err != nil {
		return err
	}

	return nil
}

func (dao *walletDAO) Update(dc *gorm.DB, cond condition.IUpdate) error {
	if err := dc.Model(model.WalletMember{}).
		Where(cond.Where()).
		Updates(cond.Update()).Error; err != nil {
		return errs.ConvertDB(err)
	}

	return nil
}
