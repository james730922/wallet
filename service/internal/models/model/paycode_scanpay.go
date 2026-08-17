package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type PayCodeScanPay struct {
	ID      int64
	Amount  decimal.Decimal
	Expired *time.Time
}
