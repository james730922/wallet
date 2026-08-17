package resp

import "github.com/james730922/wallet/service/internal/models/model"

type DepositCategoryGet struct {
	// 入款分類代碼
	ID int64 `json:"id,string"`
	// 分類名
	Name string `json:"name"`
	// 入款方法標記
	Type     int    `json:"type,string"`
	TypeName string `json:"type_name"`
	// 圖片路徑
	Image string `json:"image"`
	// 啟用狀態
	Status model.DepositCategoryStatus `json:"status,string"`
	// 排序
	Sort int `json:"sort"`
}

type DepositCategoryOptions struct {
	// 入款分類代碼
	Id int64 `json:"id,string"`
	// 分類名
	Name string `json:"name"`
	// 入款類型
	Type int `json:"type,string"`
}
