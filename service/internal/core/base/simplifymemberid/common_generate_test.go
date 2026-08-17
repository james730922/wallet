package simplifymemberid

import (
	"testing"

	"github.com/james730922/wallet/service/internal/utils/errs"
)

func TestFormatSimplifyID(t *testing.T) {
	tests := []struct {
		name    string
		count   int64
		want    string
		wantErr error
	}{
		{name: "minimum", count: 1, want: "00001"},
		{name: "five digits", count: 99999, want: "99999"},
		{name: "six digits", count: 100000, want: "100000"},
		{name: "maximum", count: maxSimplifyID, want: "99999999"},
		{name: "invalid zero", count: 0, wantErr: errs.CommonNoMemberID},
		{name: "exhausted", count: maxSimplifyID + 1, wantErr: errs.MemberSimplifyIDExhausted},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatSimplifyID(tt.count)
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Fatalf("formatSimplifyID() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("formatSimplifyID() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("formatSimplifyID() = %q, want %q", got, tt.want)
			}
		})
	}
}
