package scanpay

import (
	"errors"
	"testing"

	"github.com/shopspring/decimal"

	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func TestReusableOrderIDReturnsCanonicalTransactionOrder(t *testing.T) {
	recordID := int64(11)
	cond := &condition.OrderScanPayCreateCond{
		MemberID: 7,
		Amount:   decimal.NewFromInt(100),
		RecordID: &recordID,
	}
	order := &model.OrderScanPay{
		ID:       101,
		MemberID: 7,
		Amount:   decimal.NewFromInt(100),
		Status:   model.OrderScanPayStatusEnumTransaction,
		RecordID: &recordID,
	}

	id, err := reusableOrderID(order, cond)
	if err != nil || id != order.ID {
		t.Fatalf("reusableOrderID() = (%d, %v), want (%d, nil)", id, err, order.ID)
	}
}

func TestReusableOrderIDRejectsDifferentMember(t *testing.T) {
	recordID := int64(11)
	cond := &condition.OrderScanPayCreateCond{
		MemberID: 8,
		Amount:   decimal.NewFromInt(100),
		RecordID: &recordID,
	}
	order := &model.OrderScanPay{
		ID:       101,
		MemberID: 7,
		Amount:   decimal.NewFromInt(100),
		Status:   model.OrderScanPayStatusEnumTransaction,
		RecordID: &recordID,
	}

	if _, err := reusableOrderID(order, cond); !errors.Is(err, errs.ScanPayAddRecordFailed) {
		t.Fatalf("cross-member reusableOrderID() error = %v", err)
	}
}

func TestReusableOrderIDRejectsFailureOrder(t *testing.T) {
	recordID := int64(11)
	cond := &condition.OrderScanPayCreateCond{
		MemberID: 7,
		Amount:   decimal.NewFromInt(100),
		RecordID: &recordID,
	}
	order := &model.OrderScanPay{
		ID:       101,
		MemberID: 7,
		Amount:   decimal.NewFromInt(100),
		Status:   model.OrderScanPayStatusEnumFailure,
		RecordID: &recordID,
	}

	if _, err := reusableOrderID(order, cond); !errors.Is(err, errs.ScanPayOrderFailure) {
		t.Fatalf("failure order reusableOrderID() error = %v, want ScanPayOrderFailure", err)
	}
}

func TestTerminalPaymentFailureClassification(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "insufficient balance is final", err: errs.WalletMemberUpdateBalanceIsNegative, want: true},
		{name: "invalid decimal is final", err: errs.CommonAmountDecimalPlacesError, want: true},
		{name: "wallet invariant needs manual review", err: errs.WalletMemberAmountUnreasonable, want: false},
		{name: "database error is uncertain", err: errs.DBOperationFailed, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTerminalPaymentFailure(tt.err); got != tt.want {
				t.Fatalf("isTerminalPaymentFailure(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
