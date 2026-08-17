package auth

import (
	"context"

	"github.com/james730922/wallet/service/internal/models"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/pb/zqbapis"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/ctxs"
	"github.com/james730922/wallet/service/internal/utils/errs"
	"github.com/james730922/wallet/service/internal/utils/password"
)

func newSecurityPasswdForget() *securityPasswdForget {
	return &securityPasswdForget{}
}

type securityPasswdForget struct{}

func (hd *securityPasswdForget) Handler(ctx context.Context, req *zqbapis.SecurityPasswdForgetReq) (*zqbapis.SecurityPasswdForgetResp, error) {
	if err := hd.validateReq(req); err != nil {
		logger.ApLog().Error(err)
		return nil, err
	}

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

	if err := hd.validateMember(ctx, req, memberInfo); err != nil {
		logger.ApLog().Error(err)
		return nil, err
	}

	// 刪token
	if err := self.Token.DeleteMemberForgetSecurityPasswdByID(ctx, memberID); err != nil {
		logger.ApLog().Error(err)
		return nil, errs.AuthTokenDelFailed
	}

	// 發token
	token, err := self.Token.GenMemberForgetSecurityPasswd(ctx, memberID)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, errs.AuthTokenCreateFailed
	}

	return &zqbapis.SecurityPasswdForgetResp{
		Status: models.CommonStatusSuccess,
		Token:  token,
	}, nil
}

func (hd *securityPasswdForget) validateReq(req *zqbapis.SecurityPasswdForgetReq) error {
	// 判斷是否有輸入資料
	if req.Mobile == "" {
		return errs.MobileEmpty
	}

	if req.QqAccount == "" {
		return errs.QQEmpty
	}

	if req.Name == "" {
		return errs.NameEmpty
	}

	if req.LoginPasswd == "" {
		return errs.PasswdEmpty
	}
	return nil
}

func (hd *securityPasswdForget) validateMember(ctx context.Context, req *zqbapis.SecurityPasswdForgetReq, memberInfo *model.MemberMapping) error {
	if req.QqAccount != memberInfo.QQ {
		return errs.DataNotQualified
	}

	if req.Name != memberInfo.Name {
		return errs.DataNotQualified
	}

	if req.Mobile != memberInfo.Mobile {
		return errs.DataNotQualified
	}

	valid, needsUpgrade, err := password.Verify(memberInfo.Passwd, req.LoginPasswd, memberInfo.Salt)
	if err != nil || !valid {
		return errs.DataNotQualified
	}
	if needsUpgrade {
		memberInfo.Passwd, err = password.Hash(req.LoginPasswd)
		if err != nil {
			return errs.DBUpdateFailed
		}
		if err := packet.Member.UpdateMapping(ctx, memberInfo); err != nil {
			return errs.DBUpdateFailed
		}
	}

	return nil
}
