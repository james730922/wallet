package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type ScanPayMapping struct {
	RecordID              int64           // K幣掃碼申請紀錄ID
	Merchant              string          // 商戶
	MerchantOrderID       string          // 商戶單號(外部單號)
	MerchantMemberID      string          // 外部會員id
	MerchantMemberAccount string          // 外部會員帳號
	MerchantMemberName    *string         // 外部匯款會員名稱
	Amount                decimal.Decimal // 金額
	Remark                *string         // 註記
	AddedTime             time.Time       // 新增時間
}

func (ScanPayMapping) TableName() string {
	return "scanpay_mapping"
}
