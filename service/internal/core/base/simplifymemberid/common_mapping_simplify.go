package simplifymemberid

import (
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/thirdparty/cache"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

type commonMappingSimplify struct{}

func newCommonMappingSimplify() *commonMappingSimplify {
	return &commonMappingSimplify{}
}

func (hd *commonMappingSimplify) Handler(id int64) (string, error) {
	cache, ok := hd.getCache(id)
	if !ok {
		tmp, err := hd.getFromDB(packet.DBSlave.New(), id)
		if err != nil {
			return "", err
		}
		cache = tmp
	}

	return cache, nil
}

func (hd *commonMappingSimplify) getCache(id int64) (string, bool) {
	tmp, ok := cache.Cache().Get(hd.getKey(id))
	if !ok {
		return "", false
	}

	cache, ok := tmp.(string)
	if !ok {
		return "", false
	}

	return cache, true
}

func (hd *commonMappingSimplify) getFromDB(dc *gorm.DB, id int64) (string, error) {
	simpleID, err := dao.ISimplifyMemberID.FirstSimplifyID(dc, id)
	if err != nil {
		if err == errs.DBNoRow {
			return "", errs.MemberNotFound
		}
		return "", err
	}

	cache.Cache().Set(hd.getKey(id), simpleID, -1)
	return simpleID, nil
}

func (hd *commonMappingSimplify) getKey(id int64) string {
	return cache.KeyMemberOriginalID + "-" + strconv.FormatInt(id, 10)
}
