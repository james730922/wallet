package user

import (
	"context"

	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/pb/zqbapis"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/ctxs"
	"github.com/james730922/wallet/service/internal/utils/errs"
	"github.com/james730922/wallet/service/internal/utils/password"
)

func newMemberCheckScanPayPasswd() *checkMemberScanPayPasswd {
	return &checkMemberScanPayPasswd{}
}

type checkMemberScanPayPasswd struct{}

func (hd *checkMemberScanPayPasswd) Handler(ctx context.Context, pwd string) error {
	memberID, ok := ctxs.GetMemberID(ctx)
	if !ok {
		logger.ApLog().Error(errs.CommonNoMemberID)
		return errs.CommonNoMemberID
	}

	memberInfo, err := self.Member.FindMobile(ctx, memberID)
	if err != nil {
		logger.ApLog().Error(err)
		return errs.MemberNotFound
	}

	// 验证密码
	if err := hd.checkScanPayPasswd(ctx, memberInfo, pwd); err != nil {
		return err
	}

	return nil
}

func (hd *checkMemberScanPayPasswd) validateReq(req *zqbapis.CheckScanPayPasswdReq) error {
	// 判斷是否有輸入資料
	if req.Passwd == "" {
		return errs.PasswdEmpty
	}

	return nil
}

func (hd *checkMemberScanPayPasswd) checkScanPayPasswd(ctx context.Context, memberInfo *model.MemberMapping, passwdText string) error {
	valid, needsUpgrade, err := password.Verify(memberInfo.SecurityPasswd, passwdText, memberInfo.Salt)
	if err != nil || !valid {
		return errs.RegistrationCheckScanPayPasswdError
	}
	if needsUpgrade {
		hash, err := password.Hash(passwdText)
		if err != nil {
			return errs.DBUpdateFailed
		}
		memberInfo.SecurityPasswd = hash
		if err := self.Member.UpdateMapping(ctx, memberInfo); err != nil {
			return errs.DBUpdateFailed
		}
	}

	return nil
}
