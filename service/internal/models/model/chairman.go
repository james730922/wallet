package model

import "time"

type Chairman struct {
	ID        int64     // 識別碼
	AddedTime time.Time // 新增時間
}

func (*Chairman) TableName() string {
	return "chairman"
}
