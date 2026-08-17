package scanpay

func NewSanPayCache() *scanPayCache {
	return &scanPayCache{}
}

type scanPayCache struct {
}

func (scanPayCache) CacheScanPayRecordMapping() *cacheScanPayRecordMapping {
	return &cacheScanPayRecordMapping{}
}
