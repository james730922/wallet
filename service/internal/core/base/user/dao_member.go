package user

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/db"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newMemberDAO() iMemberDAO {
	return &memberDAO{}
}

type iMemberDAO interface {
	// 建立會員
	Insert(dc *gorm.DB, cond *model.Member) error
	// 取得會員資料
	First(dc *gorm.DB, id int64) (*model.Member, error)
	// 更新會員資料
	Update(dc *gorm.DB, cond condition.IUpdate) error
	// 建立mapping
	InsertMapping(dc *gorm.DB, cond *model.MemberMapping) error
	// 找會員
	GetMapping(dc *gorm.DB, cond condition.IQuery) (*model.MemberMapping, error)
	// Update mapping
	UpdateMapping(dc *gorm.DB, mapping *model.MemberMapping) error
	SetInitialSecurityPasswd(dc *gorm.DB, memberID int64, passwd string) (bool, error)
	IncrementFailedAttempt(dc *gorm.DB, memberID int64) error
	ResetFailedAttempts(dc *gorm.DB, memberID int64) error
}

type memberDAO struct {
}

func (dao *memberDAO) Insert(dc *gorm.DB, cond *model.Member) error {
	if err := dc.Create(cond).Error; err != nil {
		return errs.ConvertDB(err)
	}
	return nil
}

func (dao *memberDAO) First(dc *gorm.DB, id int64) (*model.Member, error) {
	var result model.Member

	if err := dc.Where("id = ?", id).Find(&result).Error; err != nil {
		return nil, errs.ConvertDB(err)
	}

	return &result, nil
}

func (dao *memberDAO) Update(dc *gorm.DB, cond condition.IUpdate) error {
	if err := dc.Model(model.Member{}).
		Scopes(db.ParseWhere(cond.Where())).
		Updates(cond.Update()).
		Error; err != nil {
		return errs.ConvertDB(err)
	}
	return nil
}

func (dao *memberDAO) InsertMapping(dc *gorm.DB, cond *model.MemberMapping) error {
	if err := dc.Create(cond).Error; err != nil {
		return errs.ConvertDB(err)
	}

	return nil
}

func (dao *memberDAO) GetMapping(dc *gorm.DB, cond condition.IQuery) (*model.MemberMapping, error) {
	var result model.MemberMapping

	if err := dc.Where(cond.Where()).First(&result).Error; err != nil {
		return nil, errs.ConvertDB(err)
	}

	return &result, nil
}

func (dao *memberDAO) UpdateMapping(dc *gorm.DB, mapping *model.MemberMapping) error {
	if err := dc.Save(mapping).Error; err != nil {
		return errs.ConvertDB(err)
	}

	return nil
}

func (dao *memberDAO) SetInitialSecurityPasswd(dc *gorm.DB, memberID int64, passwd string) (bool, error) {
	result := dc.Model(&model.MemberMapping{}).
		Where("id = ? AND COALESCE(security_passwd, '') = ''", memberID).
		Updates(map[string]interface{}{
			"security_passwd": passwd,
			"updated_time":    time.Now().UTC(),
		})
	if result.Error != nil {
		return false, errs.ConvertDB(result.Error)
	}
	return result.RowsAffected == 1, nil
}

func (dao *memberDAO) IncrementFailedAttempt(dc *gorm.DB, memberID int64) error {
	result := dc.Model(&model.Member{}).
		Where("id = ?", memberID).
		UpdateColumn("failed_attempt_count", gorm.Expr("failed_attempt_count + 1"))
	if result.Error != nil {
		return errs.ConvertDB(result.Error)
	}
	return nil
}

func (dao *memberDAO) ResetFailedAttempts(dc *gorm.DB, memberID int64) error {
	result := dc.Model(&model.Member{}).
		Where("id = ?", memberID).
		UpdateColumn("failed_attempt_count", 0)
	if result.Error != nil {
		return errs.ConvertDB(result.Error)
	}
	return nil
}
