package transaction

var (
	dao *storage
)

func newDAO() {
	dao = &storage{
		Transaction: newTransactionDAO(),
	}
}

type storage struct {
	Transaction ITransactionDAO
}
