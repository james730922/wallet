package deposit

import (
	"context"
	"strconv"
	"time"

	"github.com/shopspring/decimal"

	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/errs"
	"github.com/james730922/wallet/service/internal/utils/signs"
)

type depositCommonOrder struct{}

func (de *depositCommonOrder) Handler(ctx context.Context, cond *model.Deposit) (int64, error) {
	orderDeposit := de.getDeposit(cond)
	claimOwner := strconv.FormatInt(orderDeposit.ID, 10)
	claimed, err := self.DepositCache.claimDepositWithSameMemberAndAmount(cond, claimOwner)
	if err != nil {
		logger.ApLog().Errorf("claim deposit dedupe key failed: memberID=%d err=%v", cond.MemberID, err)
		return 0, errs.CommonServiceUnavailable
	}
	if !claimed {
		logger.ApLog().Warnf("duplicate deposit rejected: memberID=%d amount=%s", cond.MemberID, cond.Amount)
		return 0, errs.DepositTooOften
	}

	if err := dao.Common.Insert(packet.DB.New(), orderDeposit); err != nil {
		self.DepositCache.releaseDepositWithSameMemberAndAmount(cond, claimOwner)
		logger.ApLog().Error(err)
		return 0, errs.OrderDepositCreateFailed
	}

	return orderDeposit.ID, nil
}

func (de *depositCommonOrder) getDeposit(cond *model.Deposit) *model.Deposit {
	id := packet.Node.Generate().Int64()

	now := time.Now().UTC()

	orderDeposit := &model.Deposit{
		ID:              id,
		MemberID:        cond.MemberID,
		AccountID:       cond.AccountID,
		AccountNumber:   cond.AccountNumber,
		AccountBankCode: cond.AccountBankCode,
		CurrencyCode:    cond.CurrencyCode,
		PayName:         cond.PayName,
		Amount:          cond.Amount,
		Charge:          decimal.Zero,
		Status:          model.DepositStatusWaiting,
		AddedTime:       now,
		UpdatedTime:     now,
	}

	orderDeposit.Sign = signs.Deposit(orderDeposit)

	return orderDeposit
}
