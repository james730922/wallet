package bank

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newCategoryCommon() ICategoryCommon {
	return &categoryCommonUseCase{
		get:                 newCategoryCommonGet(),
		getItemMapByAccount: newCategoryCommonItemGetMapByAccountID(),
		getType:             newCategoryCommonTypeGet(),
	}
}

type ICategoryCommon interface {
	First(ctx context.Context, id int64) (*model.DepositCategory, error)
	List(ctx context.Context, cond *condition.DepositCategoryQuery) ([]*model.DepositCategory, error)
	GetItemMapByAccount(ctx context.Context) (map[int64]*model.BankDepositCategoryItem, error)

	ListType(ctx context.Context) ([]*model.BankDepositCategoryType, error)
	ListTypeMap(ctx context.Context) (map[int]*model.BankDepositCategoryType, error)
	GetBankDepositCategoryTypeName(ctx context.Context, id int) (string, error)
}

type categoryCommonUseCase struct {
	get                 *categoryCommonGet
	getItemMapByAccount *categoryCommonItemGetMapByAccountID
	getType             *categoryCommonTypeGet
}

func (c *categoryCommonUseCase) First(ctx context.Context, id int64) (*model.DepositCategory, error) {
	items, err := c.get.Handler(ctx, &condition.DepositCategoryQuery{ID: aws.Int64(id)})
	if err != nil {
		return nil, err
	}

	if len(items) != 1 {
		return nil, errs.DBNoRow
	}

	return items[0], nil
}

func (c *categoryCommonUseCase) List(ctx context.Context, cond *condition.DepositCategoryQuery) ([]*model.DepositCategory, error) {
	return c.get.Handler(ctx, cond)
}

func (c *categoryCommonUseCase) GetItemMapByAccount(ctx context.Context) (map[int64]*model.BankDepositCategoryItem, error) {
	return c.getItemMapByAccount.Handler(ctx)
}

func (c *categoryCommonUseCase) ListType(ctx context.Context) ([]*model.BankDepositCategoryType, error) {
	return c.getType.Handler(ctx)
}
func (c *categoryCommonUseCase) ListTypeMap(ctx context.Context) (map[int]*model.BankDepositCategoryType, error) {
	data, err := c.getType.Handler(ctx)
	if err != nil {
		return nil, err
	}

	result := make(map[int]*model.BankDepositCategoryType)
	for _, d := range data {
		result[d.ID] = d
	}

	return result, nil
}

func (c *categoryCommonUseCase) GetBankDepositCategoryTypeName(ctx context.Context, id int) (string, error) {
	bankDepositCategoryTypeMap, err := c.ListTypeMap(ctx)
	if err != nil {
		return "", err
	}
	if data, ok := bankDepositCategoryTypeMap[id]; ok {
		return data.Name, nil
	}
	return "", nil

}
