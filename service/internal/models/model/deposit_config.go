package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type DepositConfigStatus int

const (
	DepositConfigStatusClosed DepositConfigStatus = iota
	DepositConfigStatusEnabled
	DepositConfigStatusDeleted
)

type DepositConfig struct {
	Bonus       decimal.Decimal `json:"bonus"`        // 充值紅利
	AdminID     *int64          `json:"admin_id"`     // 操作者
	AddedTime   time.Time       `json:"added_time"`   // 新增時間
	UpdatedTime time.Time       `json:"updated_time"` // 修改時間
}

func (d DepositConfigStatus) StatusName() string {
	names := [...]string{
		"凍結",
		"启用",
		"刪除",
	}
	return names[d]
}

func (*DepositConfig) TableName() string {
	return "deposit_config"
}

type DepositConfigMemberLevel struct {
	MemberLevel int64               `json:"member_level,string"`
	Bonus       *decimal.Decimal    `json:"bonus"`
	Status      DepositConfigStatus `json:"status,string"`
	AdminID     *int64              `json:"admin_id"`
	Remarks     *string             `json:"remarks"`
	AddedTime   time.Time           `json:"added_time"`
	UpdatedTime time.Time           `json:"updated_time"`
}

func (DepositConfigMemberLevel) TableName() string {
	return "deposit_config_member_level"
}

type DepositConfigMember struct {
	MemberID    int64               `json:"member_id,string"`
	Bonus       *decimal.Decimal    `json:"bonus"`
	Status      DepositConfigStatus `json:"status,string"`
	AdminID     *int64              `json:"admin_id"`
	Remarks     *string             `json:"remarks"`
	AddedTime   time.Time           `json:"added_time"`
	UpdatedTime time.Time           `json:"updated_time"`
}

func (DepositConfigMember) TableName() string {
	return "deposit_config_member"
}

type DepositConfigForMember struct {
	MemberID int64           `json:"member_id,string"`
	Bonus    decimal.Decimal `json:"bonus"`
}

type DepositConfigForMemberLevel struct {
	MemberLevel int64           `json:"member_level,string"`
	Bonus       decimal.Decimal `json:"bonus"`
}
