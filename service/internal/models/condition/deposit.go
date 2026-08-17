package condition

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/james730922/wallet/service/internal/models"
)

type DepositQuery struct {
	models.Paging
	ID                *int64           `json:"id" form:"id"`                                     // 入款ID
	MemberID          *int64           `json:"member_id" form:"member_id"`                       // 會員ID
	AccountNumber     *string          `form:"account_number" json:"account_number"`             // 卡號
	Amount            *decimal.Decimal `json:"amount" form:"amount"`                             // 金額
	StartAtAmount     *decimal.Decimal `json:"start_at_amount" form:"start_at_amount"`           // 金額起始值
	EndAtAmount       *decimal.Decimal `json:"end_at_amount" form:"end_at_amount"`               // 金額結束值
	Status            *int             `json:"status,string" form:"status"`                      // 狀態
	StartAtAddedTime  *time.Time       `json:"start_at_added_time" form:"start_at_added_time"`   // 申請時間起始值
	EndAtAddedTime    *time.Time       `json:"end_at_added_time" form:"end_at_added_time"`       // 申請時間結束值
	StartAtAcceptTime *time.Time       `json:"start_at_accept_time" form:"start_at_accept_time"` // 操作時間起始值
	EndAtAcceptTime   *time.Time       `json:"end_at_accept_time" form:"end_at_accept_time"`     // 操作時間結束值
}

type DepositListByMemberCond struct {
	models.Paging
	MemberID int64 // 會員ID
	Status   *int  // 狀態
}
