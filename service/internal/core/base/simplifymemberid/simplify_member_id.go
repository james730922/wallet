package simplifymemberid

import (
	"sync"

	"github.com/jinzhu/gorm"
	"go.uber.org/dig"
)

var (
	packet simplifyIDSet
	once   sync.Once
	self   *simplifyID
)

func NewSimplifyMemberID(set simplifyIDSet) simplifyID {
	once.Do(func() {
		packet = set
		newDAO()

		self = &simplifyID{
			SimplifyMemberID: newSimplifyID(),
		}
	})

	return *self
}

type simplifyIDSet struct {
	dig.In

	DB      *gorm.DB
	DBSlave *gorm.DB `name:"dbSlave"`
}

type simplifyID struct {
	dig.Out

	SimplifyMemberID ISimplifyMemberID
}
