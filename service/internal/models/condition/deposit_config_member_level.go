package condition

import (
	"github.com/james730922/wallet/service/internal/models"
	"time"
)

type DepositConfigMemberLevelQuery struct {
	models.Paging
	MemberLevels     *[]int64   `json:"in_member_level"` //會員層級
	Status           *int       `json:"status"`
	Statuses         *[]int     `json:"in_status"`
	StartAtAddedTime *time.Time `json:"start_at_added_time"`
	EndAtAddedTime   *time.Time `json:"end_at_added_time"`
}
