package model

import (
	"time"
)

type DepositCategoryStatus int

const (
	DepositCategoryStatusDisable DepositCategoryStatus = 0
	DepositCategoryStatusEnable  DepositCategoryStatus = 1
)

func (d DepositCategoryStatus) StatusName() string {
	names := [...]string{
		"禁用",
		"启用",
	}
	return names[d]
}

type DepositCategory struct {
	ID          int64                 `sql:"id" form:"id" json:"id"`                                     // 入款分類代碼
	Name        string                `sql:"name" form:"name" json:"name"`                               // 分類名
	Type        int                   `sql:"type" form:"type" json:"type"`                               // 入款方法標記
	Image       string                `sql:"image" form:"image" json:"image"`                            // 圖片路徑
	Status      DepositCategoryStatus `sql:"status" form:"status" json:"status,string"`                  // 啟用狀態0=禁用1=啟用
	Sort        int                   `sql:"sort" form:"sort" json:"sort"`                               // 排序
	AddedTime   time.Time             `sql:"added_time" form:"added_time" json:"added_time,int64"`       // 創建時間
	UpdatedTime time.Time             `sql:"updated_time" form:"updated_time" json:"updated_time,int64"` // 修改時間
}

func (DepositCategory) TableName() string {
	return "bank_deposit_category"
}

type BankDepositCategoryItemStatus = int

const (
	BankDepositCategoryItemStatusDisable BankDepositCategoryItemStatus = 0
	BankDepositCategoryItemStatusEnable  BankDepositCategoryItemStatus = 1
)

type BankDepositCategoryItem struct {
	ID          int64                         `sql:"id"`
	CategoryID  int64                         `sql:"category_id"`
	AccountID   int64                         `sql:"account_id"`
	Status      BankDepositCategoryItemStatus `sql:"status"`
	Sort        int                           `sql:"sort"`
	AddedTime   time.Time                     `sql:"added_time"`
	UpdatedTime time.Time                     `sql:"updated_time"`
}

func (BankDepositCategoryItem) TableName() string {
	return "bank_deposit_category_item"
}

type BankDepositCategoryType struct {
	ID     int    `sql:"id"`
	Name   string `sql:"name"`
	Method string `sql:"method"` // 轉入方法
}

func (BankDepositCategoryType) TableName() string {
	return "bank_deposit_category_type"
}
