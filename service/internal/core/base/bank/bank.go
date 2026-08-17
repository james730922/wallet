package bank

import (
	"sync"

	"github.com/jinzhu/gorm"
	"go.uber.org/dig"
)

var (
	packet bankSet
)

var (
	once sync.Once
	self *bank
)

func NewBank(set bankSet) bank {
	once.Do(func() {
		packet = set

		newDAO()

		self = &bank{
			CategoryCommon: newCategoryCommon(),
			Common:         newBankCommon(),
		}
	})

	return *self
}

type bankSet struct {
	dig.In

	DB *gorm.DB
}

type bank struct {
	dig.Out

	CategoryCommon ICategoryCommon
	Common         IBankCommon
}
