package transaction

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"

	"github.com/james730922/wallet/service/internal/models"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
)

func newCommon() ICommon {
	return &commonUseCase{}
}

type ICommon interface {
	GetOneByID(dc *gorm.DB, id int64) (*model.Transaction, error)
	ListWithPaging(dc *gorm.DB, query *condition.TransactionQuery) ([]*model.Transaction, *models.PagingResult, error)
	SumWithoutPaging(dc *gorm.DB, query *condition.TransactionQuery) (decimal.Decimal, error)
	GetLatestRecord(dc *gorm.DB, memberID int64, dateTime time.Time) (*model.Transaction, error)

	insertTransaction(dc *gorm.DB, ctx context.Context, m *model.Transaction) error
	ListBySourceType(dc *gorm.DB, query *condition.TransactionQuery) ([]*model.Transaction, error)
	UpdateTransaction(dc *gorm.DB, cond *condition.Update) error
}

type commonUseCase struct {
}

func (uc *commonUseCase) GetOneByID(dc *gorm.DB, id int64) (*model.Transaction, error) {
	return dao.Transaction.First(dc, id)
}

func (uc *commonUseCase) ListWithPaging(dc *gorm.DB, query *condition.TransactionQuery) ([]*model.Transaction, *models.PagingResult, error) {
	return dao.Transaction.List(dc, query)
}

func (uc *commonUseCase) SumWithoutPaging(dc *gorm.DB, query *condition.TransactionQuery) (decimal.Decimal, error) {
	return dao.Transaction.SumWithoutPaging(dc, query)
}

func (uc *commonUseCase) GetLatestRecord(dc *gorm.DB, memberID int64, dateTime time.Time) (*model.Transaction, error) {
	return dao.Transaction.GetLatestRecord(dc, memberID, dateTime)
}

func (uc *commonUseCase) insertTransaction(dc *gorm.DB, ctx context.Context, m *model.Transaction) error {
	return dao.Transaction.InsertTransaction(dc, ctx, m)
}

func (uc *commonUseCase) ListBySourceType(dc *gorm.DB, query *condition.TransactionQuery) ([]*model.Transaction, error) {
	return dao.Transaction.ListBySourceType(dc, query)
}

func (uc *commonUseCase) UpdateTransaction(dc *gorm.DB, cond *condition.Update) error {
	return dao.Transaction.UpdateTransaction(dc, cond)
}
