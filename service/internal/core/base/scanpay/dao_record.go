package scanpay

import (
	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/db"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newScanPayRecordDAO() IScanPayRecordDAO {
	return &scanPayRecordDAO{}
}

type IScanPayRecordDAO interface {
	First(dc *gorm.DB, id int64) (*model.ScanPayRecord, error)
	FirstForUpdate(dc *gorm.DB, id int64) (*model.ScanPayRecord, error)
	Update(dc *gorm.DB, cond condition.IUpdate) error
}

type scanPayRecordDAO struct {
	model.ScanPayRecord
}

func (dao *scanPayRecordDAO) First(dc *gorm.DB, id int64) (*model.ScanPayRecord, error) {
	var result model.ScanPayRecord
	err := dc.
		Where("id = ?", id).
		First(&result).Error
	if err != nil {
		return &result, err
	}
	return &result, err
}

func (dao *scanPayRecordDAO) FirstForUpdate(dc *gorm.DB, id int64) (*model.ScanPayRecord, error) {
	var result model.ScanPayRecord

	if err := dc.Set(db.GormSetSelectForUpdate()).
		Where("id = ?", id).First(&result).Error; err != nil {
		return nil, errs.ConvertDB(err)
	}

	return &result, nil
}

func (dao *scanPayRecordDAO) Update(dc *gorm.DB, cond condition.IUpdate) error {
	cmd := dc.Table(dao.TableName()).
		Scopes(db.ParseWhere(cond.Where())).
		Updates(cond.Update())

	if err := cmd.Error; err != nil {
		return errs.ConvertDB(err)
	}
	if cmd.RowsAffected == 0 {
		return errs.DBUpdateZeroRow
	}

	return nil
}
