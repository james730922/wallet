package signs

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
)

// TestWalletMember 簽名驗證用 func ，如過要測試要另外寫
func TestWalletMember(t *testing.T) {
	// 要依環境調整鹽

	logger.TestMock()

	type args struct {
		s *model.WalletMember
	}

	timeParse := func(val string) time.Time {
		t, err := time.Parse("2006-01-02 15:04:05.999999999", val)
		if err != nil {
			panic(err)
		}
		return t
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			// 重新簽名的方法，例如發生簽名異常時使用
			name: "重新簽名",
			args: args{s: &model.WalletMember{
				MemberID:     1418496741044391936,
				Balance:      decimal.RequireFromString("8677.07"),
				TotalAmount:  decimal.RequireFromString("9077.07"),
				Amount:       decimal.RequireFromString("9077.07"),
				Bonus:        decimal.Zero,
				FrozenAmount: decimal.NewFromInt(400),
				AddedTime:    timeParse("2021-07-23 09:02:32.128849"),
				UpdatedTime:  timeParse("2021-09-15 15:28:30.609117"),
			}},
		},
		{
			// 到對應環境抓一筆正確的資料下來，驗證 sing 出來是否一樣
			name: "簽名正確測試",
			args: args{s: &model.WalletMember{
				MemberID:     1357236628237586432,
				Balance:      decimal.RequireFromString("2620.44"),
				TotalAmount:  decimal.RequireFromString("2620.44"),
				Amount:       decimal.RequireFromString("2620.44"),
				Bonus:        decimal.Zero,
				FrozenAmount: decimal.Zero,
				AddedTime:    timeParse("2021-02-04 07:56:43.077646"),
				UpdatedTime:  timeParse("2021-08-01 01:09:26.628838"),
			}},
			want: "ce821201bf05d6530b3feb506709cb8a5dca2ec65e966a6867701af78b5ebf371d1049c6b3f314c37c2d324858b7029cc4a606465f327997b44dc71f0a2ec9d1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WalletMember(tt.args.s)
			if tt.name == "重新簽名" {
				t.Logf("重新簽名 got = %v", got)
			}

			t.Logf("got == want: %v\ngot: %v\nwant: %v", got == tt.want, got, tt.want)
		})
	}
}
