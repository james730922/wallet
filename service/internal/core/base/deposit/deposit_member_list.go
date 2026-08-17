package deposit

import (
	"context"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/shopspring/decimal"
	"github.com/james730922/wallet/service/internal/models"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/pb/zqbapis"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/ctxs"
	"github.com/james730922/wallet/service/internal/utils/errs"
	"github.com/james730922/wallet/service/internal/utils/tools"
)

func newDepositList() *depositList {
	return &depositList{}
}

type depositList struct{}

func (de *depositList) Handler(ctx context.Context, req *zqbapis.DepositListReq) (*zqbapis.DepositListResp, error) {
	query, err := de.parseReq(ctx, req)
	if err != nil {
		return nil, err
	}

	depositList, paging, err := dao.Common.List(packet.DB.New(), condition.NewQuery(query))
	if err != nil {
		logger.ApLog().Errorf("List error:%s, depositListReq:%v", err.Error(), req)
		return nil, err
	}

	bonusByDepositID, config, err := de.loadDisplayData(depositList)
	if err != nil {
		return nil, err
	}

	resp := de.parseResp(depositList, paging, bonusByDepositID, config)
	return resp, nil
}

func (de *depositList) parseReq(ctx context.Context, req *zqbapis.DepositListReq) (*condition.DepositQuery, error) {
	memberID, ok := ctxs.GetMemberID(ctx)
	if !ok {
		logger.ApLog().Error(errs.CommonNoMemberID)
		return nil, errs.CommonNoMemberID
	}
	queryCond := &condition.DepositQuery{
		Paging:   models.NewPagingFromProto(req.Paging),
		MemberID: aws.Int64(memberID),
	}
	if req.Status != model.DepositStatusAll {
		queryCond.Status = aws.Int(int(req.Status))
	}
	return queryCond, nil
}

func (de *depositList) loadDisplayData(deposits []*model.Deposit) (map[int64]*model.OrderBonus, *model.DepositConfigForMember, error) {
	bonusByDepositID := make(map[int64]*model.OrderBonus)
	acceptedIDs := make([]int64, 0, len(deposits))
	var memberID int64
	needsConfig := false

	for _, deposit := range deposits {
		memberID = deposit.MemberID
		if deposit.Status == model.DepositStatusAccept {
			acceptedIDs = append(acceptedIDs, deposit.ID)
		} else {
			needsConfig = true
		}
	}

	if len(acceptedIDs) > 0 {
		sourceType := model.OrderBonusSourceTypeEnumDeposit
		bonuses, err := self.DepositBonus.ListDepositBonus(packet.DB.New(), &condition.OrderBonusQuery{
			MemberID:       &memberID,
			SourceType:     &sourceType,
			SourceOrderIDs: &acceptedIDs,
		})
		if err != nil {
			logger.ApLog().Error(err)
			return nil, nil, err
		}
		for _, bonus := range bonuses {
			bonusByDepositID[bonus.SourceOrderID] = bonus
		}
	}

	var config *model.DepositConfigForMember
	if needsConfig && len(deposits) > 0 {
		var err error
		config, err = self.DepositConfig.GetDepositConfigForMember(memberID)
		if err != nil && err != errs.DBNoRow {
			logger.ApLog().Error(err)
			return nil, nil, err
		}
	}

	return bonusByDepositID, config, nil
}

func (de *depositList) parseResp(
	depositList []*model.Deposit,
	paging *models.PagingResult,
	bonusByDepositID map[int64]*model.OrderBonus,
	config *model.DepositConfigForMember,
) *zqbapis.DepositListResp {
	var respData []*zqbapis.DepositList
	for _, d := range depositList {
		tmp := &zqbapis.DepositList{
			Id:           strconv.FormatInt(d.ID, 10),
			AddedTime:    d.AddedTime.Format(time.RFC3339),
			UpdateTime:   d.UpdatedTime.Format(time.RFC3339),
			CurrencyCode: d.CurrencyCode,
			Amount:       tools.DecimalToFloat64(d.Amount),
			Status:       d.Status.MemberName(),
		}
		// 取得紅利的資料
		switch d.Status {
		case model.DepositStatusAccept:
			if bonusOrder, ok := bonusByDepositID[d.ID]; ok {
				tmp.BonusAmount = tools.DecimalToFloat64(bonusOrder.Amount.RoundBank(2))
				tmp.BonusRate = tools.DecimalToFloat64(bonusOrder.SourceRate.Mul(decimal.NewFromInt(100)).RoundBank(2))
			}
		default:
			if config != nil {
				tmp.BonusRate = tools.DecimalToFloat64(config.Bonus.Mul(decimal.NewFromInt(100)).RoundBank(2))
				tmp.BonusAmount = tools.DecimalToFloat64(config.Bonus.Mul(d.Amount).RoundBank(2))
			}
		}
		respData = append(respData, tmp)
	}
	page := &zqbapis.PagingResult{
		PageIndex: int32(paging.Index),
		PageSize:  int32(paging.Size),
		TotalPage: int32(paging.TotalPage),
		TotalRow:  int64(paging.TotalRow),
	}
	resp := &zqbapis.DepositListResp{
		Paging: page,
		Data:   respData,
	}
	return resp
}
