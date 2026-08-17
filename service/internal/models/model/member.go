package model

import "time"

type MemberStatus int

const (
	MemberStatusFreezed   MemberStatus = 0
	MemberStatusEnabled   MemberStatus = 1
	MemberStatusSuspended MemberStatus = 2
)

func (j MemberStatus) Name() string {
	names := [...]string{
		"帐号冻结",
		"帐号启用",
		"帐号暂停",
	}

	return names[j]
}

type Member struct {
	ID                 int64        `gorm:"column:id"`                   // 會員ID
	LevelCode          int64        `gorm:"column:level_code"`           // 會員級別
	Status             MemberStatus `gorm:"column:status"`               // 用戶登入狀態
	LastLoginTime      *time.Time   `gorm:"column:last_login_time"`      // 最後登陸時間
	FailedAttemptCount int          `gorm:"column:failed_attempt_count"` // 嘗試登入錯誤次數
	Remark             *string      `gorm:"column:remark"`               // 備注
	AdminID            *int64       `gorm:"column:admin_id"`             // 操作者
	AddedTime          time.Time    `gorm:"column:added_time"`           // 創建時間
	UpdatedTime        time.Time    `gorm:"column:updated_time"`         // 修改時間
}

func (Member) TableName() string {
	return "member"
}

type MemberMapping struct {
	ID             int64     `gorm:"column:id"`                                             // 會員ID
	CountryCode    string    `gorm:"column:country_code"`                                   // 國碼
	Mobile         string    `gorm:"column:mobile"`                                         // 手機號碼
	Name           string    `gorm:"column:name" column:"member_mapping.name" field:"會員姓名"` // 中文姓名
	QQ             string    `gorm:"column:qq"`                                             // QQ號
	Passwd         string    `gorm:"column:passwd"`                                         // 密碼
	Salt           string    `gorm:"column:salt"`                                           // 密碼鹽
	SimplifyID     string    `gorm:"column:simplify_id"`                                    // 簡化會員ID
	SecurityPasswd string    `gorm:"column:security_passwd"`                                // 安全密码
	WalletAddress  string    `gorm:"column:wallet_address"`                                 // 錢包地址
	AddedTime      time.Time `gorm:"column:added_time"`
	UpdatedTime    time.Time `gorm:"column:updated_time"`
}

func (MemberMapping) TableName() string {
	return "member_mapping"
}

type MemberRepository struct {
	Token string
	ID    int64
}
