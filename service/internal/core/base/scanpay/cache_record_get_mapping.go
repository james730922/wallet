package scanpay

import (
	"context"
	"sync"

	"github.com/jinzhu/gorm"

	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
)

type cacheScanPayRecordMapping struct {
	mx sync.Mutex
	m  sync.Map
}

func (c *cacheScanPayRecordMapping) Get(ctx context.Context, dc *gorm.DB, recordId int64) (*model.ScanPayMapping, error) {
	if v, ok := c.m.Load(recordId); ok {
		return v.(*model.ScanPayMapping), nil
	}

	c.mx.Lock()
	defer c.mx.Unlock()

	if v, ok := c.m.Load(recordId); ok {
		return v.(*model.ScanPayMapping), nil
	}

	scanPayRecordMapping, err := self.Record.GetScanPayMapping(dc.New(), &condition.ScanPayMappingQuery{
		RecordID: &recordId,
	})

	if err != nil {
		logger.ApLog().Warnf("err:%v,[cacheScanPayRecordMapping][get]recordId:%v", err, recordId)
		return nil, err
	}

	c.m.Store(recordId, scanPayRecordMapping)

	return scanPayRecordMapping, nil
}
