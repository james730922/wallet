package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransactionSourceType int

// Values are explicit so removing retired source types does not renumber persisted records.
const (
	TransactionSourceTypeDeposit        TransactionSourceType = 1
	TransactionSourceTypeDepositBonus   TransactionSourceType = 6
	TransactionSourceTypeScanPayConfirm TransactionSourceType = 20
)

func GetAllTransactionSourceType() []TransactionSourceType {
	return []TransactionSourceType{
		TransactionSourceTypeDeposit,
		TransactionSourceTypeDepositBonus,
		TransactionSourceTypeScanPayConfirm,
	}
}

func (s TransactionSourceType) Name() string {
	switch s {
	case TransactionSourceTypeDeposit:
		return "充值"
	case TransactionSourceTypeDepositBonus:
		return "充值红利"
	case TransactionSourceTypeScanPayConfirm:
		return "扫码支付"
	default:
		return ""
	}
}

func (s TransactionSourceType) ClassName() string {
	switch s {
	case TransactionSourceTypeDepositBonus, TransactionSourceTypeDeposit:
		return "充值"
	case TransactionSourceTypeScanPayConfirm:
		return "扫码支付"
	default:
		return ""
	}
}

type Transaction struct {
	ID                    int64                 `gorm:"column:id"`
	MemberID              int64                 `gorm:"column:member_id"`
	SourceType            TransactionSourceType `gorm:"column:source_type"`
	SourceID              int64                 `gorm:"column:source_id"`
	CurrencyCode          string                `gorm:"column:currency_code"`
	Amount                decimal.Decimal       `gorm:"column:amount"`
	CurrentTotalAmount    decimal.Decimal       `gorm:"column:current_total_amount"`
	ChangedTotalAmount    decimal.Decimal       `gorm:"column:changed_total_amount"`
	CurrentBalance        decimal.Decimal       `gorm:"column:current_balance"`
	ChangedBalance        decimal.Decimal       `gorm:"column:changed_balance"`
	AddedTime             time.Time             `gorm:"column:added_time"`
	UpdatedTime           time.Time             `gorm:"column:updated_time"`
	Sign                  string                `gorm:"column:sign"`
	Remarks               string                `gorm:"column:remarks"`
	Merchant              string                `gorm:"column:merchant"`
	MerchantMemberAccount string                `gorm:"column:merchant_member_account"`
}

func (Transaction) TableName() string { return "transaction" }
