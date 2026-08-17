package scanpaymember

import (
	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/core/base/paycode"
	"github.com/james730922/wallet/service/internal/core/base/scanpay"
	"github.com/james730922/wallet/service/internal/core/base/user"
	"go.uber.org/dig"
	"sync"
)

var (
	once   sync.Once
	packet scanPaySet
	self   *scanPayOut
)

func NewScanPay(set scanPaySet) scanPayOut {
	once.Do(func() {
		packet = set

		self = &scanPayOut{
			ScanPayMember: newScanPayMember(),
		}
	})
	return *self
}

type scanPaySet struct {
	dig.In
	DB            *gorm.DB
	ScanPayRecord scanpay.IRecordCommon
	ScanPayCommon scanpay.IScanPayCommon
	PayCode       paycode.IScanPay
	UserMember    user.IMemberCommon
	MemberLevel   user.IMemberLevelCommon
}

type scanPayOut struct {
	dig.Out
	ScanPayMember IMember
}
