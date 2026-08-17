package condition

type DepositCategoryQuery struct {
	ID   *int64  `json:"id"`   // 入款分類代碼
	Name *string `json:"name"` // 分類名
}
