package scanpay

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/event"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/thirdparty/observability"
	"github.com/james730922/wallet/service/internal/utils/errs"
	"github.com/james730922/wallet/service/internal/utils/signs"
	"github.com/james730922/wallet/service/internal/utils/tools"
)

type scanPayCommonOrderPay struct{}

func (r *scanPayCommonOrderPay) Handler(ctx context.Context, cond *condition.OrderScanPayToPayCond) (int64, error) {

	orderScanPay := &model.OrderScanPay{}
	remarks := ""
	var paymentErr error
	tx := func(dc *gorm.DB) error {
		// 取得掃碼支付單
		order, err := dao.ScanPayDAO.FirstForUpdate(dc, cond)
		if err != nil {
			logger.ApLog().Error(err)
			return errs.OrderScanPaySignValidateFailed
		}
		// 檢查簽名
		if order.Sign != signs.OrderScanPay(order) {
			return errs.OrderScanPaySignValidateFailed
		}
		orderScanPay = order

		// 確認 order 單並進行扣款
		switch order.Status {
		case model.OrderScanPayStatusEnumTransaction:
			// 繼續進行
		case model.OrderScanPayStatusEnumCancel:
			return errs.ScanPayOrderCancel
		case model.OrderScanPayStatusEnumSuccess:
			return errs.ScanPayOrderAlreadyDone
		case model.OrderScanPayStatusEnumFailure:
			// Failure order 已永久作廢，不得再次扣款。
			return errs.ScanPayOrderFailure
		default:
			return errs.ScanPayRecordStatusError
		}

		// 再次確認 record 狀態並 lock
		record, err := dao.Record.FirstForUpdate(dc, aws.Int64Value(order.RecordID))
		if err != nil {
			logger.ApLog().Error(err)
			return errs.OrderScanPaySignValidateFailed
		}
		switch record.Status {
		case model.ScanPayRecordStatusWaiting:
			return errs.ScanPayRecordStatusError
		case model.ScanPayRecordStatusTransaction:
			// 繼續執行
		case model.ScanPayRecordStatusCancel:
			remarks = r.errMessage(errs.ScanPayOrderCancel)
			return errs.ScanPayOrderCancel
		case model.ScanPayRecordStatusDone:
			remarks = r.errMessage(errs.ScanPayOrderAlreadyDone)
			return errs.ScanPayOrderAlreadyDone
		case model.ScanPayRecordStatusFailure:
			return errs.ScanPayOrderFailure
		default:
			return errs.ScanPayRecordStatusError
		}

		// 錢包餘額扣除
		if err := r.transaction(dc, ctx, order); err != nil {
			if !errors.Is(err, errs.WalletMemberUpdateBalanceIsNegative) && !errors.Is(err, errs.WalletMemberAmountUnreasonable) {
				logger.ApLog().Errorf("err: %s, cond: %s", err, tools.JsonMarshalString(cond))
			}
			// 備註掃碼失敗原因
			remarks = r.errMessage(err)
			if !isTerminalPaymentFailure(err) {
				// DB、簽名或 ledger 等系統不確定錯誤不可擅自宣告付款失敗。
				return err
			}

			// Failure order 直接永久作廢，不得重試或重用。
			// order_scanpay 與 scanpay_record 必須在同一 transaction 中一起轉為 Failure。
			if finalizeErr := r.finalizePaymentFailure(dc, order, record, remarks); finalizeErr != nil {
				return finalizeErr
			}
			paymentErr = err
			return nil
		}

		//更新狀態訂單為成功
		if err := r.updateOrderStatus(dc, order, model.OrderScanPayStatusEnumSuccess, remarks); err != nil {
			logger.ApLog().Errorf("err: %s, orderId: %d", err, order.ID)
			return err
		}

		//更新全部掃碼狀態為已確認
		if err := r.updateRecordStatus(dc, record, model.ScanPayRecordStatusDone, ""); err != nil {
			logger.ApLog().Errorf("err: %s, orderId: %d", err, order.ID)
			return err
		}

		return nil
	}

	if err := packet.DB.Transaction(tx); err != nil {
		observability.RecordScanPayTransactionFailure("payment_transaction", err)
		return 0, err
	}
	if paymentErr != nil {
		return 0, paymentErr
	}

	event.Event().Publish(event.TopicWalletMemberChange, model.WalletMemberWithChangeNotify{MemberID: orderScanPay.MemberID})
	return orderScanPay.ID, nil
}
func (r *scanPayCommonOrderPay) transaction(dc *gorm.DB, ctx context.Context, orderScanPay *model.OrderScanPay) error {
	mapping, err := r.getScanPayMapping(dc, orderScanPay.RecordID)
	if err != nil {
		return err
	}

	return packet.Transaction.OrderScanPayAdd(dc, ctx, orderScanPay, mapping)
}

func (r *scanPayCommonOrderPay) getScanPayMapping(dc *gorm.DB, recordID *int64) (*model.ScanPayMapping, error) {
	mapping, err := self.Record.GetScanPayMapping(dc, &condition.ScanPayMappingQuery{
		RecordID: recordID,
	})
	if err != nil {
		logger.ApLog().Error(err)
		return nil, err
	}

	return mapping, nil
}

func (r *scanPayCommonOrderPay) updateOrderStatus(dc *gorm.DB, orderScanPay *model.OrderScanPay, orderEnum model.OrderScanPayStatusEnum, remarks string) error {
	if orderScanPay == nil || orderScanPay.Status != model.OrderScanPayStatusEnumTransaction {
		return errs.ScanPayRecordStatusError
	}
	now := time.Now().UTC()
	updatedOrder := *orderScanPay
	updatedOrder.Status = orderEnum
	updatedOrder.UpdatedTime = now
	updatedOrder.Remarks = &remarks
	updatedOrder.Sign = signs.OrderScanPay(&updatedOrder)

	updateCond := &condition.OrderScanPayUpdate{
		Status:      aws.Int(int(orderEnum)),
		Sign:        aws.String(updatedOrder.Sign),
		UpdatedTime: aws.Time(updatedOrder.UpdatedTime),
		Remarks:     &remarks,
	}

	if orderEnum == model.OrderScanPayStatusEnumSuccess {
		updateCond.SuccessTime = aws.Time(now)
		updatedOrder.SuccessTime = aws.Time(now)
	}

	currentStatus := int(model.OrderScanPayStatusEnumTransaction)
	cond := condition.NewUpdate(
		updateCond,
		&condition.OrderScanPayUpdate{
			ID:     &orderScanPay.ID,
			Status: &currentStatus,
		})

	if err := dao.ScanPayDAO.Update(dc, cond); err != nil {
		return err
	}
	*orderScanPay = updatedOrder

	return nil
}

func (r *scanPayCommonOrderPay) updateRecordStatus(dc *gorm.DB, record *model.ScanPayRecord, status model.ScanPayRecordStatus, remarks string) error {
	if record == nil || record.Status != model.ScanPayRecordStatusTransaction {
		return errs.ScanPayRecordStatusError
	}
	now := time.Now().UTC()
	update := &condition.ScanPayRecordUpdate{
		Status:      aws.Int(status),
		UpdatedTime: aws.Time(now),
	}
	if remarks != "" {
		update.Remarks = &remarks
	}
	currentStatus := model.ScanPayRecordStatusTransaction
	cond := condition.NewUpdate(update, &condition.ScanPayRecordUpdate{
		ID:     &record.ID,
		Status: &currentStatus,
	})

	if err := dao.Record.Update(dc, cond); err != nil {
		return err
	}
	record.Status = status
	record.UpdatedTime = now
	if remarks != "" {
		record.Remarks = &remarks
	}

	return nil
}

func (r *scanPayCommonOrderPay) finalizePaymentFailure(dc *gorm.DB, order *model.OrderScanPay, record *model.ScanPayRecord, remarks string) error {
	if order == nil || record == nil || order.RecordID == nil || *order.RecordID != record.ID {
		return errs.ScanPayRecordStatusError
	}

	// 只有 ledger 不存在時才能宣告失敗，避免已扣款卻把訂單作廢。
	if _, err := packet.Transaction.GetOrderScanPayTransaction(dc, order.ID); err == nil {
		return errs.ScanPayRecordStatusError
	} else if !errors.Is(err, errs.DBNoRow) {
		return err
	}

	if err := r.updateOrderStatus(dc, order, model.OrderScanPayStatusEnumFailure, remarks); err != nil {
		return err
	}
	return r.updateRecordStatus(dc, record, model.ScanPayRecordStatusFailure, remarks)
}

func isTerminalPaymentFailure(err error) bool {
	return errors.Is(err, errs.WalletMemberUpdateBalanceIsNegative) ||
		errors.Is(err, errs.CommonAmountDecimalPlacesError)
}

func (r *scanPayCommonOrderPay) errMessage(err error) string {
	errorMsg := ""
	var errorCode errs.ErrorCode
	e, ok := errs.Parse(err)
	if ok {
		errorCode = e.GetCode()
	}

	switch errorCode {
	case errs.WalletMemberSignValidateFailed.GetCode():
		errorMsg = "钱包签名验证失败"
	case errs.WalletMemberUpdateBalanceIsNegative.GetCode():
		errorMsg = "余额不足"
	case errs.ScanPayOrderAlreadyDone.GetCode():
		errorMsg = "已被完成"
	case errs.ScanPayOrderCancel.GetCode():
		errorMsg = "已被取消"
	case errs.CommonAmountDecimalPlacesError.GetCode():
		errorMsg = "金额小数位错误"
	}
	return errorMsg
}
