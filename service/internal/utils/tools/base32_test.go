package tools

import (
	"testing"
)

func Test_base32WithInt64_Encode(t *testing.T) {
	type args struct {
		memberID int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "正常",
			args: args{
				memberID: 1321464636729462784,
			},
			want: "AAIIB3YVZFLBE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := newBase32WithInt64()
			if got := b.Encode(tt.args.memberID); got != tt.want {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_base32WithInt64_Decode(t *testing.T) {
	type args struct {
		memberID string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "正常",
			args: args{
				memberID: "AAIIB3YVZFLBE",
			},
			want:    1321464636729462784,
			wantErr: false,
		},
		{
			name: "錯誤",
			args: args{
				memberID: "AAIIB3YVZFLB",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := newBase32WithInt64()
			got, err := b.Decode(tt.args.memberID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Decode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
