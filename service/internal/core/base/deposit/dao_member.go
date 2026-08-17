package deposit

import (
	"github.com/jinzhu/gorm"
	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newMemberDAO() IMemberDAO {
	return &memberDAO{}

}

type IMemberDAO interface {
	SelectCategoryItem(dc *gorm.DB) ([]*model.BankDepositCategoryItemView, error)
}

type memberDAO struct {
	model.BankAccount
}

func (dao *memberDAO) SelectCategoryItem(dc *gorm.DB) ([]*model.BankDepositCategoryItemView, error) {
	result := make([]*model.BankDepositCategoryItemView, 0, 5)

	err := dc.Raw("SELECT `dc`.`id` AS `id`," +
		"`dc`.`name` AS `category_name`," +
		"`dc`.`status` AS `category_status`," +
		"`dci`.`sort` AS `sort`," +
		"`acc`.`id` AS `account_id`," +
		"`acc`.`number`," +
		"`acc`.`currency_code`," +
		"`acc`.`levels`," +
		"`acc`.`name`," +
		"`acc`.`bank_code`," +
		"`acc`.`bank_branch`," +
		"`acc`.`min_amount`," +
		"`acc`.`max_amount`," +
		"`acc`.`qrcode`," +
		"`acc`.`status` AS `account_status`," +
		"`acc`.`visible` AS `account_visible`," +
		"`acc`.`updated_time` AS `account_updated_time`," +
		"`code`.`image` AS `bank_image`" +
		"FROM `bank_deposit_category` AS `dc`" +
		"JOIN `bank_deposit_category_item` AS `dci`" +
		"ON `dc`.`id` = `dci`.`category_id`" +
		"JOIN `bank_account` AS `acc`" +
		"ON `dci`.`account_id` = `acc`.`id`" +
		"JOIN `bank_code` AS `code`" +
		"ON `acc`.`bank_code` = `code`.`code`").Find(&result).Error

	if err != nil {
		return nil, errs.ConvertDB(err)
	}

	return result, nil
}
