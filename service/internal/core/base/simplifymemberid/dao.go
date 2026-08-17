package simplifymemberid

var (
	dao *storage
)

func newDAO() {
	dao = &storage{
		ISimplifyMemberID: newSimplifyMemberIDDAO(),
	}
}

type storage struct {
	ISimplifyMemberID ISimplifyMemberIDDAO
}
