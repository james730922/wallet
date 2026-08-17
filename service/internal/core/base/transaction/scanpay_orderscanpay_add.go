package transaction

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/thirdparty/observability"
	"github.com/james730922/wallet/service/internal/utils/errs"
	"github.com/james730922/wallet/service/internal/utils/signs"
)

type scanPayCommonAdd struct{}

func (de *scanPayCommonAdd) Handler(dc *gorm.DB, ctx context.Context, orderScanPay *model.OrderScanPay, mapping *model.ScanPayMapping) error {
	walletMemberResult, err := de.updateWalletMember(dc, orderScanPay)
	if err != nil {
		observability.RecordScanPayTransactionFailure("wallet_debit", err)
		if err != errs.WalletMemberUpdateBalanceIsNegative || err != errs.WalletMemberAmountUnreasonable {
			logger.ApLog().Errorf("updateWalletMember err: %s", err)
		}
		return err
	}

	err = de.transaction(dc, ctx, orderScanPay, walletMemberResult, mapping)
	if err != nil {
		observability.RecordScanPayTransactionFailure("ledger_insert", err)
		logger.ApLog().Errorf("transaction err: %s", err)
		return err
	}

	return nil
}

func (de *scanPayCommonAdd) updateWalletMember(dc *gorm.DB, orderScanPay *model.OrderScanPay) (*model.WalletMemberModifyResult, error) {
	cond := &condition.WalletMemberUpdateCond{
		MemberID: orderScanPay.MemberID,
		Amount:   orderScanPay.Amount,
	}

	result, err := packet.Wallet.ScanPaySubAmount(dc, cond)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (de *scanPayCommonAdd) transaction(dc *gorm.DB, ctx context.Context, orderScanPay *model.OrderScanPay, walletMemberResult *model.WalletMemberModifyResult, mapping *model.ScanPayMapping) error {
	txnID := packet.Node.Generate().Int64()

	now := time.Now().UTC()

	orderTraction := &model.Transaction{
		ID:                    txnID,
		MemberID:              orderScanPay.MemberID,
		SourceType:            model.TransactionSourceTypeScanPayConfirm,
		SourceID:              orderScanPay.ID,
		CurrencyCode:          "RMB",
		Amount:                walletMemberResult.Amount,
		CurrentTotalAmount:    walletMemberResult.CurrentTotalAmount,
		ChangedTotalAmount:    walletMemberResult.ChangedTotalAmount,
		CurrentBalance:        walletMemberResult.CurrentBalance,
		ChangedBalance:        walletMemberResult.ChangedBalance,
		AddedTime:             now,
		UpdatedTime:           now,
		Merchant:              mapping.Merchant,
		MerchantMemberAccount: mapping.MerchantMemberAccount,
	}

	orderTraction.Sign = signs.Transaction(orderTraction)

	err := self.Common.insertTransaction(dc, ctx, orderTraction)
	if err != nil {
		return err
	}

	return nil
}
