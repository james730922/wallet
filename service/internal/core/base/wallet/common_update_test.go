package wallet

import (
	"testing"

	"github.com/shopspring/decimal"

	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
)

func Test_commonUpdate_validate(t *testing.T) {
	logger.TestMock()

	type args struct {
		walletMember *model.WalletMember
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "正常",
			args: args{
				walletMember: &model.WalletMember{
					MemberID:     1232454,
					Balance:      decimal.NewFromInt(150),
					TotalAmount:  decimal.NewFromInt(230),
					Amount:       decimal.NewFromInt(200),
					Bonus:        decimal.NewFromInt(30),
					FrozenAmount: decimal.NewFromInt(80),
				},
			},
			wantErr: false,
		},
		{
			name: "異常-金額小於0",
			args: args{
				walletMember: &model.WalletMember{
					MemberID:     1232454,
					Balance:      decimal.NewFromInt(-150),
					TotalAmount:  decimal.NewFromInt(230),
					Amount:       decimal.NewFromInt(100),
					Bonus:        decimal.NewFromInt(30),
					FrozenAmount: decimal.NewFromInt(80),
				},
			},
			wantErr: true,
		},
		{
			name: "異常-小數位不正確",
			args: args{
				walletMember: &model.WalletMember{
					MemberID:     1232454,
					Balance:      decimal.RequireFromString("150.001"),
					TotalAmount:  decimal.NewFromInt(230),
					Amount:       decimal.NewFromInt(100),
					Bonus:        decimal.NewFromInt(30),
					FrozenAmount: decimal.NewFromInt(80),
				},
			},
			wantErr: true,
		},
		{
			name: "異常-加總錯誤",
			args: args{
				walletMember: &model.WalletMember{
					MemberID:     1232454,
					Balance:      decimal.NewFromInt(150),
					TotalAmount:  decimal.NewFromInt(230),
					Amount:       decimal.NewFromInt(100),
					Bonus:        decimal.NewFromInt(30),
					FrozenAmount: decimal.NewFromInt(80),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &commonUpdate{}
			if err := w.validateWallet(tt.args.walletMember); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
