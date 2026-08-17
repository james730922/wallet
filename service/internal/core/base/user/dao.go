package user

var (
	dao *storage
)

func newDAO() {
	dao = &storage{
		Member:      newMemberDAO(),
		MemberLevel: newMemberLevelDAO(),
	}
}

type storage struct {
	Member      iMemberDAO
	MemberLevel iMemberLevelDAO
}
