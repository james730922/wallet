package condition

import (
	"time"

	"github.com/shopspring/decimal"
)

type OrderScanPayUpdate struct {
	ID          *int64     `json:"id"`           // 掃碼訂單號
	Status      *int       `json:"status"`       // 訂單狀態
	SuccessTime *time.Time `json:"success_time"` // 成功時間
	UpdatedTime *time.Time `json:"updated_time"` // 修改時間
	Sign        *string    `json:"sign"`         // 簽名
	Remarks     *string    `json:"remarks"`      // 備註
}

type OrderScanPayFirstForUpdateQuery struct {
	ID       *int64           `json:"id"`
	MemberID *int64           `json:"member_id"`
	Amount   *decimal.Decimal `json:"amount"`
}

type OrderScanPayCreateCond struct {
	MemberID        int64
	Amount          decimal.Decimal
	Brand           *string
	MerchantOrderID *string
	SourceOrderID   *string
	Content         *string
	RecordID        *int64
}

type OrderScanPayToPayCond struct {
	ID       int64            `json:"id"`        // 會員ID
	MemberID *int64           `json:"member_id"` // 會員ID
	Amount   *decimal.Decimal `json:"amount"`    // 金額
}
