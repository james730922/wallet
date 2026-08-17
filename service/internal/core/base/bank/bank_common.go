package bank

import (
	"context"

	"github.com/james730922/wallet/service/internal/utils/errs"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
)

func newBankCommon() IBankCommon {
	bankCommon := &bankCommonUseCase{
		_map: newBankCommonList(),
	}
	return bankCommon
}

type IBankCommon interface {
	GetByCode(ctx context.Context, code string) *model.BankCode
	First(ctx context.Context, cond *condition.BankCodeQuery) (*model.BankCode, error)
	List(ctx context.Context, cond *condition.BankCodeQuery) ([]*model.BankCode, error)
	Map(ctx context.Context, cond *condition.BankCodeQuery) (map[string]*model.BankCode, error)
}

type bankCommonUseCase struct {
	_map *bankCommonList
}

func (b *bankCommonUseCase) GetByCode(ctx context.Context, code string) *model.BankCode {
	data, err := b.List(ctx, &condition.BankCodeQuery{Code: aws.String(code)})
	if err != nil {
		return nil
	}
	if len(data) == 0 {
		return nil
	}

	return data[0]
}

func (b *bankCommonUseCase) First(ctx context.Context, cond *condition.BankCodeQuery) (*model.BankCode, error) {
	data, err := b.List(ctx, cond)
	if err != nil {
		return nil, err
	}

	if len(data) != 1 {
		return nil, errs.CommonNoData
	}

	return data[0], nil
}

func (b *bankCommonUseCase) List(ctx context.Context, cond *condition.BankCodeQuery) ([]*model.BankCode, error) {
	return b._map.Handler(ctx, cond)
}

func (b *bankCommonUseCase) Map(ctx context.Context, cond *condition.BankCodeQuery) (map[string]*model.BankCode, error) {
	data, err := b.List(ctx, cond)
	if err != nil {
		return nil, err
	}

	res := make(map[string]*model.BankCode, len(data))
	for _, v := range data {
		res[v.Code] = v
	}

	return res, nil
}
