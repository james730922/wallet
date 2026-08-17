package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type DepositStatus int

const DepositStatusAll = -1

const (
	DepositStatusWaiting DepositStatus = iota
	DepositStatusAccept
	DepositStatusReject
	DepositStatusPause
)

func GetAllDepositStatus() []DepositStatus {
	return []DepositStatus{
		DepositStatusWaiting,
		DepositStatusAccept,
		DepositStatusReject,
	}
}

func (s DepositStatus) Name() string {
	names := [...]string{
		"待审核",
		"已入款",
		"已拒绝",
		"暂停",
	}

	return names[s]
}

// 給予APP充值狀態顯示
func (s DepositStatus) MemberName() string {
	names := [...]string{
		"待确认",
		"成功",
		"取消",
		"暂停",
	}

	return names[s]
}

type Deposit struct {
	ID              int64           // 入款序號
	MemberID        int64           // 會員ID
	AccountID       int64           // 銀行帳號ID
	AccountNumber   string          // 卡號
	AccountBankCode string          // 銀行碼
	CurrencyCode    string          // 支付幣別
	PayName         string          // 入款人戶名
	Amount          decimal.Decimal // 金額
	Charge          decimal.Decimal // 收費
	Status          DepositStatus   // 狀態
	AcceptTime      *time.Time      // 確認時間
	CancelTime      *time.Time      // 取消時間
	AdminID         *int64          // 操作人員
	AddedTime       time.Time       // 創建時間
	UpdatedTime     time.Time       // 修改時間
	Remarks         string          // 註記
	Sign            string          // 簽名
}

func (Deposit) TableName() string {
	return "deposit"
}
