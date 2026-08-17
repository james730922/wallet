package simplifymemberid

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

type commonGen struct{}

func newCommonGen() *commonGen {
	return &commonGen{}
}

func (hd *commonGen) Handler() (string, error) {
	var count int64
	var simplifyID string
	err := packet.DB.Transaction(func(tx *gorm.DB) error {
		var err error
		count, err = dao.ISimplifyMemberID.GenSeq(tx)
		if err != nil {
			return err
		}
		simplifyID, err = formatSimplifyID(count)
		return err
	})
	if err != nil {
		logger.ApLog().Errorf("dao.ISimplifyMemberID.GenSeq(), err:%v", err)
		return "", err
	}

	return simplifyID, nil
}

func formatSimplifyID(count int64) (string, error) {
	if count <= 0 {
		return "", errs.CommonNoMemberID
	}
	if count > maxSimplifyID {
		return "", errs.MemberSimplifyIDExhausted
	}

	return fmt.Sprintf("%0*d", width, count), nil
}
