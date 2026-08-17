package scanpay

import (
	"context"

	"github.com/james730922/wallet/service/internal/models/condition"
)

func newScanPayCommon() IScanPayCommon {
	return &scanPayCommonUseCase{}
}

type IScanPayCommon interface {
	CreateOrder(ctx context.Context, cond *condition.OrderScanPayCreateCond) (int64, error)
	OrderPay(ctx context.Context, cond *condition.OrderScanPayToPayCond) (int64, error)
}

type scanPayCommonUseCase struct {
	add *scanPayCommonOrder
	pay *scanPayCommonOrderPay
}

// 下單
func (uc *scanPayCommonUseCase) CreateOrder(ctx context.Context, cond *condition.OrderScanPayCreateCond) (int64, error) {
	return uc.add.Handler(ctx, cond)
}

// 付款
func (uc *scanPayCommonUseCase) OrderPay(ctx context.Context, cond *condition.OrderScanPayToPayCond) (int64, error) {
	return uc.pay.Handler(ctx, cond)
}
