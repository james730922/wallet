package deposit

import (
	"context"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
)

func newDepositCommon() IDepositCommon {
	return &depositCommonUseCase{
		categoryItem: newDepositCommonCategoryItem(),
	}
}

type IDepositCommon interface {
	Order(ctx context.Context, cond *model.Deposit) (int64, error)
	CategoryItem(ctx context.Context, cond *condition.BankDepositCategoryItemViewCond) ([]*model.BankDepositCategoryItemView, error)
	GetAccount(ctx context.Context, cond *condition.BankAccountQuery) (*model.BankAccount, error)
}

type depositCommonUseCase struct {
	order        *depositCommonOrder
	categoryItem *depositCommonCategoryItem
}

// 會員入款
func (de *depositCommonUseCase) Order(ctx context.Context, cond *model.Deposit) (int64, error) {
	return de.order.Handler(ctx, cond)
}

func (de *depositCommonUseCase) CategoryItem(ctx context.Context, cond *condition.BankDepositCategoryItemViewCond) ([]*model.BankDepositCategoryItemView, error) {
	return de.categoryItem.Handler(ctx, cond)
}

func (de *depositCommonUseCase) GetAccount(ctx context.Context, cond *condition.BankAccountQuery) (*model.BankAccount, error) {
	return dao.Account.First(packet.DB.New(), condition.NewQuery(cond))
}
