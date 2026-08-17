package model

import "testing"

func TestScanPayRecordStatusNameFailure(t *testing.T) {
	if got := ScanPayRecordStatusName(ScanPayRecordStatusFailure); got != "失败" {
		t.Fatalf("ScanPayRecordStatusName(Failure) = %q, want %q", got, "失败")
	}
}

func TestScanPayRecordStatusNameUnknownDoesNotPanic(t *testing.T) {
	if got := ScanPayRecordStatusName(999); got != "" {
		t.Fatalf("ScanPayRecordStatusName(unknown) = %q, want empty", got)
	}
}
