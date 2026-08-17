package scanpay

import (
	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/db"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newScanPayMappingDAO() IScanPayMappingDAO {
	return &scanPayMappingDAO{}
}

type IScanPayMappingDAO interface {
	First(dc *gorm.DB, query condition.IQuery) (*model.ScanPayMapping, error)
}

type scanPayMappingDAO struct {
	model.ScanPayMapping
}

func (dao *scanPayMappingDAO) First(dc *gorm.DB, query condition.IQuery) (*model.ScanPayMapping, error) {
	result := &model.ScanPayMapping{}

	if err := dc.
		Scopes(db.ParseWhere(query.Where())).
		First(result).Error; err != nil {
		return nil, errs.ConvertDB(err)
	}

	return result, nil
}
