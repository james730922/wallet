# Demo 安全說明

- 會員密碼與安全密碼使用 Argon2id；既有 SHA-512 密碼在成功驗證後自動升級。
- 會員登入後以 Redis token 驗證受保護 API；token 使用固定到期時間，登入與修改密碼時會原子輪替並撤銷舊 token。
- 匿名登入密碼重設入口已移除；修改登入密碼必須提供目前密碼。
- 掃碼資料使用 AES-256-GCM、隨機 nonce 與版本化格式；竄改內容會驗證失敗。
- SQL 操作沿用參數化查詢與交易處理。
- 前台回應使用獨立 response model，不直接暴露資料庫 model。
- demo 設定只允許 localhost 依賴位址，並移除既有環境 IP、帳密及服務連結。
- 即使關閉 Captcha，匿名註冊仍會使用 Redis 進行 IP 與 device 雙層限流；client 可透過 `X-Device-ID` 提供穩定裝置識別值。
- 外部檔案服務的 login、upload、delete、list 與 health request 均有 deadline 與 transport timeout。

本機 API 只監聽 `127.0.0.1`，且未掛載 `/debug/pprof`。正式部署前必須替換所有 demo 密碼／金鑰（掃碼 key 可用 `WALLET_SCANPAY_KEY` 注入）、使用 TLS、限制資料庫網路範圍，並建立獨立的秘密管理機制。Reverse proxy 必須移除 client 偽造的 `X-Forwarded-For`，再寫入真實來源 IP。
