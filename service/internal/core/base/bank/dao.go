package bank

var (
	dao *storage
)

func newDAO() {
	dao = &storage{
		Bank:     newBankDAO(),
		Category: newCategoryDAO(),
	}
}

type storage struct {
	Bank     IBankDAO
	Category ICategoryDAO
}
