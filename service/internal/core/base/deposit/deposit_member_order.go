package deposit

import (
	"context"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/pb/zqbapis"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/ctxs"
	"github.com/james730922/wallet/service/internal/utils/errs"
	"github.com/james730922/wallet/service/internal/utils/tools"
)

func newDepositOrder() *depositOrder {
	return &depositOrder{}
}

type depositOrder struct{}

func (de *depositOrder) Handler(ctx context.Context, req *zqbapis.DepositOrderReq) (*zqbapis.DepositOrderResp, error) {
	memberID, ok := ctxs.GetMemberID(ctx)
	if !ok {
		return nil, errs.CommonNoMemberID
	}

	// 檢查會員所屬之會員分組是否正常
	if err := de.checkMemberLevelIsAvailable(ctx, memberID); err != nil {
		return nil, err
	}

	// 取得銀行資料
	bankAccount, err := de.getBankAccount(ctx, req.AccountId)
	if err != nil {
		return nil, err
	}

	// 驗證參數
	if err := de.validate(ctx, req, bankAccount); err != nil {
		return nil, err
	}

	// 編輯入款條件
	cond := de.genOrderCond(ctx, memberID, req, bankAccount)

	// 入款
	id, err := self.DepositCommon.Order(ctx, cond)
	if err != nil {
		logger.ApLog().Warnf("err:%v, cond:%v", err, tools.JsonMarshalString(cond))
		return nil, err
	}

	resp := &zqbapis.DepositOrderResp{
		Id: strconv.FormatInt(id, 10),
	}

	return resp, nil
}
func (de *depositOrder) checkMemberLevelIsAvailable(ctx context.Context, memberID int64) error {
	memberInfo, err := packet.Member.Get(ctx, memberID)
	if err != nil {
		return err
	}

	mLevelInfo, err := packet.MemberLevel.First(packet.DB.New(), &condition.MemberLevelQuery{
		ID: &memberInfo.LevelCode,
	})
	if err != nil {
		return err
	}

	if mLevelInfo.Feature == model.MemberLevelFeatureBlackList {
		return errs.CommonMemberLevelIsBlackList
	}

	return nil
}

func (de *depositOrder) getBankAccount(ctx context.Context, bankAccountIDStr string) (*model.BankAccount, error) {
	bankAccountID, err := strconv.ParseInt(bankAccountIDStr, 10, 64)
	if err != nil {
		logger.ApLog().Warnf("err:%v, bankAccountIDStr:%v", err, bankAccountIDStr)
		return nil, errs.CommonRequestParamParseFailed
	}

	bankAccount, err := self.DepositCommon.GetAccount(ctx, &condition.BankAccountQuery{ID: &bankAccountID})
	if err != nil {
		logger.ApLog().Warnf("err:%v, bankAccountID:%v", err, bankAccountID)
		return nil, errs.CommonNoData
	}

	return bankAccount, nil
}

func (de *depositOrder) validate(ctx context.Context, req *zqbapis.DepositOrderReq, account *model.BankAccount) error {

	// 檢查入款人戶名
	if err := de.checkPayName(ctx, req); err != nil {
		return err
	}

	// 檢查入款金額
	if err := de.checkAmount(req, account); err != nil {
		return err
	}

	return nil
}

func (de *depositOrder) checkPayName(ctx context.Context, req *zqbapis.DepositOrderReq) error {
	accountID, err := strconv.ParseInt(req.AccountId, 10, 64)
	if err != nil {
		logger.ApLog().Warnf("err:%v, accountID:%v", err, req.AccountId)
		return errs.DepositErrorParamType
	}

	// 取得入款帳戶的分類
	cats, err := packet.CategoryCommon.GetItemMapByAccount(ctx)
	if err != nil {
		logger.SysLog().Errorf("err:%v", err)
		return errs.DepositErrorParamType
	}

	cat, ok := cats[accountID]
	if !ok {
		logger.ApLog().Warnf("cat:%v, accountID:%v", tools.JsonMarshalString(cats), req.AccountId)
		return errs.DepositErrorParamType
	}

	catType, err := packet.CategoryCommon.First(ctx, cat.CategoryID)
	if err != nil {
		logger.ApLog().Warnf("err:%v, categoryID:%v", err, cat.CategoryID)
		return errs.DepositErrorParamType
	}

	// 若入款帳戶為銀行，則入款人戶名不可為空
	if catType.Type == model.BankAccountTypeBank && strings.TrimSpace(req.PayName) == "" {
		return errs.DepositErrorParamPayName
	}

	return nil
}

func (de *depositOrder) checkAmount(req *zqbapis.DepositOrderReq, account *model.BankAccount) error {
	amount := decimal.NewFromFloat(req.Amount)
	if !amount.IsPositive() {
		logger.ApLog().Warnf("amount:%v <= 0", req.Amount)
		return errs.DepositErrorParamAmount
	}

	if !amount.Equal(amount.Truncate(0)) {
		logger.ApLog().Warnf("amount:%v isn't int", req.Amount)
		return errs.DepositAmountShouldBeInt
	}

	// 檢查是否超過限制範圍
	if !account.MinAmount.IsZero() && account.MinAmount.GreaterThan(amount) {
		logger.ApLog().Warnf("account.MinAmount:%v > amount:%v", account.MinAmount, req.Amount)
		return errs.DepositAmountExceedRange
	}

	if !account.MaxAmount.IsZero() && account.MaxAmount.LessThan(amount) {
		logger.ApLog().Warnf("account.MaxAmount:%v < amount:%v", account.MaxAmount, req.Amount)
		return errs.DepositAmountExceedRange
	}

	return nil
}

func (de *depositOrder) genOrderCond(ctx context.Context, memberID int64, req *zqbapis.DepositOrderReq, account *model.BankAccount) *model.Deposit {
	return &model.Deposit{
		MemberID:        memberID,
		AccountID:       account.ID,
		AccountNumber:   account.Number,
		AccountBankCode: account.BankCode,
		CurrencyCode:    account.CurrencyCode,
		PayName:         req.PayName,
		Amount:          decimal.NewFromFloat(req.Amount),
	}
}
