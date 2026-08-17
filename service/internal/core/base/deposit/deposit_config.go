package deposit

import (
	"context"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newDepositConfig() IDepositConfig {
	return &depositConfigUseCase{}
}

type IDepositConfig interface {
	// 取得某會員的紅利設定（已經過綜合判斷）
	GetDepositConfigForMember(memberID int64) (*model.DepositConfigForMember, error)
	// 取得某會員層級的紅利設定（已經過綜合判斷）
	GetDepositConfigForMemberLevel(memberLevel int64) (*model.DepositConfigForMemberLevel, error)

	// 取得-全站
	GetDepositBaseConfig(ctx context.Context) (*model.DepositConfig, error)

	// 取得-依據會員級別
	GetDepositMemberLevelConfig(query *condition.DepositConfigMemberLevelQuery) (*model.DepositConfigMemberLevel, error)
	// 取得-依據會員個人
	GetDepositMemberConfig(memberID int64) (*model.DepositConfigMember, error)
}

type depositConfigUseCase struct {
	forMember      *depositConfigForMember
	forMemberLevel *depositConfigForMemberLevel
	getBase        *depositConfigGetBase
	getMemberLevel *depositConfigGetMemberLevel
}

func (uc *depositConfigUseCase) GetDepositConfigForMember(memberID int64) (*model.DepositConfigForMember, error) {
	return uc.forMember.Handler(memberID)
}

func (uc *depositConfigUseCase) GetDepositConfigForMemberLevel(memberLevel int64) (*model.DepositConfigForMemberLevel, error) {
	return uc.forMemberLevel.Handler(memberLevel)
}

func (uc *depositConfigUseCase) GetDepositBaseConfig(ctx context.Context) (*model.DepositConfig, error) {
	return uc.getBase.Handler(ctx)
}

func (uc *depositConfigUseCase) GetDepositMemberLevelConfig(query *condition.DepositConfigMemberLevelQuery) (*model.DepositConfigMemberLevel, error) {
	configs, err := uc.getMemberLevel.Handler(query)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, err
	}

	// 由於此方法只返回一筆層級資訊，故在 return
	if len(configs) != 0 {
		if configs[0] != nil {
			return configs[0], nil
		}
	}
	return nil, errs.DBNoRow
}

func (uc *depositConfigUseCase) GetDepositMemberConfig(memberID int64) (*model.DepositConfigMember, error) {
	config, err := dao.ConfigMember.First(packet.DB.New(), memberID)
	if err != nil {
		return nil, errs.ConvertDB(err)
	}

	return config, nil
}
