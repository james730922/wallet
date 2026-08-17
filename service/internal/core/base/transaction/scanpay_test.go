package transaction

import "testing"

func TestNewScanPayCommonInitializesPaymentHandler(t *testing.T) {
	uc := newScanPayCommon()
	if uc.scanPayOrderAdd == nil {
		t.Fatal("newScanPayCommon() left scanPayOrderAdd nil")
	}
}
