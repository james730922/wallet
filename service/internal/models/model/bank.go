package model

import (
	"time"
)

type BankCodeStatus int

const (
	BankCodeStatusDisable BankCodeStatus = 0
	BankCodeStatusEnable  BankCodeStatus = 1
)

func (d BankCodeStatus) StatusName() string {
	names := [...]string{
		"禁用",
		"启用",
	}
	return names[d]
}

type BankCode struct {
	Code        string         `gorm:"code" form:"code" json:"code"`
	Name        string         `gorm:"name" form:"name" json:"name"`
	Image       string         `gorm:"image" form:"image" json:"image"`
	Status      BankCodeStatus `gorm:"status" form:"status" json:"status,string"`
	AddedTime   time.Time      `gorm:"added_time" form:"added_time" json:"added_time"`
	UpdatedTime time.Time      `gorm:"updated_time" form:"updated_time" json:"updated_time"`
}

func (BankCode) TableName() string {
	return "bank_code"
}
