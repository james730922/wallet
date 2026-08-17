package condition

import "time"

type BankCodeQuery struct {
	Code     *string `form:"code"        json:"code"`            // 銀行碼
	BankName *string `form:"bank_name"   json:"bank_name"`       //銀行名稱
	Status   *int    `form:"status,string" json:"status,string"` // 啟用狀態
}

type BankCodeUpdate struct {
	Code       *string    `json:"code"`          // 銀行碼
	Image      *string    `json:"image"`         // 圖片
	Status     *int       `json:"status,string"` // 啟用狀態
	UpdateTime *time.Time `json:"updated_time"`  // 更新時間
}
