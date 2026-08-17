package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type ScanPayRecordStatus = int

const (
	ScanPayRecordStatusWaiting     ScanPayRecordStatus = 0
	ScanPayRecordStatusTransaction ScanPayRecordStatus = 1
	ScanPayRecordStatusDone        ScanPayRecordStatus = 2
	ScanPayRecordStatusCancel      ScanPayRecordStatus = 3
	ScanPayRecordStatusFailure     ScanPayRecordStatus = 4
)

func ScanPayRecordStatusName(s ScanPayRecordStatus) string {
	names := map[ScanPayRecordStatus]string{
		ScanPayRecordStatusWaiting:     "待确认",
		ScanPayRecordStatusTransaction: "交易中",
		ScanPayRecordStatusDone:        "完成",
		ScanPayRecordStatusCancel:      "已取消",
		ScanPayRecordStatusFailure:     "失败",
	}

	return names[s]
}

type ScanPayRecord struct {
	ID            int64               // 期數
	Amount        decimal.Decimal     // 金額
	Status        ScanPayRecordStatus // 狀態 0=待確認 1=交易中 2=已確認 3=已取消 4=付款失敗且作廢
	Brand         string              // 品牌
	MerchantID    string              // 收銀台單號
	SourceOrderID string              // 串接方單號
	Content       string              // 支付二維碼內容
	ExpiredTime   time.Time           // 二維碼過期時間
	AddedTime     time.Time           // 創建時間
	CancelTime    *time.Time          // 取消時間
	UpdatedTime   time.Time           // 修改時間
	AdminID       *int64              // 操作者
	Remarks       *string             // 註解
}

func (ScanPayRecord) TableName() string {
	return "scanpay_record"
}
