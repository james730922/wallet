SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE `bank_code` (
  `code` varchar(10) NOT NULL COMMENT '銀行碼',
  `name` varchar(50) NOT NULL COMMENT '銀行名稱',
  `image` varchar(255) NOT NULL DEFAULT '' COMMENT '圖片路徑',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0=禁用 1=啟用',
  `added_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='銀行列表';

CREATE TABLE `bank_account` (
  `id` bigint(20) NOT NULL COMMENT '銀行帳戶ID',
  `type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '存款帳戶類型',
  `number` varchar(30) NOT NULL COMMENT '帳號',
  `bank_code` varchar(10) NOT NULL COMMENT '銀行碼',
  `bank_branch` varchar(100) NOT NULL DEFAULT '' COMMENT '銀行支行',
  `name` varchar(50) NOT NULL COMMENT '戶名',
  `currency_code` char(3) NOT NULL DEFAULT 'RMB' COMMENT '支付幣別',
  `levels` varchar(255) NOT NULL DEFAULT '0' COMMENT '可使用的會員層級，0=全部',
  `receive_limit` decimal(18,2) NOT NULL DEFAULT '0.00' COMMENT '收款上限，0=無上限',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '0=禁用 1=啟用',
  `visible` tinyint(4) NOT NULL DEFAULT '1' COMMENT '0=隱藏 1=顯示',
  `min_amount` decimal(18,2) NOT NULL DEFAULT '0.00' COMMENT '單筆最低金額',
  `max_amount` decimal(18,2) NOT NULL DEFAULT '0.00' COMMENT '單筆最高金額，0=無上限',
  `admin_id` bigint(20) DEFAULT NULL,
  `remark` varchar(255) NOT NULL DEFAULT '',
  `qrcode` varchar(500) NOT NULL DEFAULT '',
  `added_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_bank_account_number` (`number`),
  KEY `idx_bank_account_status_visible` (`status`,`visible`),
  KEY `idx_bank_account_bank_code` (`bank_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='前台存款帳戶';

CREATE TABLE `bank_deposit_category_type` (
  `id` int(11) NOT NULL,
  `name` varchar(50) NOT NULL,
  `method` varchar(45) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='存款分類類型';

CREATE TABLE `bank_deposit_category` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `type` int(11) NOT NULL COMMENT 'bank_deposit_category_type.id',
  `image` varchar(255) NOT NULL DEFAULT '',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0=禁用 1=啟用',
  `sort` int(11) NOT NULL DEFAULT '0',
  `added_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`),
  KEY `idx_deposit_category_status_sort` (`status`,`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='存款分類';

CREATE TABLE `bank_deposit_category_item` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `category_id` int(11) NOT NULL,
  `account_id` bigint(20) NOT NULL,
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '0=禁用 1=啟用',
  `sort` int(11) NOT NULL DEFAULT '0',
  `added_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_deposit_category_account` (`category_id`,`account_id`),
  KEY `idx_deposit_category_item_account` (`account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='存款分類帳戶關聯';

CREATE TABLE `member_level` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(40) NOT NULL,
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '0=禁用 1=啟用',
  `member_count` bigint(20) NOT NULL DEFAULT '0',
  `sort` int(11) NOT NULL DEFAULT '0',
  `note` varchar(50) NOT NULL DEFAULT '',
  `default` tinyint(4) NOT NULL DEFAULT '0' COMMENT '1=預設層級',
  `visible` tinyint(4) NOT NULL DEFAULT '1',
  `admin_id` bigint(20) DEFAULT NULL,
  `feature` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0=正常 1=黑名單',
  `added_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`),
  KEY `idx_member_level_default` (`default`),
  KEY `idx_member_level_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='會員層級';

CREATE TABLE `member` (
  `id` bigint(20) NOT NULL,
  `level_code` bigint(20) NOT NULL DEFAULT '1',
  `status` int(11) NOT NULL DEFAULT '1' COMMENT '0=凍結 1=可用 2=停權',
  `last_login_time` datetime(6) DEFAULT NULL,
  `failed_attempt_count` int(11) NOT NULL DEFAULT '0',
  `remark` text,
  `admin_id` bigint(20) DEFAULT NULL,
  `added_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`),
  KEY `idx_member_level_code` (`level_code`),
  KEY `idx_member_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='會員';

CREATE TABLE `member_mapping` (
  `id` bigint(20) NOT NULL,
  `country_code` varchar(7) NOT NULL,
  `mobile` varchar(20) NOT NULL,
  `name` varchar(20) NOT NULL,
  `qq` varchar(20) NOT NULL,
  `passwd` char(128) NOT NULL,
  `salt` char(36) NOT NULL,
  `simplify_id` varchar(8) NOT NULL,
  `security_passwd` char(128) NOT NULL DEFAULT '',
  `wallet_address` char(42) NOT NULL,
  `added_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_member_qq` (`qq`),
  UNIQUE KEY `uk_member_mobile` (`country_code`,`mobile`),
  UNIQUE KEY `uk_member_simplify_id` (`simplify_id`),
  UNIQUE KEY `uk_member_wallet_address` (`wallet_address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='會員登入與基本資料';

CREATE TABLE `member_simple_id_seq` (
  `id` tinyint(3) unsigned NOT NULL,
  `current_value` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='會員簡化ID計數';

CREATE TABLE `wallet_member` (
  `member_id` bigint(20) NOT NULL,
  `balance` decimal(18,2) NOT NULL DEFAULT '0.00',
  `total_amount` decimal(18,2) NOT NULL DEFAULT '0.00',
  `amount` decimal(18,2) NOT NULL DEFAULT '0.00',
  `bonus` decimal(18,2) NOT NULL DEFAULT '0.00',
  `frozen_amount` decimal(18,2) NOT NULL DEFAULT '0.00',
  `sign` char(128) NOT NULL,
  `added_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='會員錢包';

CREATE TABLE `wallet_journal` (
  `id` bigint(20) NOT NULL,
  `member_id` bigint(20) NOT NULL,
  `balance` decimal(18,2) NOT NULL,
  `total_amount` decimal(18,2) NOT NULL,
  `amount` decimal(18,2) NOT NULL,
  `bonus` decimal(18,2) NOT NULL,
  `frozen_amount` decimal(18,2) NOT NULL,
  `wallet_sign` char(128) NOT NULL,
  `added_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`),
  KEY `idx_wallet_journal_member_added` (`member_id`,`added_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='錢包資產異動日誌';

CREATE TABLE `deposit_config` (
  `id` int(11) NOT NULL,
  `bonus` decimal(7,4) NOT NULL DEFAULT '0.0000',
  `admin_id` bigint(20) DEFAULT NULL,
  `added_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='全域存款紅利設定';

CREATE TABLE `deposit_config_member_level` (
  `member_level` bigint(20) NOT NULL,
  `bonus` decimal(7,4) DEFAULT NULL,
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '0=關閉 1=啟用 2=刪除',
  `admin_id` bigint(20) DEFAULT NULL,
  `remarks` tinytext,
  `added_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`member_level`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='會員層級存款紅利設定';

CREATE TABLE `deposit_config_member` (
  `member_id` bigint(20) NOT NULL,
  `bonus` decimal(7,4) DEFAULT NULL,
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '0=關閉 1=啟用 2=刪除',
  `admin_id` bigint(20) DEFAULT NULL,
  `remarks` tinytext,
  `added_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='會員個人存款紅利設定';

CREATE TABLE `deposit` (
  `id` bigint(20) NOT NULL,
  `member_id` bigint(20) NOT NULL,
  `account_id` bigint(20) NOT NULL,
  `account_number` varchar(30) NOT NULL,
  `account_bank_code` varchar(10) NOT NULL,
  `currency_code` char(3) NOT NULL,
  `pay_name` varchar(20) NOT NULL,
  `amount` decimal(18,2) NOT NULL,
  `charge` decimal(18,2) NOT NULL DEFAULT '0.00',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0=待確認 1=成功 2=取消 3=暫停',
  `accept_time` datetime(6) DEFAULT NULL,
  `cancel_time` datetime(6) DEFAULT NULL,
  `admin_id` bigint(20) DEFAULT NULL,
  `remarks` tinytext,
  `sign` char(128) NOT NULL,
  `added_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`),
  KEY `idx_deposit_member_added` (`member_id`,`added_time`),
  KEY `idx_deposit_member_status` (`member_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='會員存款單';

CREATE TABLE `order_bonus` (
  `id` bigint(20) NOT NULL,
  `member_id` bigint(20) NOT NULL,
  `currency_code` char(3) NOT NULL,
  `amount` decimal(18,2) NOT NULL,
  `calculate_type` tinyint(4) NOT NULL DEFAULT '0',
  `source_type` tinyint(4) NOT NULL DEFAULT '0',
  `source_order_id` bigint(20) NOT NULL,
  `source_amount` decimal(18,2) NOT NULL,
  `source_rate` decimal(7,4) NOT NULL DEFAULT '0.0000',
  `status` tinyint(4) NOT NULL DEFAULT '0',
  `accept_time` datetime(6) DEFAULT NULL,
  `accept_admin_id` bigint(20) DEFAULT NULL,
  `sign` char(128) NOT NULL,
  `added_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`),
  KEY `idx_order_bonus_member` (`member_id`),
  KEY `idx_order_bonus_source` (`source_type`,`source_order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='存款紅利單';

CREATE TABLE `scanpay_record` (
  `id` bigint(20) NOT NULL,
  `amount` decimal(18,2) NOT NULL DEFAULT '0.00',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0=待確認 1=交易中 2=完成 3=取消 4=付款失敗作廢',
  `brand` varchar(20) NOT NULL DEFAULT '品牌',
  `merchant_id` varchar(50) NOT NULL DEFAULT '',
  `source_order_id` varchar(50) NOT NULL DEFAULT '',
  `content` varchar(500) NOT NULL,
  `expired_time` datetime(6) NOT NULL,
  `cancel_time` datetime(6) DEFAULT NULL,
  `admin_id` bigint(20) DEFAULT NULL,
  `remarks` varchar(255) DEFAULT NULL,
  `added_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`),
  KEY `idx_scanpay_record_status_expired` (`status`,`expired_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='掃碼支付來源紀錄';

CREATE TABLE `scanpay_mapping` (
  `record_id` bigint(20) NOT NULL,
  `merchant` varchar(20) NOT NULL DEFAULT '品牌',
  `merchant_order_id` varchar(50) NOT NULL,
  `merchant_member_id` varchar(50) NOT NULL,
  `merchant_member_account` varchar(50) NOT NULL DEFAULT '',
  `merchant_member_name` varchar(50) DEFAULT NULL,
  `amount` decimal(18,2) NOT NULL DEFAULT '0.00',
  `remark` varchar(255) DEFAULT NULL,
  `added_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`record_id`),
  UNIQUE KEY `uk_scanpay_merchant_order` (`merchant`,`merchant_order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='掃碼支付來源對應';

CREATE TABLE `order_scanpay` (
  `id` bigint(20) NOT NULL,
  `member_id` bigint(20) NOT NULL,
  `amount` decimal(18,2) NOT NULL,
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0=交易中 1=成功 2=失敗 3=取消',
  `record_id` bigint(20) NOT NULL,
  `brand` varchar(20) NOT NULL DEFAULT '品牌',
  `merchant_order_id` varchar(50) NOT NULL DEFAULT '',
  `source_order_id` varchar(50) NOT NULL DEFAULT '',
  `content` varchar(500) NOT NULL DEFAULT '',
  `success_time` datetime(6) DEFAULT NULL,
  `cancel_time` datetime(6) DEFAULT NULL,
  `sign` char(128) NOT NULL,
  `admin_id` bigint(20) DEFAULT NULL,
  `remarks` varchar(255) DEFAULT NULL,
  `added_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`),
  KEY `idx_order_scanpay_member_added` (`member_id`,`added_time`),
  KEY `idx_order_scanpay_status_updated` (`status`,`updated_time`),
  UNIQUE KEY `uk_order_scanpay_record` (`record_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='會員掃碼支付單';

CREATE TABLE `transaction` (
  `id` bigint(20) NOT NULL,
  `member_id` bigint(20) NOT NULL,
  `source_type` tinyint(4) NOT NULL COMMENT '1=存款 6=存款紅利 20=掃碼支付',
  `source_id` bigint(20) NOT NULL,
  `currency_code` char(3) NOT NULL,
  `amount` decimal(18,2) NOT NULL,
  `current_total_amount` decimal(18,2) NOT NULL,
  `changed_total_amount` decimal(18,2) NOT NULL,
  `current_balance` decimal(18,2) NOT NULL,
  `changed_balance` decimal(18,2) NOT NULL,
  `sign` char(128) NOT NULL,
  `remarks` tinytext,
  `merchant` varchar(20) DEFAULT NULL,
  `merchant_member_account` varchar(50) DEFAULT NULL,
  `added_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_transaction_source` (`source_type`,`source_id`),
  KEY `idx_transaction_member_source` (`member_id`,`source_type`),
  KEY `idx_transaction_member_added` (`member_id`,`added_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='前台交易流水';

SET FOREIGN_KEY_CHECKS = 1;
