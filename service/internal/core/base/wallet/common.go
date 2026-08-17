package wallet

import (
	"github.com/jinzhu/gorm"

	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/thirdparty/observability"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

var (
	common *commonUseCase
)

func newCommon() ICommon {
	common = &commonUseCase{
		update: newCommonUpdate(),
	}
	return common
}

type ICommon interface {
	/*入款*/
	//  AddAmount 入款
	AddAmount(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error)

	/*獎勵*/
	//  AddBonus 獲得獎勵金
	AddBonus(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error)

	/*人工調整帳務*/
	// SubtractAmount 減少款項
	SubtractAmount(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error)
	// SubtractBonus 減少紅利
	SubtractBonus(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error)

	// FreezeAmount 凍結金額
	FreezeAmount(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error)
	// ReleaseAmount 退回預扣款
	ReleaseAmount(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error)
	// SubAmount 扣除凍結金額
	SubFreezeAmount(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error)

	// BuyPoint 購買遊戲幣
	BuyPoint(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error)

	/*掃碼支付*/
	ScanPaySubAmount(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error)

	/*更新賣家錢包 (wellPay 使用)*/
	PointInternal(dc *gorm.DB, cond *condition.WalletMemberUpdateCond, isCancel bool) (*model.WalletMemberModifyResult, error)
}

type commonUseCase struct {
	update *commonUpdate
}

func (w *commonUseCase) getForUpdate(dc *gorm.DB, memberID int64) (*model.WalletMember, error) {
	memberWallet, err := dao.Wallet.FirstForUpdate(dc, memberID)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, errs.WalletMemberNoFound
	}

	return memberWallet, nil
}

func (w *commonUseCase) AddAmount(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error) {
	result, err := w.update.AddAmount(dc, cond)
	recordWalletFailure("add_amount", err)
	return result, err
}

func (w *commonUseCase) AddBonus(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error) {
	result, err := w.update.AddBonus(dc, cond)
	recordWalletFailure("add_bonus", err)
	return result, err
}

func (w *commonUseCase) SubtractAmount(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error) {
	result, err := w.update.SubtractAmount(dc, cond)
	recordWalletFailure("subtract_amount", err)
	return result, err
}
func (w *commonUseCase) SubtractBonus(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error) {
	result, err := w.update.SubtractBonus(dc, cond)
	recordWalletFailure("subtract_bonus", err)
	return result, err
}

func (w *commonUseCase) FreezeAmount(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error) {
	result, err := w.update.FreezeAmount(dc, cond)
	recordWalletFailure("freeze_amount", err)
	return result, err
}

func (w *commonUseCase) ReleaseAmount(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error) {
	result, err := w.update.ReleaseAmount(dc, cond)
	recordWalletFailure("release_amount", err)
	return result, err
}

func (w *commonUseCase) SubFreezeAmount(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error) {
	result, err := w.update.SubFreezeAmount(dc, cond)
	recordWalletFailure("subtract_frozen_amount", err)
	return result, err
}

func (w *commonUseCase) BuyPoint(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error) {
	result, err := w.update.BuyPoint(dc, cond)
	recordWalletFailure("buy_point", err)
	return result, err
}

func (w *commonUseCase) ScanPaySubAmount(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error) {
	result, err := w.update.ScanPaySubAmount(dc, cond)
	recordWalletFailure("scanpay_debit", err)
	return result, err
}

func (w *commonUseCase) PointInternal(dc *gorm.DB, cond *condition.WalletMemberUpdateCond, isCancel bool) (*model.WalletMemberModifyResult, error) {
	result, err := w.update.PointInternal(dc, cond, isCancel)
	recordWalletFailure("point_internal", err)
	return result, err
}

func recordWalletFailure(operation string, err error) {
	if err != nil {
		observability.RecordWalletTransactionFailure(operation, err)
	}
}
