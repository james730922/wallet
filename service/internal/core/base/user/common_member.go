package user

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"github.com/james730922/wallet/service/internal/core/base/wallet"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/errs"
	"github.com/james730922/wallet/service/internal/utils/password"
)

func newMemberCommon() IMemberCommon {
	return &memberCommonUseCase{}
}

type IMemberCommon interface {
	// 用手機號找會員
	Find(ctx context.Context, countryCode, mobile string) (*model.Member, error)
	// 找會員手機
	FindMobile(ctx context.Context, id int64) (*model.MemberMapping, error)
	// 取得會員資料
	Get(ctx context.Context, id int64) (*model.Member, error)
	// 建立會員以密碼登錄
	CreateWithPasswd(ctx context.Context, countryCode, mobile, passwd, qq, name string) (*model.Member, error)
	// 更新登入狀態或最後登入時間
	Update(ctx context.Context, cond *condition.MemberUpdateCond) error
	// 找mapping
	FindMapping(ctx context.Context, cond *condition.MemberMappingQuery) (*model.MemberMapping, error)
	FindMappingWithMobile(ctx context.Context, countryCode, mobile string) (*model.MemberMapping, error)
	FindMappingWithQQ(ctx context.Context, qq string) (*model.MemberMapping, error)
	UpdateMapping(ctx context.Context, mapping *model.MemberMapping) error
	SetInitialSecurityPasswd(ctx context.Context, memberID int64, passwd string) (bool, error)
	RecordLoginFailure(ctx context.Context, memberID int64) error
	ResetLoginFailures(ctx context.Context, memberID int64) error
	// 驗證掃碼支付密碼
	CheckScanPayPasswd(ctx context.Context, pwd string) error
}

type memberCommonUseCase struct {
	checkScanPayPasswd *checkMemberScanPayPasswd
}

func (cm *memberCommonUseCase) Find(ctx context.Context, countryCode, mobile string) (*model.Member, error) {
	cond := condition.MemberMappingQuery{
		CountryCode: aws.String(countryCode),
		Mobile:      aws.String(mobile),
	}
	mapping, err := dao.Member.GetMapping(packet.DB.New(), condition.NewQuery(cond))
	if err != nil {
		if err == errs.DBNoRow {
			return nil, err
		}
		logger.ApLog().Errorf("GetMappingByMobile params: [countryCode: %s, mobile: %s], err: %v", countryCode, mobile, err)
		return nil, errs.CommonNoData
	}

	member, err := dao.Member.First(packet.DB.New(), mapping.ID)
	if err != nil {
		if err == errs.DBNoRow {
			return nil, err
		}
		logger.ApLog().Errorf("Get params: [countryCode: %s, mobile: %s], err: %v", countryCode, mobile, err)
		return nil, errs.CommonNoData
	}

	return member, nil
}

func (cm *memberCommonUseCase) FindMobile(ctx context.Context, id int64) (*model.MemberMapping, error) {
	cond := condition.MemberMappingQuery{
		ID: aws.Int64(id),
	}

	return dao.Member.GetMapping(packet.DB.New(), condition.NewQuery(cond))
}

func (cm *memberCommonUseCase) Get(ctx context.Context, id int64) (*model.Member, error) {
	member, err := dao.Member.First(packet.DB.New(), id)
	if err != nil {
		if err == errs.DBNoRow {
			return nil, errs.MemberNotFound
		}

		logger.ApLog().Errorf("Get params: [id: %d], err: %v", id, err)
		return nil, err
	}
	return member, nil
}

func (cm *memberCommonUseCase) CreateWithPasswd(ctx context.Context, countryCode, mobile, passwd, qq, name string) (*model.Member, error) {
	memberID := packet.Node.Generate().Int64()

	now := time.Now().UTC()
	salt := uuid.New().String()
	r := rand.NewSource(time.Now().UTC().UnixNano())
	walletAddress := wallet.GenerateWalletAddress(int(r.Int63()))
	simpleID, err := packet.SimplifyMemberID.Generate()
	if err != nil {
		logger.ApLog().Error(err)
		return nil, err
	}
	passwdHash, err := password.Hash(passwd)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, errs.MemberCreateFailed
	}

	tx := func(tx *gorm.DB) error {
		if err := dao.Member.InsertMapping(tx, &model.MemberMapping{
			ID:            memberID,
			CountryCode:   countryCode,
			Mobile:        mobile,
			Name:          name,
			QQ:            qq,
			Passwd:        passwdHash,
			Salt:          salt,
			SimplifyID:    simpleID,
			AddedTime:     now,
			UpdatedTime:   now,
			WalletAddress: walletAddress,
		}); err != nil {
			logger.ApLog().Error(err)
			return errors.New(fmt.Sprintf("Insert mapping params: [countryCode: %s, mobile: %s], err: %v", countryCode, mobile, err))
		}

		// 取得默認會員級別
		memberLevel, err := self.MemberLevel.First(tx, &condition.MemberLevelQuery{Default: aws.Int(1)})
		if err != nil {
			logger.ApLog().Error(err)
			return errors.New(fmt.Sprintf("Find default member level err: %v", err))
		}

		// 更新會員分類總數
		cond := &condition.MemberLevelUpdate{
			ID:          aws.Int64(memberLevel.ID),
			MemberCount: aws.Int64(memberLevel.MemberCount + 1),
		}
		if err := self.MemberLevel.Update(tx, cond); err != nil {
			logger.ApLog().Error(err)
			return errors.New(fmt.Sprintf("Update default member level count err: %v", err))
		}

		if err := dao.Member.Insert(tx, &model.Member{
			ID:                 memberID,
			LevelCode:          memberLevel.ID,
			Status:             model.MemberStatusEnabled,
			FailedAttemptCount: 0,
			Remark:             aws.String(""),
			AddedTime:          now,
			UpdatedTime:        now,
		}); err != nil {
			logger.ApLog().Error(err)
			return errors.New(fmt.Sprintf("Insert member params: [countryCode: %s, mobile: %s], err: %v", countryCode, mobile, err))
		}

		if err := packet.Wallet.Create(tx, memberID); err != nil {
			logger.ApLog().Error(err)
			return errors.New(fmt.Sprintf("Insert wallet params: [countryCode: %s, mobile: %s], err: %v", countryCode, mobile, err))
		}

		return nil
	}

	if err := packet.DB.Transaction(tx); err != nil {
		logger.ApLog().Error(err)
		return nil, errs.MemberCreateFailed
	}

	member, err := dao.Member.First(packet.DB.New(), memberID)
	if err != nil {
		logger.ApLog().Errorf("Get params: [countryCode: %s, mobile: %s], err: %v", countryCode, mobile, err)
		return nil, errs.MemberCreateFailed
	}

	return member, nil
}

func (cm *memberCommonUseCase) Update(ctx context.Context, cond *condition.MemberUpdateCond) error {
	if cond.ID == nil {
		return errs.FrameworkIllegalParameter
	}

	updCond := condition.NewUpdate(
		&condition.MemberUpdate{
			Status:             cond.Status,
			LastLoginTime:      cond.LastLoginTime,
			FailedAttemptCount: cond.FailedAttemptCount,
			UpdatedTime:        aws.Time(time.Now().UTC()),
		},
		&condition.MemberUpdate{ID: cond.ID},
	)
	if err := dao.Member.Update(packet.DB.New(), updCond); err != nil {
		logger.ApLog().Errorf("Update member login state err: %v", err)
		return errs.CommonNoData
	}
	return nil
}
func (cm *memberCommonUseCase) FindMapping(ctx context.Context, cond *condition.MemberMappingQuery) (*model.MemberMapping, error) {
	mapping, err := dao.Member.GetMapping(packet.DB.New(), condition.NewQuery(cond))
	if err != nil {
		if err == errs.DBNoRow {
			return nil, err
		}
		logger.ApLog().Errorf("FindMapping params: [cond: %v], err: %v", cond, err)
		return nil, errs.CommonNoData
	}

	return mapping, nil
}

func (cm *memberCommonUseCase) FindMappingWithMobile(ctx context.Context, countryCode, mobile string) (*model.MemberMapping, error) {
	cond := condition.MemberMappingQuery{
		CountryCode: aws.String(countryCode),
		Mobile:      aws.String(mobile),
	}
	mapping, err := dao.Member.GetMapping(packet.DB.New(), condition.NewQuery(cond))
	if err != nil {
		if err == errs.DBNoRow {
			return nil, err
		}
		logger.ApLog().Errorf("FindMappingWithMobile params: [countryCode: %s, mobile: %s], err: %v", countryCode, mobile, err)
		return nil, errs.CommonNoData
	}

	return mapping, nil
}

func (cm *memberCommonUseCase) FindMappingWithQQ(ctx context.Context, qq string) (*model.MemberMapping, error) {
	mapping, err := dao.Member.GetMapping(packet.DB.New(), condition.NewQuery(condition.MemberMappingQuery{QQ: aws.String(qq)}))
	if err != nil {
		if err == errs.DBNoRow {
			return nil, err
		}
		logger.ApLog().Error(err)
		return nil, errs.CommonNoData
	}

	return mapping, nil
}

func (cm *memberCommonUseCase) UpdateMapping(ctx context.Context, mapping *model.MemberMapping) error {
	return dao.Member.UpdateMapping(packet.DB.New(), mapping)
}

func (cm *memberCommonUseCase) SetInitialSecurityPasswd(ctx context.Context, memberID int64, passwd string) (bool, error) {
	return dao.Member.SetInitialSecurityPasswd(packet.DB.New(), memberID, passwd)
}

func (cm *memberCommonUseCase) RecordLoginFailure(ctx context.Context, memberID int64) error {
	return dao.Member.IncrementFailedAttempt(packet.DB.New(), memberID)
}

func (cm *memberCommonUseCase) ResetLoginFailures(ctx context.Context, memberID int64) error {
	return dao.Member.ResetFailedAttempts(packet.DB.New(), memberID)
}

func (cm *memberCommonUseCase) CheckScanPayPasswd(ctx context.Context, pwd string) error {
	return cm.checkScanPayPasswd.Handler(ctx, pwd)
}
