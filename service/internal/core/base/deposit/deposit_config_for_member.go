package deposit

import (
	"context"

	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

type depositConfigForMember struct {
}

func (hd *depositConfigForMember) Handler(memberID int64) (*model.DepositConfigForMember, error) {

	// 取得全站佣金設定
	baseConf, err := self.DepositConfig.GetDepositBaseConfig(context.TODO())
	if err != nil {
		logger.ApLog().Error(err)
		return nil, err
	}

	result := &model.DepositConfigForMember{
		MemberID: memberID,
		Bonus:    baseConf.Bonus,
	}

	// 取得會員層級佣金設定
	member, err := packet.Member.Get(context.TODO(), memberID)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, err
	}

	if memberLevelConf, err := self.DepositConfig.GetDepositMemberLevelConfig(&condition.DepositConfigMemberLevelQuery{
		MemberLevels: &[]int64{member.LevelCode},
		Statuses:     &[]int{int(model.DepositConfigStatusClosed), int(model.DepositConfigStatusEnabled)},
	}); err != nil {
		if err != errs.DBNoRow {
			logger.ApLog().Error(err)
			return nil, err
		}
	} else if memberLevelConf.Status == model.DepositConfigStatusEnabled {
		if memberLevelConf.Bonus != nil {
			result.Bonus = *memberLevelConf.Bonus
		}
	}

	// 取得會員個別佣金設定
	if memberConf, err := self.DepositConfig.GetDepositMemberConfig(memberID); err != nil {
		if err != errs.DBNoRow {
			logger.ApLog().Error(err)
			return nil, err
		}
	} else if memberConf.Status == model.DepositConfigStatusEnabled {
		if memberConf.Bonus != nil {
			result.Bonus = *memberConf.Bonus
		}
	}

	return result, nil
}
