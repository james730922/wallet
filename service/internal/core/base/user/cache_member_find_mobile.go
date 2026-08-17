package user

import (
	"context"
	"sync"

	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
)

type cacheMemberFindMobile struct {
	mx sync.Mutex
	m  sync.Map
}

func (c *cacheMemberFindMobile) Get(ctx context.Context, memberID int64) (*model.MemberMapping, error) {
	if v, ok := c.m.Load(memberID); ok {
		return v.(*model.MemberMapping), nil
	}

	c.mx.Lock()
	defer c.mx.Unlock()

	if v, ok := c.m.Load(memberID); ok {
		return v.(*model.MemberMapping), nil
	}

	member, err := self.Member.FindMobile(ctx, memberID)
	if err != nil {
		logger.ApLog().Warnf("err:%v,[cacheMemberFindMobile][get]memberID:%v", err, memberID)
		return nil, err
	}

	c.m.Store(memberID, member)

	return member, nil
}
