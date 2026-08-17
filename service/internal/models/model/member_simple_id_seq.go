package model

type MemberSimpleIDSeq struct {
	ID           uint8 `gorm:"column:id"`
	CurrentValue int64 `gorm:"column:current_value"`
}

func (m *MemberSimpleIDSeq) TableName() string {
	return "member_simple_id_seq"
}
