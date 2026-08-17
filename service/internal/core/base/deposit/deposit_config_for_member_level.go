package deposit

import (
	"context"

	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

type depositConfigForMemberLevel struct {
}

func (hd *depositConfigForMemberLevel) Handler(memberLevel int64) (*model.DepositConfigForMemberLevel, error) {
	// 取得全站佣金設定
	baseConf, err := self.DepositConfig.GetDepositBaseConfig(context.TODO())
	if err != nil {
		logger.ApLog().Error(err)
		return nil, err
	}

	result := &model.DepositConfigForMemberLevel{
		MemberLevel: memberLevel,
		Bonus:       baseConf.Bonus,
	}

	// 取得層級佣金設定
	if memberLevelConf, err := self.DepositConfig.GetDepositMemberLevelConfig(&condition.DepositConfigMemberLevelQuery{
		MemberLevels: &[]int64{memberLevel},
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

	return result, nil
}
