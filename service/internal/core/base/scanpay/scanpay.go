package scanpay

import (
	"github.com/bwmarrin/snowflake"
	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/core/base/paycode"
	"github.com/james730922/wallet/service/internal/core/base/transaction"
	"go.uber.org/dig"
	"sync"
)

var (
	once   sync.Once
	packet scanPaySet
	self   *scanPay
)

const (
	journalContentAdd     = "新增扫码交易"
	journalContentSuccess = "扫码交易完成"
	journalContentFailure = "扫码交易失败"
)

func NewScanPay(set scanPaySet) scanPay {
	once.Do(func() {
		packet = set

		newDAO()

		self = &scanPay{
			ScanPayCommon: newScanPayCommon(),
			Record:        newRecordCommon(),
			Reconciler:    newScanPayReconciler(),
		}
	})

	return *self
}

type scanPaySet struct {
	dig.In

	DB          *gorm.DB
	Node        *snowflake.Node
	PayCode     paycode.IScanPay
	Transaction transaction.IScanPayCommon
}

type scanPay struct {
	dig.Out
	ScanPayCommon IScanPayCommon
	Record        IRecordCommon
	Reconciler    IReconciler
}
