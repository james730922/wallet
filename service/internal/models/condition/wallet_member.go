package condition

import (
	"time"

	"github.com/shopspring/decimal"
)

type WalletMemberUpdate struct {
	MemberID     *int64           `json:"member_id"`     // '會員id'
	Balance      *decimal.Decimal `json:"balance"`       // '錢包餘額'
	TotalAmount  *decimal.Decimal `json:"total_amount"`  // '錢包餘額'
	Amount       *decimal.Decimal `json:"amount"`        // '入款餘額'
	Bonus        *decimal.Decimal `json:"bonus"`         // '紅利餘額'
	FrozenAmount *decimal.Decimal `json:"frozen_amount"` // '凍結額'
	Sign         *string          `json:"sign"`          // '簽名'
	UpdatedTime  *time.Time       `json:"updated_time"`  // '更新時間'
}

type WalletMemberUpdateCond struct {
	MemberID int64           // '會員id'
	Amount   decimal.Decimal // '加減的金額'
}
