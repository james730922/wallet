package auth

import (
	"context"

	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/pb/zqbapis"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/ctxs"
	"github.com/james730922/wallet/service/internal/utils/errs"
	"github.com/james730922/wallet/service/internal/utils/password"
)

func newSecurityPasswdIdentifier() *securityPasswdIdentifier {
	return &securityPasswdIdentifier{}
}

type securityPasswdIdentifier struct{}

// 驗證驗證碼
func (hd *securityPasswdIdentifier) Handler(ctx context.Context, req *zqbapis.SecurityPasswdIdentifierReq) (*zqbapis.SecurityPasswdIdentifierResp, error) {
	memberID, ok := ctxs.GetMemberID(ctx)
	if !ok {
		return nil, errs.CommonNoMemberID
	}

	memberInfo, err := packet.Member.FindMobile(ctx, memberID)
	if err != nil {
		logger.ApLog().Errorf("memberID:%d, err: %s", memberID, err)
		return nil, errs.MemberNotFound
	}

	// 验证密码
	if err := hd.validatePasswd(ctx, memberInfo, req.Passwd); err != nil {
		return nil, err
	}

	// 刪token
	if err := self.Token.DeleteMemberSecurityPasswdByID(ctx, memberID); err != nil {
		logger.ApLog().Error(err)
		return nil, errs.AuthTokenDelFailed
	}

	// 發token
	token, err := self.Token.GenMemberSecurityPasswd(ctx, memberID)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, errs.AuthTokenCreateFailed
	}

	return &zqbapis.SecurityPasswdIdentifierResp{
		Token: token,
	}, nil
}

func (hd *securityPasswdIdentifier) validatePasswd(ctx context.Context, memberInfo *model.MemberMapping, passwdText string) error {
	valid, needsUpgrade, err := password.Verify(memberInfo.SecurityPasswd, passwdText, memberInfo.Salt)
	if err != nil || !valid {
		return errs.RegistrationCheckScanPayPasswdError
	}
	if needsUpgrade {
		memberInfo.SecurityPasswd, err = password.Hash(passwdText)
		if err != nil {
			return errs.DBUpdateFailed
		}
		if err := packet.Member.UpdateMapping(ctx, memberInfo); err != nil {
			return errs.DBUpdateFailed
		}
	}

	return nil
}
