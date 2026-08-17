package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type OrderBonusCalculateTypeEnum = int

const (
	OrderBonusCalculateTypeEnumCash OrderBonusCalculateTypeEnum = iota
	OrderBonusCalculateTypeEnumRate
)

type OrderBonusSourceTypeEnum = int

const (
	OrderBonusSourceTypeEnumNormal OrderBonusSourceTypeEnum = iota
	OrderBonusSourceTypeEnumDeposit
)

type OrderBonusStatusEnum = int

const (
	OrderBonusStatusEnumWaiting OrderBonusStatusEnum = iota
	OrderBonusStatusEnumAccept
	OrderBonusStatusEnumReject
)

type OrderBonus struct {
	ID            int64
	MemberID      int64
	CurrencyCode  string
	Amount        decimal.Decimal
	CalculateType OrderBonusCalculateTypeEnum
	SourceType    OrderBonusSourceTypeEnum
	SourceOrderID int64
	SourceAmount  decimal.Decimal
	SourceRate    decimal.Decimal
	Status        OrderBonusStatusEnum
	AcceptTime    *time.Time
	AcceptAdminID *int64
	AddedTime     time.Time
	UpdatedTime   time.Time
	Sign          string
}

func (OrderBonus) TableName() string {
	return "order_bonus"
}
