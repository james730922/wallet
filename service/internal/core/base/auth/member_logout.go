package auth

import (
	"context"

	"github.com/james730922/wallet/service/internal/pb/zqbapis"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/ctxs"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newLoginMemberLogout() *loginMemberLogout {
	return &loginMemberLogout{}
}

type loginMemberLogout struct{}

// 登出
func (hd *loginMemberLogout) Handler(ctx context.Context) (*zqbapis.LogoutResp, error) {
	resp := &zqbapis.LogoutResp{
		Status: zqbapis.LogoutStatus_Failure,
	}
	memberID, ok := ctxs.GetMemberID(ctx)
	if !ok {
		logger.ApLog().Error(errs.CommonNoMemberID)
		return resp, errs.CommonNoMemberID
	}

	// 取得會員
	member, err := packet.Member.Get(ctx, memberID)
	if err != nil {
		logger.ApLog().Error(err)
		return resp, errs.MemberNotFound
	}

	// 刪token
	if err := self.Token.DeleteMemberByID(ctx, memberID); err != nil {
		logger.ApLog().Error(err)
		return resp, errs.AuthTokenDelFailed
	}

	// 暫存表
	self.memberRepository.DeleteByID(member.ID)
	resp.Status = zqbapis.LogoutStatus_Success

	return resp, nil
}
