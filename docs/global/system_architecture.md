# Demo 系統架構

服務只提供會員 API。MySQL 是唯一帳務真相來源，保存前台業務資料、交易流水與錢包異動日誌；Redis 只保存會員 token 與快取。錢包更新、交易流水與 `wallet_journal` 使用同一個 MySQL transaction。

簡化會員 ID 由 MySQL `member_simple_id_seq` 原子配發，使用短 transaction 與 `LAST_INSERT_ID()` 保證多個 API instance 不會取得相同編號。會員註冊失敗時允許序號跳號。

本機模式不連接既有 DEV、UAT 或正式站。API 只監聽 `127.0.0.1`，MySQL 與 Redis 也全部使用 localhost；檔案服務只保留 localhost 預留位址。

服務不啟動後台、收銀台或背景排程程序。
