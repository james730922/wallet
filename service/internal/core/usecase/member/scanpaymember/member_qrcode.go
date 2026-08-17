package scanpaymember

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/pb/zqbapis"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/errs"
	"github.com/james730922/wallet/service/internal/utils/tools"
)

const maxQRCodeImageBytes int64 = 4 << 20 // Base64 解碼後最多 4 MiB。

func newMemberQRCode() *scanPayQRCode {
	return &scanPayQRCode{}
}

type scanPayQRCode struct{}

func (hd *scanPayQRCode) Handler(ctx context.Context, req *zqbapis.ScanPayScanReq) (*zqbapis.ScanPayScanResp, error) {
	// 檢查與處理參數
	if err := hd.checkReq(req); err != nil {
		return nil, err
	}

	// 解析掃碼圖片
	qrString, err := hd.parseQrCode(ctx, req.QrCode)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, err
	}

	// 解析QRCode
	payCodeScanPay, err := hd.parsingContent(qrString)
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

func (hd *scanPayQRCode) checkReq(req *zqbapis.ScanPayScanReq) error {
	if req.QrCode == "" {
		return errs.CommonRequestParamParseFailed
	}
	return nil
}

func (hd *scanPayQRCode) parseQrCode(ctx context.Context, qrcode string) (string, error) {
	image, err := decodeQRCodeImage(qrcode)
	if err != nil {
		logger.ApLog().Error(err)
		return "", err
	}

	formatImg, err := tools.ImageCheck(image)
	if err != nil {
		logger.ApLog().Error(err)
		return "", errs.ScanQRCodeVerifyFailed
	}

	qrcodeContent, err := tools.IsQrCode(formatImg)
	if err != nil {
		logger.ApLog().Error(err)
		return "", errs.ScanQRCodeContentFailed
	}

	// success
	logger.ApLog().Debugf("qrcode msg = %s", qrcodeContent)
	return qrcodeContent, nil
}

func decodeQRCodeImage(encoded string) ([]byte, error) {
	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(encoded))
	image, err := ioutil.ReadAll(io.LimitReader(decoder, maxQRCodeImageBytes+1))
	if err != nil {
		return nil, errs.ScanQRCodeVerifyFailed
	}
	if int64(len(image)) > maxQRCodeImageBytes {
		return nil, errs.ScanQRCodeImageTooLarge
	}
	return image, nil
}

func (hd *scanPayQRCode) parsingContent(qrCode string) (*model.PayCodeScanPay, error) {
	scanPay := &model.PayCodeScanPay{}
	decodeQrCodeString, err := packet.PayCode.Decode(qrCode, scanPay)
	if err != nil {
		return nil, err
	}
	return decodeQrCodeString, nil
}
