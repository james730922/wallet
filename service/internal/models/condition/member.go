package condition

import (
	"time"

	"github.com/james730922/wallet/service/internal/models/model"
)

// MemberUpdateCond is restricted to state maintained by member login flows.
type MemberUpdateCond struct {
	ID                 *int64
	Status             *model.MemberStatus
	LastLoginTime      *time.Time
	FailedAttemptCount *int
}

type MemberUpdate struct {
	ID                 *int64              `json:"id"`
	Status             *model.MemberStatus `json:"status"`
	LastLoginTime      *time.Time          `json:"last_login_time"`
	FailedAttemptCount *int                `json:"failed_attempt_count"`
	UpdatedTime        *time.Time          `json:"updated_time"`
}

// MemberMappingQuery is used internally by authentication and retained payment flows.
type MemberMappingQuery struct {
	ID            *int64  `json:"id"`
	CountryCode   *string `json:"country_code"`
	Mobile        *string `json:"mobile"`
	QQ            *string `json:"qq"`
	WalletAddress *string `json:"wallet_address"`
}
