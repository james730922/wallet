package auth

import (
	"context"

	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/pb/zqbapis"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/ctxs"
	"github.com/james730922/wallet/service/internal/utils/errs"
	"github.com/james730922/wallet/service/internal/utils/password"
	"github.com/james730922/wallet/service/internal/utils/tools"
)

func newUpdatePasswd() *updatePasswd {
	return &updatePasswd{}
}

type updatePasswd struct{}

// 使用目前密碼修改登入密碼
func (hd *updatePasswd) Handler(ctx context.Context, req *zqbapis.LoginUpdatePasswdReq) (*zqbapis.LoginUpdatePasswdResp, error) {
	var token string

	if err := hd.validateReq(req); err != nil {
		logger.ApLog().Error(err)
		return nil, err
	}
	authenticatedMemberID, ok := ctxs.GetMemberID(ctx)
	if !ok {
		return nil, errs.CommonNoMemberID
	}

	// 確認會員存在
	mapping, err := hd.validateMember(ctx, authenticatedMemberID, req.CountryCode, req.Mobile, req.Passwd)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, err
	}

	// 更新密碼
	if err := hd.updatePasswd(ctx, mapping, req.NewPasswd); err != nil {
		logger.ApLog().Error(err)
		return nil, err
	}

	// 密碼異動後一律輪替 token，舊工作階段立即失效。
	token, err = self.Token.GenMember(ctx, mapping.ID)
	if err != nil {
		// 密碼已更新時不可繼續保留舊工作階段；Redis 異常時盡力撤銷。
		_ = self.Token.DeleteMemberByID(ctx, mapping.ID)
		logger.ApLog().Error(err)
		return nil, errs.AuthTokenGetFailed
	}

	resp := &zqbapis.LoginUpdatePasswdResp{
		Token:       token,
		Mobile:      req.Mobile,
		CountryCode: req.CountryCode,
	}

	return resp, nil
}

// 驗證請求是否正確
func (hd *updatePasswd) validateReq(req *zqbapis.LoginUpdatePasswdReq) error {

	// 判斷是否有輸入資料
	if req.Mobile == "" || req.CountryCode == "" {
		return errs.MobileEmpty
	}

	if req.Passwd == "" {
		return errs.PasswdEmpty
	}

	if req.NewPasswd == "" {
		return errs.NewPasswdEmpty
	}

	if req.ConfirmNewPasswd == "" {
		return errs.NewConfirmPasswdEmpty
	}

	// 驗證兩次密碼輸入是否相同
	if req.NewPasswd != req.ConfirmNewPasswd {
		return errs.PasswdNotEqual
	}

	// 驗證密碼格式
	if err := tools.MemberProfileValidate.Passwd(req.NewPasswd); err != nil {
		if err != errs.RegistrationPasswdFmtErr {
			logger.ApLog().Warnf("password validation failed: %v", err)
			return errs.CommonRequestParamInvalid
		}
		return err
	}

	return nil
}

// 驗證會員資格
func (hd *updatePasswd) validateMember(ctx context.Context, authenticatedMemberID int64, countryCode, mobile, passwd string) (*model.MemberMapping, error) {
	mapping, err := packet.Member.FindMappingWithMobile(ctx, countryCode, mobile)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, errs.MemberNotFound
	}
	if err := validatePasswordChangeOwnership(authenticatedMemberID, mapping.ID); err != nil {
		return nil, err
	}

	valid, _, err := password.Verify(mapping.Passwd, passwd, mapping.Salt)
	if err != nil || !valid {
		return nil, errs.AuthAccountOrPasswdError
	}

	return mapping, nil
}

func validatePasswordChangeOwnership(authenticatedMemberID, targetMemberID int64) error {
	if authenticatedMemberID != targetMemberID {
		return errs.AuthOperationForbidden
	}
	return nil
}

// 產生密碼和更新資料
func (hd *updatePasswd) updatePasswd(ctx context.Context, mapping *model.MemberMapping, passwd string) error {
	hash, err := password.Hash(passwd)
	if err != nil {
		return errs.DBUpdateFailed
	}
	mapping.Passwd = hash

	if err := packet.Member.UpdateMapping(ctx, mapping); err != nil {
		logger.ApLog().Error(err)
		return errs.DBUpdateFailed
	}

	return nil
}
