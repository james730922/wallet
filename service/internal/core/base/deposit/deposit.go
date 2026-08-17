package deposit

import (
	"sync"

	"github.com/go-redis/redis/v7"

	"github.com/bwmarrin/snowflake"
	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/core/base/bank"
	"github.com/james730922/wallet/service/internal/core/base/user"
	"github.com/james730922/wallet/service/internal/thirdparty/fileserver"
	"go.uber.org/dig"
)

var (
	packet depositSet
	once   sync.Once
	self   *deposit
)

func NewDeposit(set depositSet) depositOut {
	once.Do(func() {
		packet = set

		newDAO()

		self = &deposit{
			DepositCommon: newDepositCommon(),
			DepositMember: newDepositMember(),
			DepositBonus:  newBonusCommon(),
			DepositConfig: newDepositConfig(),
			DepositCache:  newDepositCache(),
		}
	})

	return depositOut{
		DepositCommon: self.DepositCommon,
		DepositMember: self.DepositMember,
		DepositBonus:  self.DepositBonus,
		DepositConfig: self.DepositConfig,
	}
}

type depositSet struct {
	dig.In

	Node           *snowflake.Node
	DB             *gorm.DB
	Redis          *redis.Client
	Member         user.IMemberCommon
	MemberLevel    user.IMemberLevelCommon
	FileServer     fileserver.IZQBFileServer
	CategoryCommon bank.ICategoryCommon
	Bank           bank.IBankCommon
}

type depositOut struct {
	dig.Out

	DepositCommon IDepositCommon
	DepositMember IDepositMember
	DepositBonus  IBonusCommon
	DepositConfig IDepositConfig
}

type deposit struct {
	DepositCommon IDepositCommon
	DepositMember IDepositMember
	DepositBonus  IBonusCommon
	DepositConfig IDepositConfig
	DepositCache  iDepositCache
}
