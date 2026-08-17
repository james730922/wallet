package cache

import (
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/james730922/wallet/service/internal/utils/conf"
)

const (
	KeyBankCode                         = "bank_code"
	KeyBankDepositCategory              = "bank_deposit_category"
	KeyBankDepositCategoryItem          = "bank_deposit_category_item"
	KeyBankDepositCategoryType          = "bank_deposit_category_type"
	KeyBankDepositCategoryItemForMember = "bank_deposit_category_item_for_member"
	KeyMemberLevel                      = "member_level"
	KeyDepositConfig                    = "deposit_config"
	KeyDepositConfigMember              = "deposit_config_member"
	KeyDepositConfigMemberLevel         = "deposit_config_member_level"
	KeyBrand                            = "brand"
	KeyBrandGroup                       = "brand_group"
	KeyBrandMap                         = "brand_map"
	KeyBrandGroupMap                    = "brand_group_map"
	KeyBrandGroupCodeMap                = "brand_group_code_map"
	KeyMemberSimpleID                   = "member_simple_id"
	KeyMemberOriginalID                 = "member_original_id"

	KeyMenuTree = "menu_tree"
)

var c *cache.Cache

func Cache() *cache.Cache {
	return c
}

func NewCache() {
	c = newCache(conf.Cache().GetDefaultExpiration())
}

func newCache(defaultExpiration time.Duration) *cache.Cache {
	return cache.New(defaultExpiration, 10*time.Second)
}
