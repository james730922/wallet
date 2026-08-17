package scanpay

import (
	"context"
	"errors"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/event"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/thirdparty/observability"
	"github.com/james730922/wallet/service/internal/utils/errs"
	"github.com/james730922/wallet/service/internal/utils/signs"
)

const (
	scanPayReconcileInterval = time.Minute
	scanPayTransactionStale  = 10 * time.Minute
	scanPayReconcileBatch    = 100
)

type IReconciler interface {
	Run(ctx context.Context)
	ReconcileOnce(ctx context.Context) error
}

type scanPayReconciler struct {
	orderPay *scanPayCommonOrderPay
}

func newScanPayReconciler() IReconciler {
	return &scanPayReconciler{orderPay: &scanPayCommonOrderPay{}}
}

// Run 定期巡檢卡在 Transaction 的掃碼單。巡檢只核對已存在的 wallet ledger，絕不重新扣款。
func (r *scanPayReconciler) Run(ctx context.Context) {
	r.reconcileAndLog(ctx)
	ticker := time.NewTicker(scanPayReconcileInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.reconcileAndLog(ctx)
		}
	}
}

func (r *scanPayReconciler) reconcileAndLog(ctx context.Context) {
	if err := r.ReconcileOnce(ctx); err != nil {
		logger.ApLog().Errorf("scanpay reconciliation failed: %v", err)
	}
}

func (r *scanPayReconciler) ReconcileOnce(ctx context.Context) error {
	ids, err := dao.ScanPayDAO.ListStaleTransactionIDs(
		packet.DB.New(),
		time.Now().UTC().Add(-scanPayTransactionStale),
		scanPayReconcileBatch,
	)
	if err != nil {
		return err
	}

	var firstErr error
	for _, id := range ids {
		memberID, walletChanged, reconcileErr := r.reconcileOne(id)
		if reconcileErr != nil {
			observability.RecordScanPayTransactionFailure("reconciliation", reconcileErr)
			logger.ApLog().Errorf("scanpay reconciliation order_id=%d failed: %v", id, reconcileErr)
			if firstErr == nil {
				firstErr = reconcileErr
			}
			continue
		}
		if walletChanged {
			event.Event().Publish(event.TopicWalletMemberChange, model.WalletMemberWithChangeNotify{MemberID: memberID})
		}
	}

	return firstErr
}

func (r *scanPayReconciler) reconcileOne(orderID int64) (memberID int64, walletChanged bool, err error) {
	err = packet.DB.Transaction(func(dc *gorm.DB) error {
		order, txErr := dao.ScanPayDAO.FirstByIDForUpdate(dc, orderID)
		if txErr != nil {
			return txErr
		}
		if order.Status != model.OrderScanPayStatusEnumTransaction {
			return nil
		}
		if order.Sign != signs.OrderScanPay(order) || order.RecordID == nil {
			return errs.OrderScanPaySignValidateFailed
		}

		record, txErr := dao.Record.FirstForUpdate(dc, *order.RecordID)
		if txErr != nil {
			return txErr
		}
		if record.Status != model.ScanPayRecordStatusTransaction {
			// 狀態不一致時不自動改寫，保留人工查核，避免錯帳。
			return errs.ScanPayRecordStatusError
		}

		walletTransaction, txErr := packet.Transaction.GetOrderScanPayTransaction(dc, order.ID)
		switch {
		case txErr == nil:
			if walletTransaction.MemberID != order.MemberID || walletTransaction.Sign != signs.Transaction(walletTransaction) {
				return errs.WalletMemberSignValidateFailed
			}
			if txErr = r.orderPay.updateOrderStatus(dc, order, model.OrderScanPayStatusEnumSuccess, journalContentSuccess); txErr != nil {
				return txErr
			}
			if txErr = r.orderPay.updateRecordStatus(dc, record, model.ScanPayRecordStatusDone, ""); txErr != nil {
				return txErr
			}
			memberID = order.MemberID
			walletChanged = true
			return nil
		case errors.Is(txErr, errs.DBNoRow):
			// 逾時且無 wallet ledger 代表扣款未成立。Failure order 直接永久作廢，不得重試。
			return r.orderPay.finalizePaymentFailure(dc, order, record, journalContentFailure+"：交易逾時且無錢包流水")
		default:
			return txErr
		}
	})

	return memberID, walletChanged, err
}
