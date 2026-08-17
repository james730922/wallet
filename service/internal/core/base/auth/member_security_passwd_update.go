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
	"github.com/james730922/wallet/service/internal/utils/tools"
)

func newSecurityPasswdUpdate() *securityPasswdUpdate {
	return &securityPasswdUpdate{}
}

type securityPasswdUpdate struct{}

func (hd *securityPasswdUpdate) Handler(ctx context.Context, req *zqbapis.SecurityPasswdUpdateReq) (*zqbapis.SecurityPasswdUpdateResp, error) {
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

	updateMapping := true
	switch req.Type {
	// 确认原参数后更新
	case zqbapis.SecurityPasswdUpdateType_UpdatePwd:
		// 驗證原密碼 > 修改
		if err := hd.checkOriginPasswd(ctx, memberInfo, req); err != nil {
			return nil, err
		}
	case zqbapis.SecurityPasswdUpdateType_FirstToAdd:
		// 只允許尚未設定安全密碼的會員執行，條件式 UPDATE 可避免併發繞過。
		if err := hd.validateFirstToAdd(memberInfo); err != nil {
			return nil, err
		}
		hash, err := password.Hash(req.NewPasswd)
		if err != nil {
			return nil, errs.DBUpdateFailed
		}
		updated, err := packet.Member.SetInitialSecurityPasswd(
			ctx,
			memberInfo.ID,
			hash,
		)
		if err != nil {
			logger.ApLog().Error(err)
			return nil, errs.DBUpdateFailed
		}
		if !updated {
			return nil, errs.AuthOperationForbidden
		}
		updateMapping = false
	case zqbapis.SecurityPasswdUpdateType_ForgetToUpdate:
		// 驗證 token > 修改
		tokenMember, err := self.Token.AuthMemberForgetSecurityPasswd(ctx, req.Token)
		if err != nil {
			return nil, err
		}
		if tokenMember.ID != memberInfo.ID {
			return nil, errs.AuthOperationForbidden
		}
		// 先消耗一次性 token，避免相同 token 重複重設。
		if err := self.Token.DeleteMemberForgetSecurityPasswdByID(ctx, memberInfo.ID); err != nil {
			return nil, errs.AuthTokenDelFailed
		}
		if err := hd.forceUpdatePasswd(memberInfo, req); err != nil {
			return nil, err
		}
	default:
		return nil, errs.CommonUnknownError
	}

	if updateMapping {
		if err := packet.Member.UpdateMapping(ctx, memberInfo); err != nil {
			logger.ApLog().Errorf("Update member_mapping err: %v", err)
			return nil, errs.CommonNoData
		}
	}

	return &zqbapis.SecurityPasswdUpdateResp{
		Status: models.CommonStatusSuccess,
	}, nil
}

func (hd *securityPasswdUpdate) validateFirstToAdd(memberInfo *model.MemberMapping) error {
	if memberInfo.SecurityPasswd != "" {
		return errs.AuthOperationForbidden
	}
	return nil
}

func (hd *securityPasswdUpdate) validateReq(req *zqbapis.SecurityPasswdUpdateReq) error {
	// 判斷是否有輸入資料 => 这边不判断req.Passwd 因为只有部分方式用到，后续才会再验证
	if req.NewPasswd == "" {
		return errs.NewPasswdEmpty
	}

	if req.NewPasswd == "" {
		return errs.NewConfirmPasswdEmpty
	}

	if err := tools.CheckSecurityPasswdRule(req.NewPasswd); err != nil {
		return err
	}

	return nil
}

func (hd *securityPasswdUpdate) checkOriginPasswd(ctx context.Context, memberInfo *model.MemberMapping, req *zqbapis.SecurityPasswdUpdateReq) error {
	if req.Passwd == "" {
		return errs.PasswdEmpty
	}

	valid, _, err := password.Verify(memberInfo.SecurityPasswd, req.Passwd, memberInfo.Salt)
	if err != nil || !valid {
		return errs.RegistrationCheckOldPasswdError
	}

	memberInfo.SecurityPasswd, err = password.Hash(req.NewPasswd)
	if err != nil {
		return errs.DBUpdateFailed
	}

	return nil
}

func (hd *securityPasswdUpdate) forceUpdatePasswd(memberInfo *model.MemberMapping, req *zqbapis.SecurityPasswdUpdateReq) error {
	hash, err := password.Hash(req.NewPasswd)
	if err != nil {
		return errs.DBUpdateFailed
	}
	memberInfo.SecurityPasswd = hash
	return nil
}
