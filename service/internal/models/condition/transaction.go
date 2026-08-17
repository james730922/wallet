package condition

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/james730922/wallet/service/internal/models"
	"github.com/james730922/wallet/service/internal/models/model"
)

type TransactionQuery struct {
	models.Paging
	MemberID         *int64                       `json:"transaction.member_id"`      // 會員ID
	SourceType       *model.TransactionSourceType `json:"transaction.source_type"`    // 來源類型：0=入款
	SourceTypes      *[]int                       `json:"in_transaction.source_type"` // 來源類型s：0=入款
	SourceID         *int64                       `json:"transaction.source_id"`      // 來源單號
	CountryCode      *string                      `json:"member_mapping.country_code"`
	Mobile           *string                      `json:"member_mapping.mobile"`
	MemberName       *string                      `json:"like_member_mapping.name"`
	StartAtAmount    *decimal.Decimal             `json:"start_at_transaction.amount"`
	EndAtAmount      *decimal.Decimal             `json:"end_at_transaction.amount"`
	StartAtAddedTime *time.Time                   `json:"start_at_transaction.added_time"`
	EndAtAddedTime   *time.Time                   `json:"end_at_transaction.added_time"`
}
