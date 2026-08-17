package condition

import (
	"time"

	"github.com/james730922/wallet/service/internal/models/model"
)

type OrderBonusQuery struct {
	ID                 *int64                          `json:"id"`
	MemberID           *int64                          `json:"member_id"`
	Status             *model.OrderBonusStatusEnum     `json:"status"`
	SourceType         *model.OrderBonusSourceTypeEnum `json:"source_type"`
	SourceOrderID      *int64                          `json:"source_order_id"`
	SourceOrderIDs     *[]int64                        `json:"in_source_order_id"`
	StartAtAddedTime   *time.Time                      `json:"start_at_added_time"`
	EndAtAddedTime     *time.Time                      `json:"end_at_added_time"`
	StartAtUpdatedTime *time.Time                      `json:"start_at_updated_time"`
	EndAtUpdatedTime   *time.Time                      `json:"end_at_updated_time"`
}
