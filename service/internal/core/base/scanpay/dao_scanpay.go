package scanpay

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/jinzhu/gorm"

	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/db"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newScanPayDAO() IScanPayDAO {
	return &scanPayDAO{}
}

type IScanPayDAO interface {
	FirstForUpdate(dc *gorm.DB, cond *condition.OrderScanPayToPayCond) (*model.OrderScanPay, error)
	FirstByIDForUpdate(dc *gorm.DB, id int64) (*model.OrderScanPay, error)
	FirstByRecordID(dc *gorm.DB, recordID int64) (*model.OrderScanPay, error)
	ListStaleTransactionIDs(dc *gorm.DB, updatedBefore time.Time, limit int) ([]int64, error)
	Insert(dc *gorm.DB, cond *model.OrderScanPay) error
	Update(dc *gorm.DB, cond condition.IUpdate) error
}

func (dao *scanPayDAO) FirstByIDForUpdate(dc *gorm.DB, id int64) (*model.OrderScanPay, error) {
	var result model.OrderScanPay
	if err := dc.
		Set(db.GormSetSelectForUpdate()).
		Where("id = ?", id).
		First(&result).Error; err != nil {
		return nil, errs.ConvertDB(err)
	}

	return &result, nil
}

type scanPayDAO struct {
	model.OrderScanPay
}

func (dao *scanPayDAO) FirstForUpdate(dc *gorm.DB, cond *condition.OrderScanPayToPayCond) (*model.OrderScanPay, error) {
	var result model.OrderScanPay
	updateCond := condition.NewQuery(
		condition.OrderScanPayFirstForUpdateQuery{
			ID:       aws.Int64(cond.ID),
			MemberID: cond.MemberID,
			Amount:   cond.Amount,
		},
	)
	if err := dc.
		Set(db.GormSetSelectForUpdate()).
		Where(updateCond.Where()).
		First(&result).Error; err != nil {
		return nil, errs.ConvertDB(err)
	}

	return &result, nil
}

func (dao *scanPayDAO) FirstByRecordID(dc *gorm.DB, recordID int64) (*model.OrderScanPay, error) {
	var result model.OrderScanPay
	if err := dc.
		Where("record_id = ?", recordID).
		First(&result).Error; err != nil {
		return nil, errs.ConvertDB(err)
	}

	return &result, nil
}

func (dao *scanPayDAO) ListStaleTransactionIDs(dc *gorm.DB, updatedBefore time.Time, limit int) ([]int64, error) {
	result := make([]int64, 0, limit)
	if err := dc.
		Model(&model.OrderScanPay{}).
		Where("status = ? AND updated_time < ?", model.OrderScanPayStatusEnumTransaction, updatedBefore).
		Order("updated_time ASC").
		Limit(limit).
		Pluck("id", &result).Error; err != nil {
		return nil, errs.ConvertDB(err)
	}

	return result, nil
}

func (dao *scanPayDAO) Insert(dc *gorm.DB, cond *model.OrderScanPay) error {
	if err := dc.Create(cond).Error; err != nil {
		return errs.ConvertDB(err)
	}

	return nil
}

func (dao *scanPayDAO) Update(dc *gorm.DB, cond condition.IUpdate) error {
	cmd := dc.Table(dao.TableName()).
		Where(cond.Where()).
		Updates(cond.Update())
	if err := cmd.Error; err != nil {
		return errs.ConvertDB(err)
	}
	if cmd.RowsAffected == 0 {
		return errs.DBUpdateZeroRow
	}

	return nil
}
