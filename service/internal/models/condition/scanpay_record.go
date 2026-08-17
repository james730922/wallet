package condition

import (
	"github.com/james730922/wallet/service/internal/models/model"
	"time"
)

type ScanPayRecordUpdate struct {
	ID          *int64                     `json:"id"`           // 掃碼支付RecordID
	Status      *model.ScanPayRecordStatus `json:"status"`       // 狀態0=待確認 1=交易中 2=已確認 3=已取消 4=失敗作廢
	UpdatedTime *time.Time                 `json:"updated_time"` // 修改時間
	Remarks     *string                    `json:"remarks"`      // 狀態變更原因
}
