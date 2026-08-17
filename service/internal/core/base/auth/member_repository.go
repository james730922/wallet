package auth

import (
	"sync"

	"github.com/james730922/wallet/service/internal/models/model"
)

type iMemberRepository interface {
	Store(token string, member *model.Member)
	GetByID(id int64) (*model.MemberRepository, bool)
	GetAllID() []int64
	DeleteByID(id int64)
}

func newMemberRepository() *memberRepository {
	return &memberRepository{
		id:    make(map[int64]*model.MemberRepository),
		token: make(map[string]int64),
	}
}

type memberRepository struct {
	mx    sync.RWMutex
	id    map[int64]*model.MemberRepository
	token map[string]int64
}

func (mr *memberRepository) Store(token string, member *model.Member) {
	mr.mx.Lock()
	mr.deleteByID(member.ID)
	mr.id[member.ID] = &model.MemberRepository{
		Token: token,
		ID:    member.ID,
	}
	mr.token[token] = member.ID
	mr.mx.Unlock()
}

func (mr *memberRepository) GetByID(id int64) (*model.MemberRepository, bool) {
	var result *model.MemberRepository
	var ok bool

	mr.mx.RLock()
	if r, o := mr.id[id]; o {
		result, ok = &model.MemberRepository{
			Token: r.Token,
			ID:    id,
		}, true
	}
	mr.mx.RUnlock()

	return result, ok
}

func (mr *memberRepository) GetAllID() []int64 {
	var result []int64

	mr.mx.RLock()
	result = make([]int64, 0, len(mr.id))
	for key := range mr.id {
		result = append(result, key)
	}
	mr.mx.RUnlock()

	return result
}

func (mr *memberRepository) DeleteByID(id int64) {
	mr.mx.Lock()
	mr.deleteByID(id)
	mr.mx.Unlock()
}

func (mr *memberRepository) deleteByID(id int64) {
	if m, ok := mr.id[id]; ok {
		delete(mr.id, id)
		delete(mr.token, m.Token)
	}
}
