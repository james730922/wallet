package condition

import "github.com/james730922/wallet/service/internal/models"

type BankAccountQuery struct {
	models.Paging
	ID         *int64   `json:"id"`
	IDs        *[]int64 `json:"in_id"`
	Type       *int     `json:"type"`
	CategoryId *int64   `json:"category_id"`
	Status     *int     `json:"status"`
	Number     *string  `json:"number"`
	Visible    *int     `json:"visible"`
	Levels     *string  `json:"levels"`
}

type BankDepositCategoryItemViewCond struct {
	CategoryStatus *int `json:"category_status"` // item狀態
	AccountStatus  *int `json:"account_status"`  // 帳戶狀態
	AccountVisible *int `json:"account_visible"` // 帳戶是否刪除
	MemberLevel    *int `json:"member_level"`    // 會員級別
}
