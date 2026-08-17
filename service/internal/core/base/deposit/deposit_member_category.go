package deposit

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newDepositCategory() *depositCategory {
	return &depositCategory{}
}

type depositCategory struct {
	mx sync.Mutex
}

func (de *depositCategory) Handler(ctx context.Context) (map[int64]string, error) {
	cond := &condition.BankDepositCategoryItemViewCond{
		AccountStatus:  aws.Int(int(model.BankAccountStatusEnable)),
		CategoryStatus: aws.Int(int(model.DepositCategoryStatusEnable)),
	}
	items, err := self.DepositCommon.CategoryItem(ctx, cond)
	if err != nil {
		return nil, errs.CommonNoData
	}

	result := make(map[int64]string)
	for _, v := range items {
		result[v.AccountID] = v.CategoryName
	}
	return result, nil
}
