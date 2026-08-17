package model

import (
	"github.com/shopspring/decimal"

	"github.com/james730922/wallet/service/internal/thirdparty/event"
)

type WalletWithdrawNotify struct {
	Topic     event.Topic
	MemberIDs []int64
	Amount    decimal.Decimal
}

type WalletWithdraw struct {
	Title   string
	Content string
}
