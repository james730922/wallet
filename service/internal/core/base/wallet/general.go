package wallet

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/errs"
	"github.com/james730922/wallet/service/internal/utils/signs"
)

func newGeneral() IGeneral {
	return &generalUseCase{}
}

type IGeneral interface {
	Get(dc *gorm.DB, memberID int64) (*model.WalletMember, error)
	GetForUpdate(dc *gorm.DB, memberID int64) (*model.WalletMember, error)
	Create(dc *gorm.DB, memberID int64) error
}

type generalUseCase struct {
}

func (uc *generalUseCase) Get(dc *gorm.DB, memberID int64) (*model.WalletMember, error) {
	memberWallet, err := dao.Wallet.First(dc, memberID)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, errs.WalletMemberNoFound
	}

	return memberWallet, nil
}

func (uc *generalUseCase) GetForUpdate(dc *gorm.DB, memberID int64) (*model.WalletMember, error) {
	memberWallet, err := dao.Wallet.FirstForUpdate(dc, memberID)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, errs.WalletMemberNoFound
	}

	return memberWallet, nil
}

func (uc *generalUseCase) Create(dc *gorm.DB, memberID int64) error {
	now := time.Now().UTC()

	walletMember := &model.WalletMember{
		MemberID:    memberID,
		Balance:     decimal.Zero,
		AddedTime:   now,
		UpdatedTime: now,
	}

	walletMember.Sign = signs.WalletMember(walletMember)

	if err := dao.Wallet.Insert(dc, walletMember); err != nil {
		logger.ApLog().Error(err)
		return errs.WalletMemberCreateFailed
	}

	return nil
}
