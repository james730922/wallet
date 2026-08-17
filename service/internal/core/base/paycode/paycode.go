package paycode

import (
	"go.uber.org/dig"
)

func NewPayment() payCodeOut {
	return payCodeOut{
		SenPay: newScanPay(),
	}
}

type payCode struct {
}

type payCodeOut struct {
	dig.Out

	SenPay IScanPay
}
