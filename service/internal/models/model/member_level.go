package model

import (
	"time"
)

const MemberLevelDefault int = 1

type MemberLevelStatus int

const (
	MemberLevelStatusDisable MemberLevelStatus = iota
	MemberLevelStatusEnabled
)

type MemberLevelVisible int

const (
	MemberLevelVisibleDisable MemberLevelVisible = iota
	MemberLevelVisibleEnabled
)

type MemberLevelFeature int

const (
	MemberLevelFeatureNormal MemberLevelFeature = iota
	MemberLevelFeatureBlackList
)

type MemberLevel struct {
	ID          int64              // 會員級別代碼
	Name        string             // 會員級別名稱
	Status      MemberLevelStatus  // 會員狀態
	MemberCount int64              // 會員分類總數
	Sort        int                // 排序
	Note        string             // 備註
	Default     int                // 默認選項
	AdminID     *int64             // 操作者
	Visible     MemberLevelVisible // 顯示於列表
	Feature     MemberLevelFeature // 特徵
	AddedTime   time.Time          // 創建時間
	UpdatedTime time.Time          // 修改時間
}

func (MemberLevel) TableName() string {
	return "member_level"
}

func (s MemberLevelStatus) Name() string {
	names := [...]string{
		"禁用",
		"启用",
	}
	return names[s]
}

func (s MemberLevelVisible) Name() string {
	names := [...]string{
		"已删除",
		"未删除",
	}
	return names[s]
}

func (s MemberLevelFeature) Name() string {
	names := [...]string{
		"正常",
		"黑名单",
	}
	return names[s]
}
