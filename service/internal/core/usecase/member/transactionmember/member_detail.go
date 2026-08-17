package transactionmember

import (
	"context"
	"strconv"
	"time"

	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/pb/zqbapis"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/ctxs"
	"github.com/james730922/wallet/service/internal/utils/errs"
	"github.com/james730922/wallet/service/internal/utils/tools"
)

type memberDetail struct{}

func (de *memberDetail) Handler(ctx context.Context, id int64) (*zqbapis.TransactionDetailResp, error) {
	txn, err := packet.Transaction.GetOneByID(packet.DB.New(), id)
	if err != nil {
		return nil, err
	}

	memberID, ok := ctxs.GetMemberID(ctx)
	if !ok {
		return nil, errs.CommonNoMemberID
	}

	if txn.MemberID != memberID {
		return nil, errs.CommonRecordOwnershipWrong
	}

	data := de.arrangeData(ctx, txn)

	resp := &zqbapis.TransactionDetailResp{
		Data: data,
	}

	return resp, nil
}

func (de *memberDetail) arrangeData(ctx context.Context, txn *model.Transaction) *zqbapis.TransactionDetail {
	simpleID, err := packet.SimplifyMemberID.MappingSimplifyID(txn.MemberID)
	if err != nil {
		logger.ApLog().Errorf("memberID:%v, err:%v", simpleID, err)
	}

	resp := &zqbapis.TransactionDetail{
		Id:                 strconv.FormatInt(txn.ID, 10),
		MemberId:           simpleID,
		SourceId:           strconv.FormatInt(txn.SourceID, 10),
		SourceType:         txn.SourceType.Name(),
		Amount:             tools.DecimalToFloat64(txn.Amount),
		CurrentTotalAmount: tools.DecimalToFloat64(txn.CurrentTotalAmount),
		ChangedTotalAmount: tools.DecimalToFloat64(txn.ChangedTotalAmount),
		CurrentBalance:     tools.DecimalToFloat64(txn.CurrentBalance),
		ChangedBalance:     tools.DecimalToFloat64(txn.ChangedBalance),
		AddedTime:          txn.AddedTime.UTC().Format(time.RFC3339Nano),
		Remarks:            txn.Remarks,
	}

	return resp
}
