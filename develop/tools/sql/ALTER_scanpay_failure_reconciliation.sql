-- 供已存在的環境啟用 ScanPay Failure/reconciliation 使用。
-- 新環境直接使用 DDL_schema.sql，不需要再執行本檔。

ALTER TABLE `scanpay_record`
  MODIFY COLUMN `status` tinyint(4) NOT NULL DEFAULT '0'
    COMMENT '0=待確認 1=交易中 2=完成 3=取消 4=付款失敗作廢';

ALTER TABLE `order_scanpay`
  ADD KEY `idx_order_scanpay_status_updated` (`status`, `updated_time`);
