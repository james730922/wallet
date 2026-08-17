package auth

import (
	"context"

	"github.com/james730922/wallet/service/internal/models"
	"github.com/james730922/wallet/service/internal/pb/zqbapis"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/ctxs"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newSecurityPasswdFirstUse() *securityPasswdFirstUse {
	return &securityPasswdFirstUse{}
}

type securityPasswdFirstUse struct{}

func (hd *securityPasswdFirstUse) Handler(ctx context.Context) (*zqbapis.SecurityPasswdFirstUseResp, error) {
	memberID, ok := ctxs.GetMemberID(ctx)
	if !ok {
		logger.ApLog().Error(errs.CommonNoMemberID)
		return nil, errs.CommonNoMemberID
	}

	memberInfo, err := packet.Member.FindMobile(ctx, memberID)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, errs.MemberNotFound
	}

	resp := &zqbapis.SecurityPasswdFirstUseResp{}
	// 验证是否第一次使用
	if ok := hd.checkFirstUse(memberInfo.SecurityPasswd); !ok {
		return resp, nil
	}

	resp.Status = models.CommonStatusSuccess
	return resp, nil
}

func (hd *securityPasswdFirstUse) checkFirstUse(passwd string) bool {
	// 查询不到密码，第一次使用
	if len(passwd) > 0 {
		return false
	}

	return true
}
