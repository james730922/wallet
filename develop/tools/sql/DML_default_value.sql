SET NAMES utf8mb4;

INSERT INTO `bank_code`
  (`code`, `name`, `image`, `status`, `added_time`, `updated_time`)
VALUES
  ('DEMO', 'Demo 銀行', '/bank/code/demo.png', 1, NOW(6), NOW(6));

INSERT INTO `bank_deposit_category_type`
  (`id`, `name`, `method`)
VALUES
  (1, '銀行轉帳', 'bank');

INSERT INTO `bank_deposit_category`
  (`id`, `name`, `type`, `image`, `status`, `sort`, `added_time`, `updated_time`)
VALUES
  (1, 'Demo 銀行卡', 1, '/bank/deposit/category/demo.png', 1, 1, NOW(6), NOW(6));

INSERT INTO `bank_account`
  (`id`, `type`, `number`, `bank_code`, `bank_branch`, `name`, `currency_code`,
   `levels`, `receive_limit`, `status`, `visible`, `min_amount`, `max_amount`,
   `admin_id`, `remark`, `qrcode`, `added_time`, `updated_time`)
VALUES
  (1, 0, '0000000000000000', 'DEMO', 'Demo 分行', '品牌 Demo 帳戶', 'RMB',
   '0', 0.00, 1, 1, 100.00, 50000.00,
   NULL, '僅供本機展示', '', NOW(6), NOW(6));

INSERT INTO `bank_deposit_category_item`
  (`id`, `category_id`, `account_id`, `status`, `sort`, `added_time`, `updated_time`)
VALUES
  (1, 1, 1, 1, 1, NOW(6), NOW(6));

INSERT INTO `deposit_config`
  (`id`, `bonus`, `admin_id`, `added_time`, `updated_time`)
VALUES
  (1, 0.0100, NULL, NOW(6), NOW(6));

INSERT INTO `member_level`
  (`id`, `name`, `status`, `member_count`, `sort`, `note`, `default`, `visible`,
   `admin_id`, `feature`, `added_time`, `updated_time`)
VALUES
  (1, 'Demo 會員', 1, 1, 1, '本機預設會員層級', 1, 1,
   NULL, 0, NOW(6), NOW(6));

INSERT INTO `member_simple_id_seq` (`id`, `current_value`) VALUES (1, 1);

-- 可直接登入的本機 demo 會員：
-- 國碼 86、手機 13800000000、登入密碼 demo1234、安全密碼 123456。
INSERT INTO `member`
  (`id`, `level_code`, `status`, `last_login_time`, `failed_attempt_count`, `remark`,
   `admin_id`, `added_time`, `updated_time`)
VALUES
  (900000000000000001, 1, 1, NULL, 0, '本機 demo 會員',
   NULL, '2026-01-01 00:00:00.000000', '2026-01-01 00:00:00.000000');

INSERT INTO `member_mapping`
  (`id`, `country_code`, `mobile`, `name`, `qq`, `passwd`, `salt`, `simplify_id`,
   `security_passwd`, `wallet_address`, `added_time`, `updated_time`)
VALUES
  (900000000000000001, '86', '13800000000', 'Demo 會員', '10000001',
   '$argon2id$v=19$m=65536,t=3,p=1$bZKe+CbrDUbT7ENJcfprQw$qlJRouEEtPyxV32n16JNofUpR2NIZswoyHM6eS38PCo',
	'demo-member-salt', '00001',
	'$argon2id$v=19$m=65536,t=3,p=1$N2dNFDd3NGihZNMQNNJqYg$dqM5ssygfslp1s79HREJWJhkRehe77B+P1CzUbydH0s',
   '0x0000000000000000000000000000000000000001',
   '2026-01-01 00:00:00.000000', '2026-01-01 00:00:00.000000');

INSERT INTO `wallet_member`
  (`member_id`, `balance`, `total_amount`, `amount`, `bonus`, `frozen_amount`, `sign`,
   `added_time`, `updated_time`)
VALUES
  (900000000000000001, 1000.00, 1000.00, 1000.00, 0.00, 0.00,
   '542657f0ffbc62f428d591b8a09646c307dc783643e4bbaf4608b54b079aabf04989fb7a68ce499aa4899f9b4ba1920adda426978027a5207d1061bbda565248',
   '2026-01-01 00:00:00.000000', '2026-01-01 00:00:00.000000');

-- 收銀台建立來源單的 API 已移除，因此提供一筆可供前台掃描／付款的固定資料。
-- 對應 QR 內容：
-- v1.fecLF2eHjRwK19DTCePRfOlmcrBoXN29ZU5wzuP2f3P6IgIn6hNEUDTnC9SRkRRHs49nvfYZlJOzq-h_bUbhCudJwpRYm3z_5mj-sO9qYZjY0be5UiBO_0i3HFicd-tvhXvy5KfKtgA0mPMRgjsk6xzxOuTnw4TfOAu4P5utop9u_zLKiLVwQx0zN7PFSFqbxvVTRQXd8be1xOXO2KI
INSERT INTO `scanpay_record`
  (`id`, `amount`, `status`, `brand`, `merchant_id`, `source_order_id`, `content`,
   `expired_time`, `cancel_time`, `admin_id`, `remarks`, `added_time`, `updated_time`)
VALUES
  (900000000000000101, 100.00, 0, '品牌', 'demo-merchant-order-001',
   'demo-source-order-001',
	'v1.fecLF2eHjRwK19DTCePRfOlmcrBoXN29ZU5wzuP2f3P6IgIn6hNEUDTnC9SRkRRHs49nvfYZlJOzq-h_bUbhCudJwpRYm3z_5mj-sO9qYZjY0be5UiBO_0i3HFicd-tvhXvy5KfKtgA0mPMRgjsk6xzxOuTnw4TfOAu4P5utop9u_zLKiLVwQx0zN7PFSFqbxvVTRQXd8be1xOXO2KI',
   '2099-01-01 00:00:00.000000', NULL, NULL, '本機 demo 掃碼支付',
   '2026-01-01 00:00:00.000000', '2026-01-01 00:00:00.000000');

INSERT INTO `scanpay_mapping`
  (`record_id`, `merchant`, `merchant_order_id`, `merchant_member_id`,
   `merchant_member_account`, `merchant_member_name`, `amount`, `remark`,
	`added_time`)
VALUES
  (900000000000000101, '品牌', 'demo-merchant-order-001', 'demo-member-001',
	'demo-account', 'Demo 會員', 100.00, '本機 demo 掃碼支付',
	'2026-01-01 00:00:00.000000');
