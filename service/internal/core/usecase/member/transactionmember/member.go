package transactionmember

import (
	"context"

	"github.com/james730922/wallet/service/internal/pb/zqbapis"
)

func newMember() *memberUseCase {
	return &memberUseCase{}
}

type IMember interface {
	List(ctx context.Context, req *zqbapis.TransactionReq) (*zqbapis.TransactionListResp, error)
	Detail(ctx context.Context, id int64) (*zqbapis.TransactionDetailResp, error)
}

type memberUseCase struct {
	list   *memberList
	detail *memberDetail
}

func (uc *memberUseCase) List(ctx context.Context, req *zqbapis.TransactionReq) (*zqbapis.TransactionListResp, error) {
	return uc.list.Handler(ctx, req)
}

func (uc *memberUseCase) Detail(ctx context.Context, id int64) (*zqbapis.TransactionDetailResp, error) {
	return uc.detail.Handler(ctx, id)
}
