package condition

type MemberLevelQuery struct {
	ID      *int64  `json:"id"`      // 會員級別代碼
	Name    *string `json:"name"`    // 會員級別名稱
	Default *int    `json:"default"` // 默認選項
	Status  *int    `json:"status"`  // 啟用禁用
	Visible *int    `json:"visible"` // 顯示於列表
}

type MemberLevelUpdate struct {
	ID          *int64 `json:"id"`
	MemberCount *int64 `json:"member_count"`
}
