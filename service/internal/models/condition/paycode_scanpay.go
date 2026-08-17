package condition

import (
	"time"

	"github.com/shopspring/decimal"
)

type PayCodeScanPay struct {
	Expired *time.Time
	ID      int64
	Amount  decimal.Decimal
}
