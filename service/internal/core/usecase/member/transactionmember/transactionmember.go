package transactionmember

import (
	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/core/base/simplifymemberid"
	"github.com/james730922/wallet/service/internal/core/base/transaction"
	"go.uber.org/dig"
	"sync"
)

var (
	packet transactionMemberSet
	once   sync.Once
	ptr    *transactionMember
)

func NewTransactionMember(set transactionMemberSet) transactionMember {
	once.Do(func() {
		packet = set

		ptr = &transactionMember{
			Member: newMember(),
		}
	})

	return *ptr
}

type transactionMemberSet struct {
	dig.In

	DB               *gorm.DB
	Transaction      transaction.ICommon
	SimplifyMemberID simplifymemberid.ISimplifyMemberID
}

type transactionMember struct {
	dig.Out

	Member IMember
}
