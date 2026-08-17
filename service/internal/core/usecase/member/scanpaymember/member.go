package scanpaymember

import (
	"context"
	"github.com/james730922/wallet/service/internal/pb/zqbapis"
)

func newScanPayMember() IMember {
	return &scanPayMemberUseCase{}
}

type IMember interface {
	Scan(ctx context.Context, req *zqbapis.ScanPayScanReq) (*zqbapis.ScanPayScanResp, error)
	Pay(ctx context.Context, req *zqbapis.ScanPayAddPayReq) (*zqbapis.ScanPayAddPayResp, error)
	UploadQRCode(ctx context.Context, req *zqbapis.ScanPayScanReq) (*zqbapis.ScanPayScanResp, error)
}

type scanPayMemberUseCase struct {
	scan   *scanPayScan
	pay    *scanPayPay
	qrcode *scanPayQRCode
}

func (uc scanPayMemberUseCase) Scan(ctx context.Context, req *zqbapis.ScanPayScanReq) (*zqbapis.ScanPayScanResp, error) {
	return uc.scan.Handler(ctx, req)
}

func (uc scanPayMemberUseCase) Pay(ctx context.Context, req *zqbapis.ScanPayAddPayReq) (*zqbapis.ScanPayAddPayResp, error) {
	return uc.pay.Handler(ctx, req)
}

func (uc scanPayMemberUseCase) UploadQRCode(ctx context.Context, req *zqbapis.ScanPayScanReq) (*zqbapis.ScanPayScanResp, error) {
	return uc.qrcode.Handler(ctx, req)
}
