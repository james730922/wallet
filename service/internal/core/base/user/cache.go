package user

func NewUserCache() *userCache {
	return &userCache{}
}

type userCache struct {
}

func (userCache) CacheMember() *cacheMember {
	return &cacheMember{}
}

func (userCache) CacheMemberMapping() *cacheMemberMapping {
	return &cacheMemberMapping{}
}

func (userCache) CacheMemberFindMobile() *cacheMemberFindMobile {
	return &cacheMemberFindMobile{}
}
