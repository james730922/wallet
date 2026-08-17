package simplifymemberid

import (
	"github.com/jinzhu/gorm"

	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newSimplifyMemberIDDAO() ISimplifyMemberIDDAO {
	return &simplifyMemberIDDAO{}
}

type ISimplifyMemberIDDAO interface {
	GenSeq(dc *gorm.DB) (int64, error)
	FirstOriginalID(dc *gorm.DB, simpleID string) (int64, error)
	FirstSimplifyID(dc *gorm.DB, originID int64) (string, error)
}

type simplifyMemberIDDAO struct{}

func (dao *simplifyMemberIDDAO) GenSeq(dc *gorm.DB) (int64, error) {
	result := dc.Exec(`
		UPDATE member_simple_id_seq
		SET current_value = LAST_INSERT_ID(current_value + 1)
		WHERE id = 1`)
	if result.Error != nil {
		return 0, errs.ConvertDB(result.Error)
	}
	if result.RowsAffected != 1 {
		return 0, errs.DBUpdateZeroRow
	}

	var count int64
	if err := dc.Raw("SELECT LAST_INSERT_ID()").Row().Scan(&count); err != nil {
		return 0, errs.ConvertDB(err)
	}

	return count, nil
}

func (dao *simplifyMemberIDDAO) FirstOriginalID(dc *gorm.DB, simpleID string) (int64, error) {
	var result model.MemberMapping

	if err := dc.
		Select("id").
		Where("simplify_id = ?", simpleID).
		First(&result).
		Error; err != nil {
		return 0, errs.ConvertDB(err)
	}

	return result.ID, nil
}

func (dao *simplifyMemberIDDAO) FirstSimplifyID(dc *gorm.DB, originID int64) (string, error) {
	var result model.MemberMapping

	if err := dc.
		Select("simplify_id").
		Where("id = ?", originID).
		First(&result).
		Error; err != nil {
		return "", errs.ConvertDB(err)
	}

	return result.SimplifyID, nil
}
