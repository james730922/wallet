package wallet

import (
	"sync"

	"github.com/bwmarrin/snowflake"
	"go.uber.org/dig"
)

var (
	packet walletSet
)

var (
	once sync.Once
	ptr  *wallet
)

func NewWallet(set walletSet) wallet {
	once.Do(func() {
		packet = set
		newDAO()

		ptr = &wallet{
			Common:  newCommon(),
			General: newGeneral(),
		}
	})

	return *ptr
}

type walletSet struct {
	dig.In

	Node *snowflake.Node
}

type wallet struct {
	dig.Out

	Common  ICommon
	General IGeneral
}
