package deposit

import (
	"context"
	"sort"
	"strconv"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/james730922/wallet/service/internal/models/condition"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/pb/zqbapis"
	"github.com/james730922/wallet/service/internal/thirdparty/fileserver"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/ctxs"
	"github.com/james730922/wallet/service/internal/utils/errs"
	"github.com/james730922/wallet/service/internal/utils/tools"
)

func newDepositMethod() *depositMemberMethods {
	return &depositMemberMethods{}
}

type depositMemberMethods struct {
	mx sync.Mutex
}

func (de *depositMemberMethods) Handler(ctx context.Context) (*zqbapis.DepositMethodResp, error) {
	// methods
	tmpMethods, err := de.getMethods(ctx)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, errs.CommonNoData
	}

	// items
	tmpItems, err := de.getItems(ctx)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, errs.CommonNoData
	}

	// 去掉沒有item 的methods
	// 照著methods 的排序塞入items
	methods := make([]*zqbapis.DepositMethod, 0, len(tmpMethods))
	items := make([]*zqbapis.DepositMethodItem, 0, len(tmpItems))
	for _, tmpMethod := range tmpMethods {
		id, err := strconv.ParseInt(tmpMethod.Id, 10, 64)
		if err != nil {
			logger.ApLog().Error(err)
			return nil, errs.CommonRequestParamParseFailed
		}
		if tmpItem, ok := tmpItems[id]; ok {
			methods = append(methods, tmpMethod)
			items = append(items, tmpItem)
		}
	}

	// 組成resp
	resp := &zqbapis.DepositMethodResp{
		Methods: methods,
		Items:   items,
	}

	if len(methods) == 0 || len(items) == 0 {
		return resp, errs.DepositNoMethod
	}

	return resp, err
}

func (de *depositMemberMethods) getMethods(ctx context.Context) ([]*zqbapis.DepositMethod, error) {
	// methods
	categories, err := packet.CategoryCommon.List(ctx, &condition.DepositCategoryQuery{})
	if err != nil {
		return nil, err
	}

	categoryType, err := packet.CategoryCommon.ListTypeMap(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*zqbapis.DepositMethod, 0, len(categories))
	for _, category := range categories {
		tmp := &zqbapis.DepositMethod{
			Id:   strconv.FormatInt(category.ID, 10),
			Name: category.Name,
			Type: de.findType(categoryType, category.Type),
			Icon: packet.FileServer.URL(category.Image, fileserver.WithURLUpdateTime(category.UpdatedTime)),
		}
		result = append(result, tmp)
	}

	return result, nil
}

func (de *depositMemberMethods) getItems(ctx context.Context) (map[int64]*zqbapis.DepositMethodItem, error) {
	accounts, err := de.getAccounts(ctx)
	if err != nil {
		return nil, err
	}

	// 將抓出來的account 分類
	result := make(map[int64]*zqbapis.DepositMethodItem)
	for _, account := range accounts {
		val, ok := result[account.CategoryID]
		if !ok {
			val = &zqbapis.DepositMethodItem{
				Id:       strconv.FormatInt(account.CategoryID, 10),
				Accounts: make([]*zqbapis.DepositAccount, 0, 5),
			}

			result[account.CategoryID] = val
		}

		// 轉換成DepositMethodItem
		depositAccount := &zqbapis.DepositAccount{
			Id:            strconv.FormatInt(account.AccountID, 10),
			BankName:      account.BankName,
			BankImage:     packet.FileServer.URL(account.BankImage),
			AccountName:   account.Name,
			AccountNumber: account.Number,
			AccountBranch: account.BankBranch,
			Min:           tools.DecimalToFloat64(account.MinAmount),
			Max:           tools.DecimalToFloat64(account.MaxAmount),
			QrCode:        packet.FileServer.URL(account.QRCode, fileserver.WithURLUpdateTime(account.AccountUpdatedTime)),
		}

		val.Accounts = append(val.Accounts, depositAccount)
	}

	return result, nil
}

// 取得所有銀行帳號
func (de *depositMemberMethods) getAccounts(ctx context.Context) ([]*model.BankDepositCategoryItemView, error) {
	// 取得會員級別
	id, ok := ctxs.GetMemberID(ctx)
	if !ok {
		return nil, errs.CommonNoMemberID
	}
	member, err := packet.Member.Get(ctx, id)
	if err != nil {
		logger.ApLog().Error(err)
		return nil, errs.CommonNoData
	}

	// 取得銀行帳號
	categoryCond := &condition.BankDepositCategoryItemViewCond{
		AccountStatus:  aws.Int(int(model.BankAccountStatusEnable)),
		CategoryStatus: aws.Int(int(model.DepositCategoryStatusEnable)),
		MemberLevel:    aws.Int(int(member.LevelCode)),
	}
	items, err := self.DepositCommon.CategoryItem(ctx, categoryCond)
	if err != nil {
		return nil, errs.CommonNoData
	}

	// 取得銀行名稱
	cond := &condition.BankCodeQuery{
		Status: aws.Int(1),
	}
	banks, err := packet.Bank.Map(ctx, cond)
	if err != nil {
		return nil, errs.CommonNoData
	}

	// packet bankName
	for _, item := range items {
		bankName := item.BankCode
		if bank, ok := banks[item.BankCode]; ok {
			bankName = bank.Name
		}
		item.BankName = bankName
	}

	sort.SliceStable(items, func(i, j int) bool { return items[i].Sort < items[j].Sort })

	return items, nil
}

func (de *depositMemberMethods) findType(m map[int]*model.BankDepositCategoryType, key int) string {
	if data, ok := m[key]; ok {
		return data.Method
	}
	return ""
}
