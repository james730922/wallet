package wallet

var (
	dao *storage
)

func newDAO() {
	dao = &storage{
		Wallet:  newWalletDAO(),
		Journal: newJournalDAO(),
	}
}

type storage struct {
	Wallet  IWalletDAO
	Journal IJournalDAO
}
