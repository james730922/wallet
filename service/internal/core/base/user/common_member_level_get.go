package user

import (
	"context"
	"sync"

	"github.com/ahmetb/go-linq/v3"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/cache"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
)

func newMemberLevelCommonGet() *memberLevelCommonGet {
	return &memberLevelCommonGet{}
}

type memberLevelCommonGet struct {
	mx sync.Mutex
}

// 取得分級
func (m *memberLevelCommonGet) Handler(ctx context.Context, cond *condition.MemberLevelQuery) ([]*model.MemberLevel, error) {
	items, ok := m.getCache(ctx)
	if !ok {
		tmp, err := m.get(ctx)
		if err != nil {
			return nil, err
		}
		items = tmp
	}

	return m.find(items, cond), nil
}

func (m *memberLevelCommonGet) getCache(ctx context.Context) ([]*model.MemberLevel, bool) {
	tmp, ok := cache.Cache().Get(cache.KeyMemberLevel)
	if !ok {
		return nil, false
	}

	return tmp.([]*model.MemberLevel), true
}

func (m *memberLevelCommonGet) get(ctx context.Context) ([]*model.MemberLevel, error) {
	m.mx.Lock()
	defer m.mx.Unlock()

	tmp, ok := m.getCache(ctx)
	if ok {
		return tmp, nil
	}

	items, err := dao.MemberLevel.List(packet.DB.New())
	if err != nil {
		return nil, err
	}

	if err := cache.Cache().Add(cache.KeyMemberLevel, items, 0); err != nil {
		logger.ApLog().Error(err)
	}

	return items, nil
}

func (m *memberLevelCommonGet) find(source []*model.MemberLevel, cond *condition.MemberLevelQuery) []*model.MemberLevel {
	var result []*model.MemberLevel

	linq.From(source).
		Where(func(i interface{}) bool {
			if cond.ID == nil {
				return true
			}
			return i.(*model.MemberLevel).ID == aws.Int64Value(cond.ID)
		}).
		Where(func(i interface{}) bool {
			if cond.Name == nil {
				return true
			}
			return i.(*model.MemberLevel).Name == aws.StringValue(cond.Name)
		}).
		Where(func(i interface{}) bool {
			if cond.Visible == nil {
				return true
			}
			return int(i.(*model.MemberLevel).Visible) == aws.IntValue(cond.Visible)
		}).
		Where(func(i interface{}) bool {
			if cond.Status == nil {
				return true
			}
			return int(i.(*model.MemberLevel).Status) == aws.IntValue(cond.Status)
		}).
		ToSlice(&result)

	return result
}

func (m *memberLevelCommonGet) copy(from []*model.MemberLevel) []*model.MemberLevel {
	to := make([]*model.MemberLevel, 0, len(from))

	for _, fromSource := range from {
		to = append(to, &model.MemberLevel{
			ID:          fromSource.ID,
			Name:        fromSource.Name,
			Status:      fromSource.Status,
			MemberCount: fromSource.MemberCount,
			Sort:        fromSource.Sort,
			Note:        fromSource.Note,
			Default:     fromSource.Default,
			AdminID:     fromSource.AdminID,
			Visible:     fromSource.Visible,
			AddedTime:   fromSource.AddedTime,
			UpdatedTime: fromSource.UpdatedTime,
		})
	}

	return to
}
