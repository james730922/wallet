package deposit

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
)

func TestDepositCategoryItemFilterUsesCategoryStatus(t *testing.T) {
	items := []*model.BankDepositCategoryItemView{
		{CategoryID: 1, AccountStatus: 1, CategoryStatus: 0},
		{CategoryID: 2, AccountStatus: 0, CategoryStatus: 1},
	}

	got := newDepositCommonCategoryItem().filter(items, &condition.BankDepositCategoryItemViewCond{
		AccountStatus:  aws.Int(0),
		CategoryStatus: aws.Int(1),
	})

	if len(got) != 1 || got[0].CategoryID != 2 {
		t.Fatalf("filter() = %#v, want category 2", got)
	}
}
