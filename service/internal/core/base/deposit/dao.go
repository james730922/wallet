package deposit

var (
	dao *storage
)

func newDAO() {
	dao = &storage{
		Common:            newCommonDAO(),
		Member:            newMemberDAO(),
		Account:           newAccountDAO(),
		Bonus:             newBonusDAO(),
		Config:            newDepositConfigDAO(),
		ConfigMember:      newDepositConfigMemberDAO(),
		ConfigMemberLevel: newDepositConfigMemberLevelDAO(),
	}
}

type storage struct {
	Common            ICommonDAO
	Member            IMemberDAO
	Account           IAccountDAO
	Bonus             IBonusDAO
	Config            IDepositConfigDAO
	ConfigMember      IDepositConfigMemberDAO
	ConfigMemberLevel IDepositConfigMemberLevelDAO
}
