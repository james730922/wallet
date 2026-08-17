package simplifymemberid

import (
	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/thirdparty/cache"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

type commonMappingOriginal struct{}

func newCommonMappingOriginal() *commonMappingOriginal {
	return &commonMappingOriginal{}
}

func (hd *commonMappingOriginal) Handler(simplifyID string) (int64, error) {
	cache, ok := hd.getCache(simplifyID)
	if !ok {
		tmp, err := hd.getFromDB(packet.DBSlave.New(), simplifyID)
		if err != nil {
			return 0, err
		}
		cache = tmp
	}

	return cache, nil
}

func (hd *commonMappingOriginal) getCache(simplifyID string) (int64, bool) {
	tmp, ok := cache.Cache().Get(hd.getKey(simplifyID))
	if !ok {
		return 0, false
	}

	cache, ok := tmp.(int64)
	if !ok {
		return 0, false
	}

	return cache, true
}

func (hd *commonMappingOriginal) getFromDB(dc *gorm.DB, simplifyID string) (int64, error) {
	id, err := dao.ISimplifyMemberID.FirstOriginalID(dc, simplifyID)
	if err != nil {
		if err == errs.DBNoRow {
			return 0, errs.MemberNotFound
		}
		return 0, err
	}

	cache.Cache().Set(hd.getKey(simplifyID), id, -1)
	return id, nil
}

func (hd *commonMappingOriginal) getKey(simplifyID string) string {
	return cache.KeyMemberSimpleID + "-" + simplifyID
}
