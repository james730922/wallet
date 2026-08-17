package deposit

import (
	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/db"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newBonusDAO() IBonusDAO {
	return &bonusDAO{}
}

type IBonusDAO interface {
	First(dc *gorm.DB, cond *condition.OrderBonusQuery) (*model.OrderBonus, error)
	List(dc *gorm.DB, cond *condition.OrderBonusQuery) ([]*model.OrderBonus, error)
}

func (dao *bonusDAO) List(dc *gorm.DB, cond *condition.OrderBonusQuery) ([]*model.OrderBonus, error) {
	query := condition.NewQuery(cond)
	result := make([]*model.OrderBonus, 0)

	if err := dc.
		Scopes(db.ParseWhere(query.Where())).
		Find(&result).
		Error; err != nil {
		return nil, errs.ConvertDB(err)
	}

	return result, nil
}

type bonusDAO struct {
	model.OrderBonus
}

func (dao *bonusDAO) First(dc *gorm.DB, cond *condition.OrderBonusQuery) (*model.OrderBonus, error) {
	query := condition.NewQuery(cond)
	var result model.OrderBonus

	if err := dc.
		Scopes(db.ParseWhere(query.Where())).
		First(&result).
		Error; err != nil {
		return nil, errs.ConvertDB(err)
	}

	return &result, nil
}
