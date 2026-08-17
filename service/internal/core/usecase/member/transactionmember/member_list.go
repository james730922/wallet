package transactionmember

import (
	"context"
	"strconv"
	"time"

	"github.com/james730922/wallet/service/internal/models"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/pb/zqbapis"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/ctxs"
	"github.com/james730922/wallet/service/internal/utils/errs"
	"github.com/james730922/wallet/service/internal/utils/tools"
)

type memberList struct{}

func (de *memberList) Handler(ctx context.Context, req *zqbapis.TransactionReq) (*zqbapis.TransactionListResp, error) {
	cond, err := de.getCond(ctx, req)
	if err != nil {
		return nil, err
	}

	txns, paging, err := packet.Transaction.ListWithPaging(packet.DB.New(), cond)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, err
	}

	resp, err := de.arrangeResp(txns, paging)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (de *memberList) getCond(ctx context.Context, req *zqbapis.TransactionReq) (*condition.TransactionQuery, error) {
	memberID, ok := ctxs.GetMemberID(ctx)
	if !ok {
		return nil, errs.CommonNoMemberID
	}

	result := &condition.TransactionQuery{
		Paging:   models.NewPagingFromProto(req.Paging),
		MemberID: &memberID,
	}

	var sourceTypes []int

	switch req.Class {
	case zqbapis.TransactionClassType_Deposit:
		sourceTypes = append(sourceTypes,
			int(model.TransactionSourceTypeDepositBonus),
			int(model.TransactionSourceTypeDeposit))
		result.SourceTypes = &sourceTypes
	case zqbapis.TransactionClassType_Scanpay:
		sourceTypes = append(sourceTypes,
			int(model.TransactionSourceTypeScanPayConfirm),
		)
		result.SourceTypes = &sourceTypes
	default:
	}

	return result, nil
}

func (de *memberList) arrangeResp(txns []*model.Transaction, paging *models.PagingResult) (*zqbapis.TransactionListResp, error) {
	var data []*zqbapis.TransactionList
	for _, v := range txns {
		data = append(data, &zqbapis.TransactionList{
			Id:             strconv.FormatInt(v.ID, 10),
			Class:          v.SourceType.ClassName(),
			SourceType:     v.SourceType.Name(),
			Amount:         tools.DecimalToFloat64(v.Amount),
			ChangedBalance: tools.DecimalToFloat64(v.ChangedBalance),
			AddedTime:      v.AddedTime.UTC().Format(time.RFC3339Nano),
			BrandName:      "",
		})
	}

	page := &zqbapis.PagingResult{
		PageIndex: int32(paging.Index),
		PageSize:  int32(paging.Size),
		TotalPage: int32(paging.TotalPage),
		TotalRow:  int64(paging.TotalRow),
	}

	resp := &zqbapis.TransactionListResp{
		Paging: page,
		Data:   data,
	}

	return resp, nil
}
