package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type MemberCreditHour struct {
	Hour          string          // 小時區間 (ex:2020010101)
	MemberID      int64           // 會員ID
	IntroducerID  int64           // 上線ID
	ChildrenCount int             //直屬人數
	DepositCount  int             // 充值次數
	DepositAmount decimal.Decimal // 充值總額
	DepositComm   decimal.Decimal // 充值佣金
	AddedTime     time.Time       // 新增時間
	UpdatedTime   time.Time       // 更新時間
}

func (MemberCreditHour) TableName() string {
	return "member_credit_hour"
}

type MemberCreditFromChildHour struct {
	Hour          string          // 小時區間 (ex:2020010101)
	MemberID      int64           // 會員ID
	IntroducerID  int64           // 上線ID
	DepositCount  int             // 充值次數
	DepositAmount decimal.Decimal // 充值總額
	DepositComm   decimal.Decimal // 充值佣金
	AddedTime     time.Time       // 新增時間
	UpdatedTime   time.Time       // 更新時間
}

func (MemberCreditFromChildHour) TableName() string {
	return "member_credit_from_child_hour"
}
