package user

import (
	"context"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"sync"
)

type cacheMember struct {
	mx sync.Mutex
	m  sync.Map
}

func (c *cacheMember) Get(ctx context.Context, memberID int64) (*model.Member, error) {
	if v, ok := c.m.Load(memberID); ok {
		return v.(*model.Member), nil
	}

	c.mx.Lock()
	defer c.mx.Unlock()

	if v, ok := c.m.Load(memberID); ok {
		return v.(*model.Member), nil
	}

	member, err := self.Member.Get(ctx, memberID)
	if err != nil {
		logger.ApLog().Errorf("err:%v,[cacheMember][get]memberID:%v", err, memberID)
		return nil, err
	}

	c.m.Store(memberID, member)

	return member, nil
}
