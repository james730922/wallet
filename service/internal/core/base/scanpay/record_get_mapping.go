package scanpay

import (
	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
)

func newGetScanPayMapping() *commonGetScanPayMapping {
	return &commonGetScanPayMapping{}
}

type commonGetScanPayMapping struct {
}

func (hd *commonGetScanPayMapping) Handler(dc *gorm.DB, query *condition.ScanPayMappingQuery) (*model.ScanPayMapping, error) {
	cond := condition.NewQuery(query)

	result, err := dao.Mapping.First(dc, cond)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, err
	}
	return result, nil
}
