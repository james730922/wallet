package signs

import (
	"bytes"
	"crypto/sha512"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/utils/conf"
)

func Hash(data []byte) [64]byte {
	var buf bytes.Buffer
	buf.Write(data)
	buf.Write(conf.Sign().GetSalt())

	return sha512.Sum512(buf.Bytes())
}

func Hex(args ...interface{}) string {
	s := fmt.Sprint(args...)
	return fmt.Sprintf("%x", Hash([]byte(s)))
}

func WalletMember(s *model.WalletMember) string {
	return Hex(
		s.MemberID,
		s.Balance,
		s.TotalAmount,
		s.Amount,
		s.Bonus,
		s.FrozenAmount,
		timeToString(s.AddedTime),
		timeToString(s.UpdatedTime),
	)
}

func OrderBonus(s *model.OrderBonus) string {
	return Hex(
		s.ID,
		s.MemberID,
		s.CurrencyCode,
		s.Amount,
		s.CalculateType,
		s.SourceType,
		s.SourceAmount,
		s.SourceRate,
		s.Status,
		timeToString(s.AddedTime),
		timeToString(s.UpdatedTime),
	)
}

func Deposit(s *model.Deposit) string {
	return Hex(
		s.ID,
		s.MemberID,
		s.AccountID,
		s.AccountNumber,
		s.AccountBankCode,
		s.CurrencyCode,
		s.PayName,
		s.Amount,
		s.Charge,
		s.Status,
		timeToString(s.AddedTime),
		timeToString(s.UpdatedTime),
	)
}

func Transaction(s *model.Transaction) string {
	return Hex(
		s.ID,
		s.MemberID,
		s.SourceType,
		s.SourceID,
		s.CurrencyCode,
		s.Amount,
		s.CurrentTotalAmount,
		s.ChangedTotalAmount,
		s.CurrentBalance,
		s.ChangedBalance,
		timeToString(s.AddedTime),
		timeToString(s.UpdatedTime),
		s.Remarks,
	)
}

func timeToString(t time.Time) string {
	return t.UTC().Truncate(time.Second).String()
}

func OrderScanPay(s *model.OrderScanPay) string {
	return Hex(
		s.ID,
		s.MemberID,
		s.Amount,
		s.Status,
		aws.Int64Value(s.RecordID),
		timeToString(s.AddedTime),
		timeToString(s.UpdatedTime),
	)
}
