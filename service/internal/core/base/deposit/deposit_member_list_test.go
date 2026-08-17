package deposit

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/james730922/wallet/service/internal/models"
	"github.com/james730922/wallet/service/internal/models/model"
)

func TestDepositListUsesBonusMatchedByDepositID(t *testing.T) {
	now := time.Now().UTC()
	deposits := []*model.Deposit{
		{ID: 101, MemberID: 7, Amount: decimal.NewFromInt(100), Status: model.DepositStatusAccept, AddedTime: now, UpdatedTime: now},
		{ID: 102, MemberID: 7, Amount: decimal.NewFromInt(200), Status: model.DepositStatusAccept, AddedTime: now, UpdatedTime: now},
	}
	bonuses := map[int64]*model.OrderBonus{
		101: {SourceOrderID: 101, Amount: decimal.NewFromInt(1), SourceRate: decimal.RequireFromString("0.01")},
		102: {SourceOrderID: 102, Amount: decimal.NewFromInt(4), SourceRate: decimal.RequireFromString("0.02")},
	}

	resp := (&depositList{}).parseResp(deposits, models.NewPagingResult(&models.Paging{Index: 1, Size: 25}, 2), bonuses, nil)
	if len(resp.Data) != 2 {
		t.Fatalf("response items = %d", len(resp.Data))
	}
	if resp.Data[0].BonusAmount != 1 || resp.Data[0].BonusRate != 1 {
		t.Fatalf("first deposit bonus = amount %v rate %v", resp.Data[0].BonusAmount, resp.Data[0].BonusRate)
	}
	if resp.Data[1].BonusAmount != 4 || resp.Data[1].BonusRate != 2 {
		t.Fatalf("second deposit bonus = amount %v rate %v", resp.Data[1].BonusAmount, resp.Data[1].BonusRate)
	}
}

func TestDepositDedupeKeyCanonicalizesAmount(t *testing.T) {
	first := depositDedupeKey(&model.Deposit{MemberID: 7, Amount: decimal.RequireFromString("100")})
	second := depositDedupeKey(&model.Deposit{MemberID: 7, Amount: decimal.RequireFromString("100.0")})
	if first != second {
		t.Fatalf("equivalent amounts produced different keys: %q != %q", first, second)
	}
}
