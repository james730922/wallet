package scanpaymember

import (
	"context"
	"strconv"

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

func newMemberPay() *scanPayPay {
	return &scanPayPay{}
}

type scanPayPay struct{}

func (hd *scanPayPay) Handler(ctx context.Context, req *zqbapis.ScanPayAddPayReq) (*zqbapis.ScanPayAddPayResp, error) {
	// 檢查參數
	if err := hd.checkReq(ctx, req); err != nil {
		return nil, err
	}

	//取得APP會員ID
	memberID, ok := ctxs.GetMemberID(ctx)
	if !ok {
		logger.ApLog().Error(errs.CommonNoMemberID)
		return nil, errs.CommonNoMemberID
	}
	//掃碼停權
	// 檢查會員所屬之會員分組是否正常
	if err := hd.checkMemberLevelIsAvailable(ctx, memberID); err != nil {
		return nil, err
	}

	//驗證請求參數
	recordID, err := strconv.ParseInt(req.Id, 10, 64)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, errs.CommonRequestParamParseFailed
	}

	record, err := packet.ScanPayRecord.GetScanPayRecord(packet.DB.New(), recordID)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, err
	}

	//驗證金額
	requestAmount := decimal.NewFromFloat(req.Amount)
	err = hd.checkAmount(requestAmount, record.Amount)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, err
	}

	//建立掃碼支付ByMember
	orderScanPayID, err := hd.checkStatusAndCreateOrder(ctx, record, memberID)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, err
	}

	//進行付款
	err = hd.orderToPay(ctx, record, memberID, orderScanPayID)
	if err != nil {
		if err != errs.WalletMemberUpdateBalanceIsNegative || err != errs.WalletMemberAmountUnreasonable {
			logger.ApLog().Errorf("memberID: %d, record: %s, orderScanPayID: %d, err: %v",
				memberID, tools.JsonMarshalString(record), orderScanPayID, err)
		}
		return nil, err
	}

	respData := &zqbapis.ScanPayAddPayResp{
		Status: models.CommonStatusSuccess,
		Amount: tools.DecimalToFloat64(record.Amount),
	}

	return respData, nil
}

func (hd *scanPayPay) checkReq(ctx context.Context, req *zqbapis.ScanPayAddPayReq) error {
	if req.Id == "" {
		return errs.CommonRequestParamInvalid
	}
	if req.Amount == 0 {
		return errs.CommonRequestParamInvalid
	}
	if req.Passwd == "" {
		return errs.CommonRequestParamInvalid
	}

	err := packet.UserMember.CheckScanPayPasswd(ctx, req.Passwd)
	if err != nil {
		logger.ApLog().Error(err)
		return err
	}
	return nil
}

func (hd *scanPayPay) checkAmount(reqAmt, recordAmt decimal.Decimal) error {
	if !reqAmt.Equal(recordAmt) {
		return errs.ScanPayAmountValidateFailed
	}
	return nil
}

func (hd *scanPayPay) checkStatusAndCreateOrder(ctx context.Context, scanPayRecord *model.ScanPayRecord, memberID int64) (int64, error) {
	switch scanPayRecord.Status {
	case model.ScanPayRecordStatusWaiting, model.ScanPayRecordStatusTransaction:
		//建立掃碼支付ByMember
		return hd.createOrderScanPay(ctx, scanPayRecord, memberID)
	case model.ScanPayRecordStatusDone:
		return 0, errs.ScanPayOrderAlreadyDone
	case model.ScanPayRecordStatusCancel:
		return 0, errs.ScanPayOrderCancel
	case model.ScanPayRecordStatusFailure:
		// Failure order/record 直接作廢，必須取得新的掃碼支付單才能再付款。
		return 0, errs.ScanPayOrderFailure
	}
	return 0, errs.ScanPayRecordStatusError
}

func (hd *scanPayPay) createOrderScanPay(ctx context.Context, scanPayRecord *model.ScanPayRecord, memberID int64) (int64, error) {
	orderScanPayID, err := packet.ScanPayCommon.CreateOrder(ctx, &condition.OrderScanPayCreateCond{
		MemberID:        memberID,
		Amount:          scanPayRecord.Amount,
		Brand:           aws.String(scanPayRecord.Brand),
		MerchantOrderID: aws.String(scanPayRecord.MerchantID),
		SourceOrderID:   aws.String(scanPayRecord.SourceOrderID),
		Content:         aws.String(scanPayRecord.Content),
		RecordID:        aws.Int64(scanPayRecord.ID),
	})
	if err != nil {
		logger.ApLog().Errorf("createOrderScanPay err: %s", err)
		return 0, err
	}
	return orderScanPayID, nil
}

func (hd *scanPayPay) orderToPay(ctx context.Context, scanPayRecord *model.ScanPayRecord, memberID, orderScanPayID int64) error {
	cond := &condition.OrderScanPayToPayCond{
		ID:       orderScanPayID,
		MemberID: aws.Int64(memberID),
		Amount:   tools.DecimalPtr(scanPayRecord.Amount),
	}
	orderScanPayId, err := packet.ScanPayCommon.OrderPay(ctx, cond)
	if err != nil {
		if err != errs.WalletMemberUpdateBalanceIsNegative || err != errs.WalletMemberAmountUnreasonable {
			logger.ApLog().Errorf("orderToPay err: %s, orderId: %s", err, orderScanPayId)
		}
		return err
	}
	return nil
}

func (hd *scanPayPay) checkMemberLevelIsAvailable(ctx context.Context, memberID int64) error {
	memberInfo, err := packet.UserMember.Get(ctx, memberID)
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
