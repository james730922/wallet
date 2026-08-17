package condition

import (
	"testing"

	"github.com/james730922/wallet/service/internal/models/model"
)

func TestOrderBonusQueryBuildsOwnershipAndSourceConditions(t *testing.T) {
	memberID := int64(11)
	sourceType := model.OrderBonusSourceTypeEnumDeposit
	sourceOrderIDs := []int64{101, 102}

	where := NewQuery(&OrderBonusQuery{
		MemberID:       &memberID,
		SourceType:     &sourceType,
		SourceOrderIDs: &sourceOrderIDs,
	}).Where()

	if got := where["member_id"]; got != &memberID {
		t.Fatalf("member_id condition = %#v", got)
	}
	if got := where["source_type"]; got != &sourceType {
		t.Fatalf("source_type condition = %#v", got)
	}
	gotIDs, ok := where["in_source_order_id"].(*[]int64)
	if !ok || len(*gotIDs) != 2 || (*gotIDs)[0] != 101 || (*gotIDs)[1] != 102 {
		t.Fatalf("source order IDs condition = %#v", where["in_source_order_id"])
	}
}
