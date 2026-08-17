package transaction

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/models/model"
)

func newScanPayCommon() *scanPayCommonUsecase {
	return &scanPayCommonUsecase{scanPayOrderAdd: &scanPayCommonAdd{}}
}

type IScanPayCommon interface {
	OrderScanPayAdd(dc *gorm.DB, ctx context.Context, orderScanPay *model.OrderScanPay, mapping *model.ScanPayMapping) error
	GetOrderScanPayTransaction(dc *gorm.DB, orderScanPayID int64) (*model.Transaction, error)
}

type scanPayCommonUsecase struct {
	scanPayOrderAdd *scanPayCommonAdd
}

func (uc *scanPayCommonUsecase) OrderScanPayAdd(dc *gorm.DB, ctx context.Context, orderScanPay *model.OrderScanPay, mapping *model.ScanPayMapping) error {
	return uc.scanPayOrderAdd.Handler(dc, ctx, orderScanPay, mapping)
}

// GetOrderScanPayTransaction 只讀取已存在的錢包流水，reconciliation 不會再執行扣款。
func (uc *scanPayCommonUsecase) GetOrderScanPayTransaction(dc *gorm.DB, orderScanPayID int64) (*model.Transaction, error) {
	return dao.Transaction.FirstBySource(dc, model.TransactionSourceTypeScanPayConfirm, orderScanPayID)
}
