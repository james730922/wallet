package deposit

import (
	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/models"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/db"
)

func newCommonDAO() ICommonDAO {
	return &commonDAO{}
}

type ICommonDAO interface {
	List(dc *gorm.DB, cond condition.IQuery) ([]*model.Deposit, *models.PagingResult, error)
	Insert(dc *gorm.DB, cond *model.Deposit) error
}

type commonDAO struct{}

func (dao *commonDAO) List(dc *gorm.DB, cond condition.IQuery) ([]*model.Deposit, *models.PagingResult, error) {
	result := make([]*model.Deposit, 0, cond.Paging().Size)

	query := dc.Model(result).Scopes(db.ParseWhere(cond.Where()))

	if err := query.
		Scopes(db.ParsePaging(cond.Paging())).
		Order("id DESC").
		Find(&result).GetErrors(); len(err) > 0 {
		return nil, nil, err[0]
	}

	var count int
	if err := query.Count(&count).GetErrors(); len(err) > 0 {
		return nil, nil, err[0]
	}

	return result, models.NewPagingResult(cond.Paging(), count), nil
}

func (dao *commonDAO) Insert(dc *gorm.DB, cond *model.Deposit) error {
	if err := dc.Create(cond).Error; err != nil {
		return err
	}

	return nil
}
