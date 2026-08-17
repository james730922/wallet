package transaction

import (
	"sync"

	"github.com/jinzhu/gorm"

	"github.com/bwmarrin/snowflake"
	"go.uber.org/dig"

	"github.com/james730922/wallet/service/internal/core/base/wallet"
)

var (
	packet transactionSet
	once   sync.Once
	self   *transaction
)

func NewTransaction(set transactionSet) transaction {
	once.Do(func() {
		packet = set

		newDAO()

		self = &transaction{
			Common:  newCommon(),
			ScanPay: newScanPayCommon(),
		}
	})

	return *self
}

type transactionSet struct {
	dig.In

	DB            *gorm.DB
	Node          *snowflake.Node
	Wallet        wallet.ICommon
	WalletGeneral wallet.IGeneral
}

type transaction struct {
	dig.Out

	Common  ICommon
	ScanPay IScanPayCommon
}
