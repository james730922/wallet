package apictrl

import (
	"go.uber.org/dig"

	"github.com/james730922/wallet/service/internal/core/base/auth"
	"github.com/james730922/wallet/service/internal/core/base/bank"
	"github.com/james730922/wallet/service/internal/core/base/deposit"
	"github.com/james730922/wallet/service/internal/core/usecase/member/scanpaymember"
	"github.com/james730922/wallet/service/internal/core/usecase/member/transactionmember"
)

type apiControllerSet struct {
	dig.In
	AuthLogin   auth.ILoginMember
	AuthToken   auth.IToken
	Deposit     deposit.IDepositMember
	Transaction transactionmember.IMember
	Bank        bank.IBankCommon
	ScanPay     scanpaymember.IMember
}

func NewController(set apiControllerSet) *Controller {
	return &Controller{
		AuthLogin:     newAuthLoginController(set.AuthLogin),
		HttpAuthToken: newAuthTokenHttpController(set.AuthToken),
		Deposit:       newDepositController(set.Deposit),
		Transaction:   newTransactionController(set.Transaction),
		Bank:          newBankController(set.Bank),
		ScanPay:       newScanPayController(set.ScanPay),
	}
}

type Controller struct {
	AuthLogin     IAuthLogin
	HttpAuthToken IAuthTokenHttpController
	Deposit       IDepositController
	Transaction   ITransactionController
	Bank          IBankController
	ScanPay       IScanPayController
}
