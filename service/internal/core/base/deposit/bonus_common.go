package deposit

import (
	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
)

func newBonusCommon() IBonusCommon {
	return &bonusCommonUseCase{}
}

type IBonusCommon interface {
	SelectDepositBonus(dc *gorm.DB, cond *condition.OrderBonusQuery) (*model.OrderBonus, error)
	ListDepositBonus(dc *gorm.DB, cond *condition.OrderBonusQuery) ([]*model.OrderBonus, error)
}

type bonusCommonUseCase struct{}

func (uc *bonusCommonUseCase) SelectDepositBonus(dc *gorm.DB, cond *condition.OrderBonusQuery) (*model.OrderBonus, error) {
	return dao.Bonus.First(dc, cond)
}

func (uc *bonusCommonUseCase) ListDepositBonus(dc *gorm.DB, cond *condition.OrderBonusQuery) ([]*model.OrderBonus, error) {
	return dao.Bonus.List(dc, cond)
}
