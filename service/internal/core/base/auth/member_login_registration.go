package auth

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/pb/zqbapis"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/errs"
	"github.com/james730922/wallet/service/internal/utils/tools"
)

func newLoginMemberRegistration() *loginMemberRegistration {
	return &loginMemberRegistration{}
}

type loginMemberRegistration struct{}

// 註冊會員
func (hd *loginMemberRegistration) Handler(ctx context.Context, req *zqbapis.LoginRegistrationReq) (*zqbapis.LoginRegistrationResp, error) {
	if err := self.registrationGuard.beforeAttempt(ctx); err != nil {
		return nil, err
	}

	if err := packet.GeeTestCaptcha.Validate(req.Challenge, req.Validate, req.SecCode); err != nil {
		return nil, err
	}

	if err := hd.validateReq(req); err != nil {
		return nil, err
	}

	// 確認會員存在
	exist, err := hd.isMemberExist(ctx, req.CountryCode, req.Mobile)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, errs.LoginMemberRegistFailed
	}

	if exist {
		return nil, errs.LoginMemberRegistDuplicated
	}

	// 創會員
	member, err := hd.createMember(ctx, req)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, errs.MemberCreateFailed
	}

	// 發token
	token, err := hd.genToken(ctx, member)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, errs.AuthTokenCreateFailed
	}

	member.LastLoginTime = aws.Time(time.Now().UTC())
	memberUpdateReq := condition.MemberUpdateCond{
		ID:            &member.ID,
		LastLoginTime: member.LastLoginTime,
	}
	if err := packet.Member.Update(ctx, &memberUpdateReq); err != nil {
		logger.ApLog().Errorf("req:%v, err:%v", tools.JsonMarshalString(memberUpdateReq), err)
		return nil, errs.DBUpdateFailed
	}

	// 暫存表
	self.memberRepository.DeleteByID(member.ID)
	self.memberRepository.Store(token, member)

	resp := &zqbapis.LoginRegistrationResp{
		Token:       token,
		Mobile:      req.Mobile,
		CountryCode: req.CountryCode,
	}

	return resp, nil
}

// 驗證請求是否正確
func (hd *loginMemberRegistration) validateReq(req *zqbapis.LoginRegistrationReq) error {
	// 判斷是否有輸入資料
	if req.Mobile == "" || req.CountryCode == "" {
		return errs.MobileEmpty
	}

	if req.Passwd == "" {
		return errs.PasswdEmpty
	}

	if req.ConfirmPasswd == "" {
		return errs.RegistrationConfirmPasswdEmpty
	}

	if req.QqAcount == "" {
		return errs.QQEmpty
	}

	if req.Name == "" {
		return errs.NameEmpty
	}

	// 驗證兩次密碼輸入是否相同
	if req.Passwd != req.ConfirmPasswd {
		return errs.PasswdNotEqual
	}

	// 驗證密碼格式
	if err := tools.MemberProfileValidate.Passwd(req.Passwd); err != nil {
		if err != errs.RegistrationPasswdFmtErr {
			logger.ApLog().Error(err)
			return errs.CommonRequestParamInvalid
		}
		return err
	}

	// 驗證手機格式
	if err := tools.MemberProfileValidate.Mobile(req.CountryCode, req.Mobile); err != nil {
		if err != errs.RegistrationMobileFmtErr {
			logger.ApLog().Error(err)
			return errs.CommonRequestParamInvalid
		}
		return err
	}

	// 驗證手機是否重複
	if _, err := packet.Member.FindMappingWithMobile(context.TODO(), req.CountryCode, req.Mobile); err == nil {
		return errs.MobileDuplicate
	} else {
		if err != errs.DBNoRow {
			logger.ApLog().Error(err)
			return err
		}
	}

	// 驗證QQ
	if err := tools.MemberProfileValidate.QQ(req.QqAcount); err != nil {
		if err != errs.RegistrationQQFmtErr {
			logger.ApLog().Error(err)
			return errs.CommonRequestParamInvalid
		}
		return err
	}

	// 驗證qq是否重複
	if _, err := packet.Member.FindMappingWithQQ(context.TODO(), req.QqAcount); err == nil {
		return errs.QQDuplicate
	} else {
		if err != errs.DBNoRow {
			logger.ApLog().Error(err)
			return err
		}
	}

	// 驗證中文名
	if err := tools.MemberProfileValidate.Name(req.Name); err != nil {
		if err != errs.RegistrationNameFmtErr {
			logger.ApLog().Error(err)
			return errs.CommonRequestParamInvalid
		}
		return err
	}

	return nil
}

// 驗證會員是否已經存在
func (hd *loginMemberRegistration) isMemberExist(ctx context.Context, countryCode, mobile string) (exist bool, err error) {
	_, err = packet.Member.Find(ctx, countryCode, mobile)
	if err == errs.DBNoRow {
		err = nil
		return
	}

	if err != nil {
		logger.ApLog().Error(err)
		return
	}

	exist = true
	return
}

// 建立新會員
func (hd *loginMemberRegistration) createMember(ctx context.Context, req *zqbapis.LoginRegistrationReq) (*model.Member, error) {
	member, err := packet.Member.CreateWithPasswd(ctx, req.CountryCode, req.Mobile, req.Passwd, req.QqAcount, req.Name)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, err
	}
	return member, nil
}

// 產token
func (hd *loginMemberRegistration) genToken(ctx context.Context, member *model.Member) (string, error) {
	token, err := token.GenMember(ctx, member.ID)
	if err != nil {
		logger.ApLog().Errorf("gen token failure, err: %v", err)
		return token, err
	}
	return token, err
}
