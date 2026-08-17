package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// JournalMemberWallet is the retained financial audit trail for wallet mutations.
// It is intentionally not exposed through the removed log-management APIs.
type JournalMemberWallet struct {
	ID           int64           `gorm:"column:id"`
	MemberID     int64           `gorm:"column:member_id"`
	Balance      decimal.Decimal `gorm:"column:balance"`
	TotalAmount  decimal.Decimal `gorm:"column:total_amount"`
	Amount       decimal.Decimal `gorm:"column:amount"`
	Bonus        decimal.Decimal `gorm:"column:bonus"`
	FrozenAmount decimal.Decimal `gorm:"column:frozen_amount"`
	WalletSign   string          `gorm:"column:wallet_sign"`
	AddedTime    time.Time       `gorm:"column:added_time"`
}

func (JournalMemberWallet) TableName() string {
	return "wallet_journal"
}
