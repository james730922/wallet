package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type BankAccountStatus int

const (
	BankAccountStatusDisable BankAccountStatus = 0
	BankAccountStatusEnable  BankAccountStatus = 1
)

func (b BankAccountStatus) StatusName() string {
	names := [...]string{
		"禁用",
		"启用",
	}
	return names[b]
}

type BankAccountType = int

const (
	BankAccountTypeBank   BankAccountType = 1
	BankAccountTypeQRCode BankAccountType = 2
)

type BankAccountLevelsEnum int

const (
	BankAccountAllLevels BankAccountLevelsEnum = iota
)

func (l BankAccountLevelsEnum) Name() string {
	names := [...]string{
		"全部会员级别",
	}

	return names[l]
}

func (l BankAccountLevelsEnum) Int() int {
	return int(l)
}

type BankAccountMethodType int

const (
	BankAccountDeposit BankAccountMethodType = iota
)

func (b BankAccountMethodType) StatusName() string {
	names := [...]string{
		"人工充值",
	}
	return names[b]
}

type BankAccount struct {
	ID           int64                 `gorm:"column:id"`            // 序號
	Type         BankAccountMethodType `gorm:"column:type"`          // 使用類別
	Number       string                `gorm:"column:number"`        // 卡號
	BankCode     string                `gorm:"column:bank_code"`     // 銀行碼
	BankBranch   string                `gorm:"column:bank_branch"`   // 銀行支行
	Name         string                `gorm:"column:name"`          // 戶名
	CurrencyCode string                `gorm:"column:currency_code"` // 支付幣別
	Levels       string                `gorm:"column:levels"`        // 會員級別
	ReceiveLimit decimal.Decimal       `gorm:"column:receive_limit"` // 收款上限0=無上限
	Status       BankAccountStatus     `gorm:"column:status"`        // 啟用狀態0=禁用1=啟用
	Visible      int                   `gorm:"visible"`              // 軟刪除與否
	MinAmount    decimal.Decimal       `gorm:"column:min_amount"`    // 單次入款最低限額
	MaxAmount    decimal.Decimal       `gorm:"column:max_amount"`    // 單次入款最高限額
	AdminID      *int64                `gorm:"column:admin_id"`      // 操作者
	Remark       string                `gorm:"column:remark"`        // 備註
	QRCode       string                `gorm:"column:qrcode"`        // 面對面QRCode
	AddedTime    time.Time             `gorm:"column:added_time"`    // 創建時間
	UpdatedTime  time.Time             `gorm:"column:updated_time"`  // 修改時間
}

func (BankAccount) TableName() string {
	return "bank_account"
}

type BankDepositCategoryItemView struct {
	CategoryID         int64           `sql:"column:id"`                   // item 序號
	CategoryName       string          `sql:"column:category_name"`        // item 名稱
	CategoryStatus     int             `sql:"column:category_status"`      // item狀態
	Sort               int32           `sql:"column:sort"`                 // item 排序
	AccountID          int64           `sql:"column:account_id"`           // 銀行ID
	Number             string          `sql:"column:number"`               // 卡號
	CurrencyCode       string          `sql:"column:currency_code"`        // 支付幣別
	Levels             string          `sql:"column:levels"`               // 會員級別
	Name               string          `sql:"column:name"`                 // 戶名
	BankCode           string          `sql:"column:bank_code"`            // 銀行碼
	BankImage          string          `sql:"column:bank_image"`           // 銀行圖片
	BankName           string          `sql:"column:bank_name"`            // 銀行名稱
	BankBranch         string          `sql:"column:bank_branch"`          // 銀行支行
	MinAmount          decimal.Decimal `sql:"column:min_amount"`           // 單次入款最低限額
	MaxAmount          decimal.Decimal `sql:"column:max_amount"`           // 單次入款最高限額
	Type               int32           `sql:"column:type"`                 // 帳號類型：0=銀行，1=微信面對面，2=支付寶面對面,3=微信轉帳，4=支付寶轉帳
	QRCode             string          `sql:"column:qrcode"`               // 面對面QRCode
	AccountStatus      int             `sql:"column:account_status"`       // 帳戶狀態
	AccountVisible     int             `sql:"column:account_visible"`      // 帳戶是否已被刪除
	AccountUpdatedTime time.Time       `sql:"column:account_updated_time"` // 帳戶更新時間
}
