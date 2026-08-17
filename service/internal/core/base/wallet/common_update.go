package wallet

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"

	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/errs"
	"github.com/james730922/wallet/service/internal/utils/signs"
	"github.com/james730922/wallet/service/internal/utils/tools"
)

func newCommonUpdate() *commonUpdate {
	return &commonUpdate{}
}

type commonUpdate struct{}

// AddAmount 加額
func (w *commonUpdate) AddAmount(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error) {
	// 轉成 decimal 並校驗金額
	dCondAmount := cond.Amount
	if err := w.validateAmount(cond.MemberID, dCondAmount); err != nil {
		return nil, err
	}

	// 取得原錢包金額
	walletMember, err := w.getAndCheckForUpdate(dc, cond.MemberID)
	if err != nil {
		logger.ApLog().Errorf("[ADD] memberID: %d, err: %v", cond.MemberID, err)
		return nil, err
	}
	dWalletMember := w.convertWalletMember(walletMember)

	result := &model.WalletMemberModifyResult{}
	result.CurrentTotalAmount = walletMember.TotalAmount
	result.ChangedTotalAmount = dWalletMember.TotalAmount.Add(dCondAmount)
	result.CurrentBalance = walletMember.Balance
	result.ChangedBalance = dWalletMember.Balance.Add(dCondAmount)
	result.Amount = cond.Amount

	walletMember.Balance = dWalletMember.Balance.Add(dCondAmount)
	walletMember.TotalAmount = dWalletMember.TotalAmount.Add(dCondAmount)
	walletMember.Amount = dWalletMember.Amount.Add(dCondAmount)

	if err := w.validateWallet(walletMember); err != nil {
		if err != errs.WalletMemberUpdateBalanceIsNegative || err != errs.WalletMemberAmountUnreasonable {
			logger.ApLog().Errorf("validateWallet memberID: %d, err: %v", cond.MemberID, err)
		}
		return nil, err
	}

	err = w.update(dc, walletMember)
	if err != nil {
		logger.ApLog().Errorf("[ADD] memberID: %d, update val: %v", tools.JsonMarshalString(walletMember))
		return nil, err
	}

	return result, nil
}

// AddBonus 獲得獎勵金
func (w *commonUpdate) AddBonus(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error) {
	// 轉成 decimal 並校驗金額
	dCondAmount := cond.Amount
	if err := w.validateAmount(cond.MemberID, dCondAmount); err != nil {
		return nil, err
	}

	// 取得原錢包金額
	walletMember, err := w.getAndCheckForUpdate(dc, cond.MemberID)
	if err != nil {
		logger.ApLog().Errorf("[ADD] memberID: %d, err: %v", cond.MemberID, err)
		return nil, err
	}
	dWalletMember := w.convertWalletMember(walletMember)

	result := &model.WalletMemberModifyResult{}
	result.CurrentTotalAmount = walletMember.TotalAmount
	result.ChangedTotalAmount = dWalletMember.TotalAmount.Add(dCondAmount)
	result.CurrentBalance = walletMember.Balance
	result.ChangedBalance = dWalletMember.Balance.Add(dCondAmount)
	result.Amount = cond.Amount

	walletMember.Balance = dWalletMember.Balance.Add(dCondAmount)
	walletMember.TotalAmount = dWalletMember.TotalAmount.Add(dCondAmount)
	walletMember.Bonus = dWalletMember.Bonus.Add(dCondAmount)

	if err := w.validateWallet(walletMember); err != nil {
		if err != errs.WalletMemberUpdateBalanceIsNegative || err != errs.WalletMemberAmountUnreasonable {
			logger.ApLog().Errorf("validateWallet memberID: %d, err: %v", cond.MemberID, err)
		}
		return nil, err
	}

	err = w.update(dc, walletMember)
	if err != nil {
		logger.ApLog().Errorf("[ADD] memberID: %d, update val: %v", tools.JsonMarshalString(walletMember))
		return nil, err
	}

	return result, nil
}

// SubtractAmount 減少款項
func (w *commonUpdate) SubtractAmount(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error) {
	// 轉成 decimal 並校驗金額
	dCondAmount := cond.Amount
	if err := w.validateAmount(cond.MemberID, dCondAmount); err != nil {
		return nil, err
	}

	// 取得原錢包金額
	walletMember, err := w.getAndCheckForUpdate(dc, cond.MemberID)
	if err != nil {
		logger.ApLog().Errorf("[SUB] memberID: %d, err: %v", cond.MemberID, err)
		return nil, err
	}
	dWalletMember := w.convertWalletMember(walletMember)

	result := &model.WalletMemberModifyResult{}
	result.CurrentTotalAmount = walletMember.TotalAmount
	result.ChangedTotalAmount = dWalletMember.TotalAmount.Sub(dCondAmount)
	result.CurrentBalance = walletMember.Balance
	result.ChangedBalance = dWalletMember.Balance.Sub(dCondAmount)
	result.Amount = cond.Amount.Neg()

	walletMember.Balance = dWalletMember.Balance.Sub(dCondAmount)
	walletMember.TotalAmount = dWalletMember.TotalAmount.Sub(dCondAmount)
	walletMember.Amount = dWalletMember.Amount.Sub(dCondAmount)

	if err := w.validateWallet(walletMember); err != nil {
		if err != errs.WalletMemberUpdateBalanceIsNegative || err != errs.WalletMemberAmountUnreasonable {
			logger.ApLog().Errorf("validateWallet memberID: %d, err: %v", cond.MemberID, err)
		}
		return nil, err
	}

	err = w.update(dc, walletMember)
	if err != nil {
		logger.ApLog().Errorf("[SUB] memberID: %d, update val: %v", tools.JsonMarshalString(walletMember))
		return nil, err
	}

	return result, nil
}

// SubtractBonus 減少紅利
func (w *commonUpdate) SubtractBonus(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error) {
	// 轉成 decimal 並校驗金額
	dCondAmount := cond.Amount
	if err := w.validateAmount(cond.MemberID, dCondAmount); err != nil {
		return nil, err
	}

	// 取得原錢包金額
	walletMember, err := w.getAndCheckForUpdate(dc, cond.MemberID)
	if err != nil {
		logger.ApLog().Errorf("[SUB] memberID: %d, err: %v", cond.MemberID, err)
		return nil, err
	}
	dWalletMember := w.convertWalletMember(walletMember)

	result := &model.WalletMemberModifyResult{}
	result.CurrentTotalAmount = walletMember.TotalAmount
	result.ChangedTotalAmount = dWalletMember.TotalAmount.Sub(dCondAmount)
	result.CurrentBalance = walletMember.Balance
	result.ChangedBalance = dWalletMember.Balance.Sub(dCondAmount)
	result.Amount = cond.Amount.Neg()

	walletMember.Balance = dWalletMember.Balance.Sub(dCondAmount)
	walletMember.TotalAmount = dWalletMember.TotalAmount.Sub(dCondAmount)
	walletMember.Bonus = dWalletMember.Bonus.Sub(dCondAmount)

	if err := w.validateWallet(walletMember); err != nil {
		if err != errs.WalletMemberUpdateBalanceIsNegative || err != errs.WalletMemberAmountUnreasonable {
			logger.ApLog().Errorf("validateWallet memberID: %d, err: %v", cond.MemberID, err)
		}
		return nil, err
	}

	err = w.update(dc, walletMember)
	if err != nil {
		logger.ApLog().Errorf("[SUB] memberID: %d, update val: %v", tools.JsonMarshalString(walletMember))
		return nil, err
	}

	return result, nil
}

// FreezeAmount 凍結金額
func (w *commonUpdate) FreezeAmount(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error) {
	// 轉成 decimal 並校驗金額
	dCondAmount := cond.Amount
	if err := w.validateAmount(cond.MemberID, dCondAmount); err != nil {
		return nil, err
	}

	// 取得原錢包金額
	walletMember, err := w.getAndCheckForUpdate(dc, cond.MemberID)
	if err != nil {
		logger.ApLog().Errorf("memberID: %d, err: %v", cond.MemberID, err)
		return nil, err
	}
	dWalletMember := w.convertWalletMember(walletMember)

	result := &model.WalletMemberModifyResult{}
	result.CurrentTotalAmount = walletMember.TotalAmount
	result.ChangedTotalAmount = walletMember.TotalAmount
	result.CurrentBalance = walletMember.Balance
	result.ChangedBalance = dWalletMember.Balance.Sub(dCondAmount)
	result.Amount = cond.Amount

	walletMember.Balance = dWalletMember.Balance.Sub(dCondAmount)
	walletMember.FrozenAmount = dWalletMember.FrozenAmount.Add(dCondAmount)

	if err := w.validateWallet(walletMember); err != nil {
		if err != errs.WalletMemberUpdateBalanceIsNegative || err != errs.WalletMemberAmountUnreasonable {
			logger.ApLog().Errorf("validateWallet memberID: %d, err: %v", cond.MemberID, err)
		}
		return nil, err
	}

	err = w.update(dc, walletMember)
	if err != nil {
		logger.ApLog().Errorf("memberID: %d, update val: %v", tools.JsonMarshalString(walletMember))
		return nil, err
	}

	return result, nil
}

// ReleaseAmount 退回預扣款
func (w *commonUpdate) ReleaseAmount(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error) {
	// 轉成 decimal 並校驗金額
	dCondAmount := cond.Amount
	if err := w.validateAmount(cond.MemberID, dCondAmount); err != nil {
		return nil, err
	}

	// 取得原錢包金額
	walletMember, err := w.getAndCheckForUpdate(dc, cond.MemberID)
	if err != nil {
		logger.ApLog().Errorf("[SUB] memberID: %d, err: %v", cond.MemberID, err)
		return nil, err
	}
	dWalletMember := w.convertWalletMember(walletMember)

	result := &model.WalletMemberModifyResult{}
	result.CurrentTotalAmount = walletMember.TotalAmount
	result.ChangedTotalAmount = walletMember.TotalAmount
	result.CurrentBalance = walletMember.Balance
	result.ChangedBalance = dWalletMember.Balance.Add(dCondAmount)
	result.Amount = cond.Amount

	walletMember.Balance = dWalletMember.Balance.Add(dCondAmount)
	walletMember.FrozenAmount = dWalletMember.FrozenAmount.Sub(dCondAmount)

	if walletMember.FrozenAmount.IsNegative() {
		logger.ApLog().Errorf("[SUB] %s memberID: %d, amount: %s, frozen_amount: %s",
			errs.WalletMemberUpdateFrozenAmountIsNegative, cond.MemberID, cond.Amount, walletMember.FrozenAmount)
		return nil, errs.WalletMemberUpdateFrozenAmountIsNegative
	}

	if err := w.validateWallet(walletMember); err != nil {
		if err != errs.WalletMemberUpdateBalanceIsNegative || err != errs.WalletMemberAmountUnreasonable {
			logger.ApLog().Errorf("validateWallet memberID: %d, err: %v", cond.MemberID, err)
		}
		return nil, err
	}

	err = w.update(dc, walletMember)
	if err != nil {
		logger.ApLog().Errorf("memberID: %d, update val: %v", tools.JsonMarshalString(walletMember))
		return nil, err
	}

	return result, nil
}

// SubAmount 扣除凍結金額
func (w *commonUpdate) SubFreezeAmount(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error) {
	// 轉成 decimal 並校驗金額
	dCondAmount := cond.Amount
	if err := w.validateAmount(cond.MemberID, dCondAmount); err != nil {
		return nil, err
	}

	// 取得原錢包金額
	walletMember, err := w.getAndCheckForUpdate(dc, cond.MemberID)
	if err != nil {
		logger.ApLog().Errorf("[SUB] memberID: %d, err: %v", cond.MemberID, err)
		return nil, err
	}
	dWalletMember := w.convertWalletMember(walletMember)

	result := &model.WalletMemberModifyResult{}
	result.CurrentTotalAmount = walletMember.TotalAmount
	result.ChangedTotalAmount = dWalletMember.TotalAmount.Sub(dCondAmount)
	result.CurrentBalance = walletMember.Balance
	result.ChangedBalance = walletMember.Balance
	result.Amount = cond.Amount

	// 扣款邏輯，目前採有獎勵金優先扣除
	switch {
	case walletMember.Bonus.GreaterThanOrEqual(cond.Amount):
		result.UsedBonus = cond.Amount
		walletMember.Bonus = dWalletMember.Bonus.Sub(dCondAmount)
	case walletMember.Bonus.LessThan(cond.Amount) && walletMember.Bonus.IsPositive():
		result.UsedBonus = walletMember.Bonus
		walletMember.Amount = dWalletMember.Amount.Sub(dCondAmount).Add(dWalletMember.Bonus)
		walletMember.Bonus = decimal.Zero
	case walletMember.Bonus.IsZero():
		walletMember.Amount = dWalletMember.Amount.Sub(dCondAmount)
	default:
		logger.ApLog().Errorf("[SubAmount] %s memberID: %d, amount: %s, bonus: %s ",
			errs.WalletMemberUpdateBonusIsNegative, cond.MemberID, cond.Amount, walletMember.Bonus)
		return nil, errs.WalletMemberUpdateBonusIsNegative
	}

	walletMember.TotalAmount = dWalletMember.TotalAmount.Sub(dCondAmount)
	walletMember.FrozenAmount = dWalletMember.FrozenAmount.Sub(dCondAmount)

	if err := w.validateWallet(walletMember); err != nil {
		if err != errs.WalletMemberUpdateBalanceIsNegative || err != errs.WalletMemberAmountUnreasonable {
			logger.ApLog().Errorf("validateWallet memberID: %d, err: %v", cond.MemberID, err)
		}
		return nil, err
	}

	err = w.update(dc, walletMember)
	if err != nil {
		logger.ApLog().Errorf("memberID: %d, update val: %v", tools.JsonMarshalString(walletMember))
		return nil, err
	}

	return result, nil
}

// BuyPoint 購買遊戲幣
func (w *commonUpdate) BuyPoint(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error) {
	// 轉成 decimal 並校驗金額
	dCondAmount := cond.Amount
	if err := w.validateAmount(cond.MemberID, dCondAmount); err != nil {
		return nil, err
	}

	// 取得原錢包金額
	walletMember, err := w.getAndCheckForUpdate(dc, cond.MemberID)
	if err != nil {
		logger.ApLog().Errorf("[ADD] memberID: %d, err: %v", cond.MemberID, err)
		return nil, err
	}
	dWalletMember := w.convertWalletMember(walletMember)

	result := &model.WalletMemberModifyResult{}
	result.CurrentTotalAmount = walletMember.TotalAmount
	result.ChangedTotalAmount = dWalletMember.TotalAmount.Add(dCondAmount)
	result.CurrentBalance = walletMember.Balance
	result.ChangedBalance = dWalletMember.Balance.Add(dCondAmount)
	result.Amount = cond.Amount

	walletMember.Balance = dWalletMember.Balance.Add(dCondAmount)
	walletMember.TotalAmount = dWalletMember.TotalAmount.Add(dCondAmount)
	walletMember.Amount = dWalletMember.Amount.Add(dCondAmount)

	if err := w.validateWallet(walletMember); err != nil {
		if err != errs.WalletMemberUpdateBalanceIsNegative || err != errs.WalletMemberAmountUnreasonable {
			logger.ApLog().Errorf("validateWallet memberID: %d, err: %v", cond.MemberID, err)
		}
		return nil, err
	}

	err = w.update(dc, walletMember)
	if err != nil {
		logger.ApLog().Errorf("[ADD] memberID: %d, update val: %v", tools.JsonMarshalString(walletMember))
		return nil, err
	}

	return result, nil
}

// 掃碼支付
func (w *commonUpdate) ScanPaySubAmount(dc *gorm.DB, cond *condition.WalletMemberUpdateCond) (*model.WalletMemberModifyResult, error) {
	// 轉成 decimal 並校驗金額
	dCondAmount := cond.Amount
	if err := w.validateAmount(cond.MemberID, dCondAmount); err != nil {
		return nil, err
	}

	// 取得原錢包金額
	walletMember, err := w.getAndCheckForUpdate(dc, cond.MemberID)
	if err != nil {
		logger.ApLog().Errorf("[SUB] memberID: %d, err: %v", cond.MemberID, err)
		return nil, err
	}
	dWalletMember := w.convertWalletMember(walletMember)

	result := &model.WalletMemberModifyResult{}
	result.CurrentTotalAmount = walletMember.TotalAmount
	result.ChangedTotalAmount = dWalletMember.TotalAmount.Sub(dCondAmount)
	result.CurrentBalance = walletMember.Balance
	result.ChangedBalance = dWalletMember.Balance.Sub(dCondAmount)
	result.Amount = cond.Amount.Neg()

	// 扣款邏輯，目前採有獎勵金優先扣除
	switch {
	case walletMember.Bonus.GreaterThanOrEqual(cond.Amount):
		result.UsedBonus = cond.Amount
		walletMember.Bonus = dWalletMember.Bonus.Sub(dCondAmount)
	case walletMember.Bonus.LessThan(cond.Amount) && walletMember.Bonus.IsPositive():
		result.UsedBonus = walletMember.Bonus
		walletMember.Amount = dWalletMember.Amount.Sub(dCondAmount).Add(dWalletMember.Bonus)
		walletMember.Bonus = decimal.Zero
	case walletMember.Bonus.IsZero():
		walletMember.Amount = dWalletMember.Amount.Sub(dCondAmount)
	default:
		logger.ApLog().Errorf("[SubAmount] %s memberID: %d, amount: %s, bonus: %s ",
			errs.WalletMemberUpdateBonusIsNegative, cond.MemberID, cond.Amount, walletMember.Bonus)
		return nil, errs.WalletMemberUpdateBonusIsNegative
	}

	walletMember.Balance = dWalletMember.Balance.Sub(dCondAmount)
	walletMember.TotalAmount = dWalletMember.TotalAmount.Sub(dCondAmount)

	if err := w.validateWallet(walletMember); err != nil {
		if err != errs.WalletMemberUpdateBalanceIsNegative || err != errs.WalletMemberAmountUnreasonable {
			logger.ApLog().Errorf("validateWallet memberID: %d, err: %v", cond.MemberID, err)
		}
		return nil, err
	}

	err = w.update(dc, walletMember)
	if err != nil {
		logger.ApLog().Errorf("memberID: %d, update val: %v", tools.JsonMarshalString(walletMember))
		return nil, err
	}

	return result, nil
}

func (w *commonUpdate) getAndCheckForUpdate(dc *gorm.DB, memberID int64) (*model.WalletMember, error) {
	walletMember, err := common.getForUpdate(dc, memberID)
	if err != nil {
		return nil, err
	}

	// 驗證錢包簽名
	if walletMember.Sign != signs.WalletMember(walletMember) {
		return nil, errs.WalletMemberSignValidateFailed
	}

	return walletMember, nil
}

func (w *commonUpdate) update(dc *gorm.DB, walletMember *model.WalletMember) error {
	walletMember.UpdatedTime = time.Now().UTC()
	walletMember.Sign = signs.WalletMember(walletMember)

	updCond := condition.NewUpdate(
		&condition.WalletMemberUpdate{
			Balance:      tools.DecimalPtr(walletMember.Balance),
			TotalAmount:  tools.DecimalPtr(walletMember.TotalAmount),
			Amount:       tools.DecimalPtr(walletMember.Amount),
			Bonus:        tools.DecimalPtr(walletMember.Bonus),
			FrozenAmount: tools.DecimalPtr(walletMember.FrozenAmount),
			Sign:         aws.String(walletMember.Sign),
			UpdatedTime:  aws.Time(walletMember.UpdatedTime),
		},
		&condition.WalletMemberUpdate{
			MemberID: aws.Int64(walletMember.MemberID),
		},
	)

	err := dao.Wallet.Update(dc, updCond)
	if err != nil {
		return errs.WalletMemberUpdateFailed
	}
	if err := w.journal(dc, walletMember); err != nil {
		return errs.WalletMemberJournalCreateFailed
	}

	return nil
}

func (w *commonUpdate) validateAmount(memberID int64, amount decimal.Decimal) error {
	if err := w.checkDecimalPlaces(amount); err != nil {
		logger.ApLog().Debug(fmt.Sprintf("[amount validate failed] err: %s, memberID :%d, amount: %s\n%v", errs.CommonAmountDecimalPlacesError, memberID, amount, string(debug.Stack())))
		return err
	}
	return nil
}

func (w *commonUpdate) validateWallet(walletMember *model.WalletMember) error {
	memberID := walletMember.MemberID

	decimalWallet := w.convertWalletMember(walletMember)

	// 錢包金額不可小於 0
	validateNegative := func(wallet *model.WalletMemberTypeDecimal) bool {
		if wallet.Balance.IsNegative() ||
			wallet.Amount.IsNegative() ||
			wallet.TotalAmount.IsNegative() ||
			wallet.Bonus.IsNegative() ||
			wallet.FrozenAmount.IsNegative() {
			return false
		}
		return true
	}

	// 錢包金額位數驗證
	validateDecimalPlaces := func(wallet *model.WalletMemberTypeDecimal) bool {
		if err := w.checkDecimalPlaces(wallet.Balance); err != nil {
			return false
		}
		if err := w.checkDecimalPlaces(wallet.Amount); err != nil {
			return false
		}
		if err := w.checkDecimalPlaces(wallet.TotalAmount); err != nil {
			return false
		}
		if err := w.checkDecimalPlaces(wallet.Bonus); err != nil {
			return false
		}
		if err := w.checkDecimalPlaces(wallet.FrozenAmount); err != nil {
			return false
		}

		return true
	}

	// 錢包 balance 驗證
	validateRelation := func(wallet *model.WalletMemberTypeDecimal) bool {
		// Balance = Amount + Bonus - FrozenAmount
		if !wallet.Balance.Equal(wallet.Amount.Add(wallet.Bonus).Sub(wallet.FrozenAmount)) {
			return false
		}

		// TotalAmount = Balance + FrozenAmount
		if !wallet.TotalAmount.Equal(wallet.Balance.Add(wallet.FrozenAmount)) {
			return false
		}

		// TotalAmount = Amount + Bonus
		if !wallet.TotalAmount.Equal(wallet.Amount.Add(wallet.Bonus)) {
			return false
		}

		return true
	}

	if !validateNegative(decimalWallet) {
		logger.ApLog().Debug(fmt.Sprintf("[wallet validate failed] err: %s, memberID :%d, wallet: %s\n%v", errs.WalletMemberUpdateBalanceIsNegative, memberID, decimalWallet, string(debug.Stack())))
		return errs.WalletMemberUpdateBalanceIsNegative
	}

	if !validateDecimalPlaces(decimalWallet) {
		logger.ApLog().Debug(fmt.Sprintf("[wallet validate failed] err: %s, memberID :%d, wallet: %s\n%v", errs.CommonAmountDecimalPlacesError, memberID, decimalWallet, string(debug.Stack())))
		return errs.CommonAmountDecimalPlacesError
	}

	if !validateRelation(decimalWallet) {
		logger.ApLog().Debug(fmt.Sprintf("[wallet validate failed] err: %s, memberID :%d, wallet: %s\n%v", errs.WalletMemberAmountUnreasonable, memberID, decimalWallet, string(debug.Stack())))
		return errs.WalletMemberAmountUnreasonable
	}

	return nil
}

func (w *commonUpdate) convertWalletMember(walletMember *model.WalletMember) *model.WalletMemberTypeDecimal {
	return model.ConvertWalletMember(walletMember)
}

func (w *commonUpdate) checkDecimalPlaces(d decimal.Decimal) error {
	if !d.Round(2).Equal(d) {
		return errs.CommonAmountDecimalPlacesError
	}

	return nil
}

func (w *commonUpdate) journal(dc *gorm.DB, walletMember *model.WalletMember) error {
	journal := &model.JournalMemberWallet{
		ID:           packet.Node.Generate().Int64(),
		MemberID:     walletMember.MemberID,
		Balance:      walletMember.Balance,
		TotalAmount:  walletMember.TotalAmount,
		Amount:       walletMember.Amount,
		Bonus:        walletMember.Bonus,
		FrozenAmount: walletMember.FrozenAmount,
		WalletSign:   walletMember.Sign,
		AddedTime:    walletMember.UpdatedTime,
	}

	if err := dao.Journal.Insert(dc, journal); err != nil {
		logger.ApLog().Error(err)
		return err
	}

	return nil
}

func (w *commonUpdate) PointInternal(dc *gorm.DB, cond *condition.WalletMemberUpdateCond, isCancel bool) (*model.WalletMemberModifyResult, error) {
	// 轉成 decimal 並校驗金額
	dCondAmount := cond.Amount
	if err := w.validateAmount(cond.MemberID, dCondAmount); err != nil {
		return nil, err
	}

	// 取得原錢包金額
	walletMember, err := w.getAndCheckForUpdate(dc, cond.MemberID)
	if err != nil {
		logger.ApLog().Errorf("[ADD] memberID: %d, err: %v", cond.MemberID, err)
		return nil, err
	}
	dWalletMember := w.convertWalletMember(walletMember)

	result := &model.WalletMemberModifyResult{}
	if isCancel {
		result.CurrentTotalAmount = walletMember.TotalAmount
		result.ChangedTotalAmount = walletMember.TotalAmount
		result.CurrentBalance = walletMember.Balance
		result.ChangedBalance = dWalletMember.Balance.Add(dCondAmount)
		result.Amount = cond.Amount

		walletMember.Balance = dWalletMember.Balance.Add(dCondAmount)
		walletMember.FrozenAmount = dWalletMember.FrozenAmount.Sub(dCondAmount)

		if walletMember.FrozenAmount.IsNegative() {
			logger.ApLog().Errorf("[SUB] %s memberID: %d, amount: %s, frozen_amount: %s",
				errs.WalletMemberUpdateFrozenAmountIsNegative, cond.MemberID, cond.Amount, walletMember.FrozenAmount)
			return nil, errs.WalletMemberUpdateFrozenAmountIsNegative
		}
	} else {

		result.CurrentTotalAmount = walletMember.TotalAmount
		result.ChangedTotalAmount = dWalletMember.TotalAmount.Sub(dCondAmount)
		result.CurrentBalance = walletMember.Balance
		result.ChangedBalance = walletMember.Balance
		result.Amount = cond.Amount

		switch {
		case walletMember.Bonus.GreaterThanOrEqual(cond.Amount):
			result.UsedBonus = cond.Amount
			walletMember.Bonus = dWalletMember.Bonus.Sub(dCondAmount)
		case walletMember.Bonus.LessThan(cond.Amount) && walletMember.Bonus.IsPositive():
			result.UsedBonus = walletMember.Bonus
			walletMember.Amount = dWalletMember.Amount.Sub(dCondAmount).Add(dWalletMember.Bonus)
			walletMember.Bonus = decimal.Zero
		case walletMember.Bonus.IsZero():
			walletMember.Amount = dWalletMember.Amount.Sub(dCondAmount)
		default:
			logger.ApLog().Errorf("[SubAmount] %s memberID: %d, amount: %s, bonus: %s ",
				errs.WalletMemberUpdateBonusIsNegative, cond.MemberID, cond.Amount, walletMember.Bonus)
			return nil, errs.WalletMemberUpdateBonusIsNegative
		}

		walletMember.TotalAmount = dWalletMember.TotalAmount.Sub(dCondAmount)
		walletMember.FrozenAmount = dWalletMember.FrozenAmount.Sub(dCondAmount)
	}

	if err := w.validateWallet(walletMember); err != nil {
		if err != errs.WalletMemberUpdateBalanceIsNegative || err != errs.WalletMemberAmountUnreasonable {
			logger.ApLog().Errorf("validateWallet memberID: %d, err: %v", cond.MemberID, err)
		}
		return nil, err
	}

	err = w.update(dc, walletMember)
	if err != nil {
		logger.ApLog().Errorf("[ADD] memberID: %d, update val: %v", tools.JsonMarshalString(walletMember))
		return nil, err
	}

	return result, nil
}
