package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type WalletMember struct {
	MemberID     int64           // '會員id'
	Balance      decimal.Decimal // '可用餘額'
	TotalAmount  decimal.Decimal // '錢包餘額'
	Amount       decimal.Decimal // '入款餘額'
	Bonus        decimal.Decimal // '紅利餘額'
	FrozenAmount decimal.Decimal // '凍結額'
	Sign         string          // '簽名'
	AddedTime    time.Time       // '創建時間'
	UpdatedTime  time.Time       // '更新時間'
}

func (WalletMember) TableName() string {
	return "wallet_member"
}

func ConvertWalletMember(walletMember *WalletMember) *WalletMemberTypeDecimal {
	return &WalletMemberTypeDecimal{
		Balance:      walletMember.Balance,
		TotalAmount:  walletMember.TotalAmount,
		Amount:       walletMember.Amount,
		Bonus:        walletMember.Bonus,
		FrozenAmount: walletMember.FrozenAmount,
	}
}

type WalletMemberTypeDecimal struct {
	Balance      decimal.Decimal // '可用餘額'
	TotalAmount  decimal.Decimal // '錢包餘額'
	Amount       decimal.Decimal // '入款餘額'
	Bonus        decimal.Decimal // '紅利餘額'
	FrozenAmount decimal.Decimal // '凍結額'
}

// 異動結果
type WalletMemberModifyResult struct {
	CurrentTotalAmount decimal.Decimal // 原資產
	ChangedTotalAmount decimal.Decimal // 異動後資產
	CurrentBalance     decimal.Decimal // 原錢包餘額
	ChangedBalance     decimal.Decimal // 異動後錢包餘額
	Amount             decimal.Decimal // 異動金額
	UsedBonus          decimal.Decimal // 使用到的bonus
}

// 異動後通知
type WalletMemberWithChangeNotify struct {
	MemberID int64
}
