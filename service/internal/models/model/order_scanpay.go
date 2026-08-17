package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type OrderScanPayStatusEnum int

const (
	OrderScanPayStatusEnumTransaction OrderScanPayStatusEnum = iota // 交易中
	OrderScanPayStatusEnumSuccess                                   // 成功
	OrderScanPayStatusEnumFailure                                   // 失敗：永久作廢，不得重試或重用
	OrderScanPayStatusEnumCancel                                    // 取消
)

func (s OrderScanPayStatusEnum) Name() string {
	names := [...]string{
		"交易中",
		"完成",
		"失败",
		"取消",
	}

	return names[s]
}

type OrderScanPay struct {
	ID              int64                  // 掃碼訂單號
	MemberID        int64                  // 會員ID
	Amount          decimal.Decimal        // 金額
	Status          OrderScanPayStatusEnum // 狀態：0=交易中 1=成功 2=失敗作廢 3=取消
	RecordID        *int64                 // 掃碼申請紀錄
	Brand           string                 // 品牌
	MerchantOrderID string                 // 收銀台單號
	SourceOrderID   string                 // 串接方單號
	Content         string                 // 支付二維碼內容
	SuccessTime     *time.Time             // 成功時間
	AddedTime       time.Time              // 創建時間
	CancelTime      *time.Time             // 取消時間
	UpdatedTime     time.Time              // 修改時間
	Sign            string                 // 簽名
	AdminID         *int64                 // 操作者
	Remarks         *string                // 註解
}

func (OrderScanPay) TableName() string {
	return "order_scanpay"
}

// 統計
type OrderScanPaySummary struct {
	Status  int             `gorm:"column:status"`  // 狀態
	Numbers int64           `gorm:"column:numbers"` // 筆數
	Amount  decimal.Decimal `gorm:"column:amount"`  // 金額
}

type OrderScanPayMapping struct {
	ID                 int64                  // 掃碼訂單號
	MemberID           int64                  // 會員ID
	Amount             decimal.Decimal        // 金額
	Status             OrderScanPayStatusEnum // 狀態：0=失敗，1=成功
	RecordID           *int64                 // 掃碼申請紀錄
	Brand              string                 // 品牌
	MerchantOrderID    string                 // 收銀台單號
	MerchantMemberName *string                // 串接方會員姓名
	SourceOrderID      string                 // 串接方單號
	Content            string                 // 支付二維碼內容
	SuccessTime        *time.Time             // 成功時間
	AddedTime          time.Time              // 創建時間
	CancelTime         *time.Time             // 取消時間
	UpdatedTime        time.Time              // 修改時間
	Sign               string                 // 簽名
	AdminID            *int64                 // 操作者
	Remarks            *string                // 註解
}

func (OrderScanPayMapping) TableName() string {
	return "order_scanpay"
}
