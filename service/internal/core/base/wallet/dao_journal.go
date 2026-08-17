package wallet

import (
	"github.com/jinzhu/gorm"

	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newJournalDAO() IJournalDAO {
	return &journalDAO{}
}

type IJournalDAO interface {
	Insert(dc *gorm.DB, journal *model.JournalMemberWallet) error
}

type journalDAO struct{}

func (dao *journalDAO) Insert(dc *gorm.DB, journal *model.JournalMemberWallet) error {
	if err := dc.Create(journal).Error; err != nil {
		return errs.ConvertDB(err)
	}

	return nil
}
