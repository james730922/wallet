package scanpay

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/errs"
	"github.com/james730922/wallet/service/internal/utils/signs"
	"github.com/james730922/wallet/service/internal/utils/tools"
)

type scanPayCommonOrder struct{}

func (r *scanPayCommonOrder) Handler(ctx context.Context, cond *condition.OrderScanPayCreateCond) (int64, error) {
	// 取得單號
	newOrderID := packet.Node.Generate().Int64()
	var orderScanPayID int64

	tx := func(dc *gorm.DB) error {
		id, create, err := r.claimOrder(dc, cond, newOrderID)
		if err != nil {
			logger.ApLog().Error(err)
			return err
		}
		orderScanPayID = id
		if !create {
			return nil
		}

		// 新增關聯單
		_, err = r.order(dc, cond, orderScanPayID)
		if err != nil {
			logger.ApLog().Errorf("err: %s, cond: %s", err, tools.JsonMarshalString(cond))
			return errs.OrderDepositCreateFailed
		}

		return nil
	}

	if err := packet.DB.Transaction(tx); err != nil {
		return 0, err
	}

	return orderScanPayID, nil
}

func (r *scanPayCommonOrder) claimOrder(dc *gorm.DB, cond *condition.OrderScanPayCreateCond, newOrderID int64) (int64, bool, error) {
	if cond.RecordID == nil {
		return 0, false, errs.CommonRequestParamInvalid
	}
	now := time.Now().UTC()

	record, err := self.Record.GetScanPayRecordForUpdate(dc, *cond.RecordID)
	if err != nil {
		logger.ApLog().Error(err)
		return 0, false, err
	}

	if !record.Amount.Equal(cond.Amount) {
		return 0, false, errs.CommonRequestParamInvalid
	}

	switch record.Status {
	case model.ScanPayRecordStatusWaiting:
		// 繼續執行
	case model.ScanPayRecordStatusTransaction:
		// record row is already locked and record_id is unique. Do not lock order
		// here: payment locks order then record, so doing the reverse would deadlock.
		order, err := dao.ScanPayDAO.FirstByRecordID(dc, record.ID)
		if err != nil {
			if err == errs.DBNoRow {
				return 0, false, errs.ScanPayRecordStatusError
			}
			return 0, false, err
		}
		id, err := reusableOrderID(order, cond)
		return id, false, err
	case model.ScanPayRecordStatusDone:
		return 0, false, errs.ScanPayOrderAlreadyDone
	case model.ScanPayRecordStatusCancel:
		return 0, false, errs.ScanPayOrderCancel
	case model.ScanPayRecordStatusFailure:
		// Failure record 已永久作廢，不得重用舊 order 或建立新 order。
		return 0, false, errs.ScanPayOrderFailure
	default:
		return 0, false, errs.ScanPayRecordStatusError
	}
	//檢查是否過期
	if record.ExpiredTime.Before(now) {
		return 0, false, errs.ScanPayRecordExpired
	}
	//改成交易中
	updateCond := condition.NewUpdate(&condition.ScanPayRecordUpdate{
		Status:      aws.Int(model.ScanPayRecordStatusTransaction),
		UpdatedTime: aws.Time(now),
	}, &condition.ScanPayRecordUpdate{
		ID: &record.ID,
	})
	if err := dao.Record.Update(dc, updateCond); err != nil {
		logger.SysLog().Warnf("err:%v, cond:%v", err, tools.JsonMarshalString(updateCond))
		return 0, false, err
	}

	return newOrderID, true, nil
}

func reusableOrderID(order *model.OrderScanPay, cond *condition.OrderScanPayCreateCond) (int64, error) {
	if order == nil || order.RecordID == nil || cond.RecordID == nil ||
		*order.RecordID != *cond.RecordID ||
		order.MemberID != cond.MemberID ||
		!order.Amount.Equal(cond.Amount) {
		return 0, errs.ScanPayAddRecordFailed
	}

	switch order.Status {
	case model.OrderScanPayStatusEnumTransaction:
		return order.ID, nil
	case model.OrderScanPayStatusEnumSuccess:
		return 0, errs.ScanPayOrderAlreadyDone
	case model.OrderScanPayStatusEnumCancel:
		return 0, errs.ScanPayOrderCancel
	case model.OrderScanPayStatusEnumFailure:
		// Failure order 已永久作廢，不得重試。
		return 0, errs.ScanPayOrderFailure
	default:
		return 0, errs.ScanPayAddRecordFailed
	}
}

func (r *scanPayCommonOrder) order(dc *gorm.DB, cond *condition.OrderScanPayCreateCond, orderScanPayID int64) (*model.OrderScanPay, error) {
	now := time.Now().UTC()
	orderScanPay := &model.OrderScanPay{
		ID:              orderScanPayID,
		MemberID:        cond.MemberID,
		Amount:          cond.Amount,
		Status:          model.OrderScanPayStatusEnumTransaction,
		RecordID:        cond.RecordID,
		Brand:           aws.StringValue(cond.Brand),
		MerchantOrderID: aws.StringValue(cond.MerchantOrderID),
		SourceOrderID:   aws.StringValue(cond.SourceOrderID),
		Content:         aws.StringValue(cond.Content),
		AddedTime:       now,
		UpdatedTime:     now,
	}

	orderScanPay.Sign = signs.OrderScanPay(orderScanPay)

	if err := dao.ScanPayDAO.Insert(dc, orderScanPay); err != nil {
		return nil, err
	}

	return orderScanPay, nil
}
