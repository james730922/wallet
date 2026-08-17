package scanpay

import (
	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
)

func newRecordCommon() IRecordCommon {
	return &recordCommonUseCase{
		getScanPayMapping: newGetScanPayMapping(),
	}

}

type IRecordCommon interface {
	GetScanPayMapping(dc *gorm.DB, query *condition.ScanPayMappingQuery) (*model.ScanPayMapping, error)
	GetScanPayRecord(dc *gorm.DB, id int64) (*model.ScanPayRecord, error)
	GetScanPayRecordForUpdate(dc *gorm.DB, id int64) (*model.ScanPayRecord, error)
}

type recordCommonUseCase struct {
	getScanPayMapping *commonGetScanPayMapping
}

func (cm *recordCommonUseCase) GetScanPayMapping(dc *gorm.DB, query *condition.ScanPayMappingQuery) (*model.ScanPayMapping, error) {
	return cm.getScanPayMapping.Handler(dc, query)
}

func (cm *recordCommonUseCase) GetScanPayRecord(dc *gorm.DB, id int64) (*model.ScanPayRecord, error) {
	return dao.Record.First(dc, id)
}

func (cm *recordCommonUseCase) GetScanPayRecordForUpdate(dc *gorm.DB, id int64) (*model.ScanPayRecord, error) {
	return dao.Record.FirstForUpdate(dc, id)
}
