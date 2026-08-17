package deposit

import (
	"context"

	"github.com/james730922/wallet/service/internal/pb/zqbapis"
)

func newDepositMember() IDepositMember {
	return &depositMemberUseCase{
		list:     newDepositList(),
		methods:  newDepositMethod(),
		order:    newDepositOrder(),
		category: newDepositCategory(),
	}
}

type IDepositMember interface {
	List(ctx context.Context, req *zqbapis.DepositListReq) (*zqbapis.DepositListResp, error)
	Methods(ctx context.Context) (*zqbapis.DepositMethodResp, error)
	Order(ctx context.Context, req *zqbapis.DepositOrderReq) (*zqbapis.DepositOrderResp, error)
	Category(ctx context.Context) (map[int64]string, error)
}

type depositMemberUseCase struct {
	list     *depositList
	methods  *depositMemberMethods
	order    *depositOrder
	category *depositCategory
}

func (de *depositMemberUseCase) List(ctx context.Context, req *zqbapis.DepositListReq) (*zqbapis.DepositListResp, error) {
	return de.list.Handler(ctx, req)
}

func (de *depositMemberUseCase) Methods(ctx context.Context) (*zqbapis.DepositMethodResp, error) {
	return de.methods.Handler(ctx)
}

func (de *depositMemberUseCase) Order(ctx context.Context, req *zqbapis.DepositOrderReq) (*zqbapis.DepositOrderResp, error) {
	return de.order.Handler(ctx, req)
}

func (de *depositMemberUseCase) Category(ctx context.Context) (map[int64]string, error) {
	return de.category.Handler(ctx)
}
