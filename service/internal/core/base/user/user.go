package user

import (
	"sync"

	"github.com/bwmarrin/snowflake"
	"github.com/jinzhu/gorm"
	"go.uber.org/dig"

	"github.com/james730922/wallet/service/internal/core/base/simplifymemberid"
	"github.com/james730922/wallet/service/internal/core/base/wallet"
)

var (
	packet userSet
)

var (
	once sync.Once
	self *user
)

func NewUser(set userSet) user {
	once.Do(func() {
		packet = set

		newDAO()
		self = &user{
			Member:      newMemberCommon(),
			MemberLevel: newMemberLevelCommon(),
		}
	})

	return *self
}

type userSet struct {
	dig.In

	DB               *gorm.DB
	Node             *snowflake.Node
	Wallet           wallet.IGeneral
	SimplifyMemberID simplifymemberid.ISimplifyMemberID
}

type user struct {
	dig.Out

	Member      IMemberCommon
	MemberLevel IMemberLevelCommon
}
