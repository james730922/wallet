package auth

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/pb/zqbapis"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/ctxs"
	"github.com/james730922/wallet/service/internal/utils/errs"
	"github.com/james730922/wallet/service/internal/utils/password"
	"github.com/james730922/wallet/service/internal/utils/tools"
)

func newLoginMemberWithPasswd() *loginMemberWithPasswd {
	return &loginMemberWithPasswd{}
}

type loginMemberWithPasswd struct{}

// 登陸會員
func (hd *loginMemberWithPasswd) Handler(ctx context.Context, req *zqbapis.LoginIdentifierPasswdReq) (*zqbapis.LoginIdentifierPasswdResp, error) {
	if err := hd.validateReq(req); err != nil {
		logger.ApLog().Error(err)
		return nil, err
	}

	clientIP, _ := ctxs.GetClientIP(ctx)
	requireCaptcha, err := self.loginGuard.beforeAttempt(req.CountryCode, req.Mobile, clientIP)
	if err != nil {
		return nil, err
	}
	if requireCaptcha {
		if !packet.GeeTestCaptcha.Enabled() {
			return nil, errs.CommonFrequentOperationError
		}
		if err := packet.GeeTestCaptcha.Validate(req.Challenge, req.Validate, req.SecCode); err != nil {
			return nil, err
		}
	}

	// 確認會員存在
	member, failedMemberID, err := hd.validateMember(ctx, req.CountryCode, req.Mobile, req.Passwd)
	if err != nil {
		hd.recordFailedAttempt(ctx, failedMemberID)
		logger.ApLog().Debugf("msg: %s err: %s", errs.LoginMemberLoginFailed, err)
		return nil, errs.LoginMemberLoginFailed
	}

	if member.Status != model.MemberStatusEnabled {
		return nil, errs.MemberDisabled
	}

	if err := self.loginGuard.resetAccount(req.CountryCode, req.Mobile); err != nil {
		logger.ApLog().Warnf("reset login attempt counter failed: memberID=%d err=%v", member.ID, err)
	}
	if err := packet.Member.ResetLoginFailures(ctx, member.ID); err != nil {
		logger.ApLog().Warnf("reset member failed-attempt count failed: memberID=%d err=%v", member.ID, err)
	}

	// 原子輪替 token，同一會員只保留最後一個登入階段。
	token, err := hd.genToken(ctx, member.ID)
	if err != nil {
		logger.ApLog().Errorf("memberID: %d, msg: %s err: %s", member.ID, errs.AuthTokenCreateFailed, err)
		return nil, errs.LoginMemberLoginFailed
	}

	logger.ApLog().Debugf("member login success: countryCode=%s, memberID=%d", req.CountryCode, member.ID)

	member.LastLoginTime = aws.Time(time.Now().UTC())
	memberUpdateReq := condition.MemberUpdateCond{
		ID:            &member.ID,
		LastLoginTime: member.LastLoginTime,
	}
	if err := packet.Member.Update(ctx, &memberUpdateReq); err != nil {
		logger.ApLog().Errorf("memberID: %d, msg: %s req:%v, err:%v", member.ID, errs.DBUpdateFailed, tools.JsonMarshalString(memberUpdateReq), err)
		return nil, errs.LoginMemberLoginFailed
	}
	// 暫存表
	self.memberRepository.DeleteByID(member.ID)
	self.memberRepository.Store(token, &model.Member{ID: member.ID})

	resp := &zqbapis.LoginIdentifierPasswdResp{
		Token:       token,
		Mobile:      req.Mobile,
		CountryCode: req.CountryCode,
	}

	return resp, nil
}

// 驗證請求是否正確
func (hd *loginMemberWithPasswd) validateReq(req *zqbapis.LoginIdentifierPasswdReq) error {
	// 判斷是否有輸入資料
	if req.Mobile == "" || req.CountryCode == "" {
		return errs.MobileEmpty
	}

	if req.Passwd == "" {
		return errs.PasswdEmpty
	}

	return nil
}

// 驗證會員資格
func (hd *loginMemberWithPasswd) validateMember(ctx context.Context, countryCode, mobile, passwd string) (*model.Member, int64, error) {
	mapping, err := packet.Member.FindMappingWithMobile(ctx, countryCode, mobile)
	if err != nil {
		return nil, 0, errs.MemberNotFound
	}

	valid, needsUpgrade, err := password.Verify(mapping.Passwd, passwd, mapping.Salt)
	if err != nil || !valid {
		return nil, mapping.ID, errs.AuthAccountOrPasswdError
	}
	if needsUpgrade {
		upgradedHash, err := password.Hash(passwd)
		if err != nil {
			return nil, mapping.ID, errs.LoginMemberLoginFailed
		}
		mapping.Passwd = upgradedHash
		if err := packet.Member.UpdateMapping(ctx, mapping); err != nil {
			logger.ApLog().Error(err)
			return nil, mapping.ID, errs.LoginMemberLoginFailed
		}
	}

	member, err := packet.Member.Get(ctx, mapping.ID)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, mapping.ID, errs.MemberNotFound
	}

	return member, mapping.ID, nil
}

func (hd *loginMemberWithPasswd) recordFailedAttempt(ctx context.Context, memberID int64) {
	if memberID != 0 {
		if err := packet.Member.RecordLoginFailure(ctx, memberID); err != nil {
			logger.ApLog().Warnf("record member failed-attempt count failed: memberID=%d err=%v", memberID, err)
		}
	}
}

// 產token
func (hd *loginMemberWithPasswd) genToken(ctx context.Context, id int64) (string, error) {
	token, err := token.GenMember(ctx, id)
	if err != nil {
		logger.ApLog().Errorf("gen token failure, err: %v", err)
		return token, err
	}
	return token, err
}
