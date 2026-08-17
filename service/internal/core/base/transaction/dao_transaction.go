package transaction

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"

	"github.com/james730922/wallet/service/internal/models"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/db"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newTransactionDAO() ITransactionDAO {
	return &transactionDAO{}
}

type ITransactionDAO interface {
	InsertTransaction(dc *gorm.DB, ctx context.Context, cond *model.Transaction) error
	List(dc *gorm.DB, cond *condition.TransactionQuery) ([]*model.Transaction, *models.PagingResult, error)
	First(dc *gorm.DB, id int64) (*model.Transaction, error)
	FirstBySource(dc *gorm.DB, sourceType model.TransactionSourceType, sourceID int64) (*model.Transaction, error)
	SumWithoutPaging(dc *gorm.DB, cond *condition.TransactionQuery) (decimal.Decimal, error)
	GetLatestRecord(dc *gorm.DB, memberID int64, dateTime time.Time) (*model.Transaction, error)
	UpdateTransaction(dc *gorm.DB, cond *condition.Update) error
	ListBySourceType(dc *gorm.DB, cond *condition.TransactionQuery) ([]*model.Transaction, error)
}

type transactionDAO struct {
}

func (dao *transactionDAO) InsertTransaction(dc *gorm.DB, _ context.Context, cond *model.Transaction) error {
	if err := dc.Create(cond).Error; err != nil {
		return errs.ConvertDB(err)
	}

	return nil
}

func (dao *transactionDAO) UpdateTransaction(dc *gorm.DB, cond *condition.Update) error {
	if err := dc.Model(model.Transaction{}).
		Scopes(db.ParseWhere(cond.Where())).
		Updates(cond.Update()).
		Error; err != nil {
		return errs.ConvertDB(err)
	}
	return nil
}

func (dao *transactionDAO) SumWithoutPaging(dc *gorm.DB, cond *condition.TransactionQuery) (decimal.Decimal, error) {
	type sumResult struct {
		Total decimal.Decimal
	}

	result := sumResult{}

	query := dc.Table(model.Transaction{}.TableName()).Select("COALESCE(SUM(transaction.amount), 0) as total")

	if cond.CountryCode != nil || cond.Mobile != nil || cond.MemberName != nil {
		query = query.Joins("INNER JOIN member_mapping ON transaction.member_id = member_mapping.id")
	}

	iQuery := condition.NewQuery(cond)

	if err := query.
		Scopes(db.ParseWhere(iQuery.Where())).
		Find(&result).
		Error; err != nil {
		logger.ApLog().Error(err)
		return decimal.Zero, errs.ConvertDB(err)
	}

	return result.Total, nil
}

func (dao *transactionDAO) List(dc *gorm.DB, cond *condition.TransactionQuery) ([]*model.Transaction, *models.PagingResult, error) {
	query := dc.Table(model.Transaction{}.TableName()).Select("transaction.*")

	if cond.CountryCode != nil || cond.Mobile != nil || cond.MemberName != nil {
		query = query.Joins("INNER JOIN member_mapping ON transaction.member_id = member_mapping.id")
	}

	iQuery := condition.NewQuery(cond)
	result := make([]*model.Transaction, 0, iQuery.Paging().GetSize())

	if err := query.
		Scopes(db.ParseWhere(iQuery.Where())).
		Scopes(db.ParsePaging(iQuery.Paging())).
		Order("id DESC").
		Find(&result).
		Error; err != nil {
		logger.ApLog().Error(err)
		return nil, nil, errs.ConvertDB(err)
	}

	var count int
	if err := query.
		Scopes(db.ParseWhere(iQuery.Where())).
		Count(&count).
		Error; err != nil {
		logger.ApLog().Error(err)
		return nil, nil, errs.ConvertDB(err)
	}

	return result, models.NewPagingResult(iQuery.Paging(), count), nil
}

func (dao *transactionDAO) ListBySourceType(dc *gorm.DB, cond *condition.TransactionQuery) ([]*model.Transaction, error) {
	query := dc.Table(model.Transaction{}.TableName())
	iQuery := condition.NewQuery(cond)
	var result []*model.Transaction

	if err := query.
		Scopes(db.ParseWhere(iQuery.Where())).
		Find(&result).
		Error; err != nil {
		logger.ApLog().Error(err)
		return nil, errs.ConvertDB(err)
	}

	return result, nil
}

func (dao *transactionDAO) First(dc *gorm.DB, id int64) (*model.Transaction, error) {
	var result model.Transaction

	if err := dc.Where("id = ?", id).First(&result).Error; err != nil {
		return nil, errs.ConvertDB(err)
	}

	return &result, nil
}

func (dao *transactionDAO) FirstBySource(dc *gorm.DB, sourceType model.TransactionSourceType, sourceID int64) (*model.Transaction, error) {
	var result model.Transaction

	if err := dc.
		Where("source_type = ? AND source_id = ?", sourceType, sourceID).
		First(&result).Error; err != nil {
		return nil, errs.ConvertDB(err)
	}

	return &result, nil
}

func (dao *transactionDAO) GetLatestRecord(dc *gorm.DB, memberID int64, dateTime time.Time) (*model.Transaction, error) {
	result := &model.Transaction{}

	sql := `SELECT 
    			*
			FROM
    			transaction
			WHERE
    			added_time <= ?
        	AND 
				member_id = ?`

	if err := dc.
		Raw(sql, dateTime, memberID).
		Order("added_time desc").
		First(result).
		Error; err != nil {
		return nil, errs.ConvertDB(err)
	}

	return result, nil
}
