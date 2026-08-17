package user

import (
	"context"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"sync"
)

type cacheMemberMapping struct {
	mx sync.Mutex
	m  sync.Map
}

func (c *cacheMemberMapping) Get(ctx context.Context, memberID int64) (*model.MemberMapping, error) {
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
		logger.ApLog().Errorf("err:%v,[cacheMemberMapping][get]memberID:%v", err, memberID)
		return nil, err
	}

	c.m.Store(memberID, member)

	return member, nil
}
