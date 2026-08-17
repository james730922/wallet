package scanpaymember

import (
	"context"
	"fmt"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/pb/zqbapis"
	"github.com/james730922/wallet/service/internal/utils/errs"
	"github.com/james730922/wallet/service/internal/utils/tools"
	"time"
)

func newMemberScan() *scanPayScan {
	return &scanPayScan{}
}

type scanPayScan struct{}

func (hd *scanPayScan) Handler(ctx context.Context, req *zqbapis.ScanPayScanReq) (*zqbapis.ScanPayScanResp, error) {
	// 檢查與處理參數
	if err := hd.checkReq(req); err != nil {
		return nil, err
	}

	// 解析QRCode
	payCodeScanPay, err := hd.parsingContent(req.QrCode)
	if err != nil {
		return nil, err
	}

	respData := &zqbapis.ScanPayScanResp{
		Id:          fmt.Sprintf("%v", payCodeScanPay.ID),
		Amount:      tools.DecimalToFloat64(payCodeScanPay.Amount),
		ExpiredTime: payCodeScanPay.Expired.Format(time.RFC3339),
	}

	return respData, nil
}

func (hd *scanPayScan) checkReq(req *zqbapis.ScanPayScanReq) error {
	if req.QrCode == "" {
		return errs.CommonRequestParamParseFailed
	}
	return nil
}

func (hd *scanPayScan) parsingContent(qrCode string) (*model.PayCodeScanPay, error) {
	scanPay := &model.PayCodeScanPay{}
	decodeQrCodeString, err := packet.PayCode.Decode(qrCode, scanPay)
	if err != nil {
		return nil, err
	}
	return decodeQrCodeString, nil
}
