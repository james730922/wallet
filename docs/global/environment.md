# 本機 Demo 環境

| 服務 | 位址 | Demo 認證 |
| --- | --- | --- |
| 前台 API | `127.0.0.1:17801` | 會員登入後取得 token |
| MySQL | `127.0.0.1:3306` | `root` / `demo_password` |
| Redis | `127.0.0.1:6379` | 無密碼 |
| 檔案服務預留位址 | `127.0.0.1:18080` | `demo` / `demo_password` |

預設前台會員：

| 欄位 | Demo 值 |
| --- | --- |
| 國碼 | `86` |
| 手機 | `13800000000` |
| 登入密碼 | `demo1234` |
| 安全密碼 | `123456` |
| 初始錢包餘額 | `1000` |

預設掃碼付款資料的金額為 `100`，QR Code 內容記錄在 `develop/tools/sql/DML_default_value.sql`。這筆資料只供一次成功付款；需要重測時請重建 demo MySQL volume。

所有值僅供本機展示，不應直接用於正式環境。

掃碼 AES-256 key 的本機 demo 值位於 `conf.d/app.conf`。正式環境請以 `WALLET_SCANPAY_KEY` 提供 Base64 編碼的 32-byte key；更換 key 時同步設定 `WALLET_SCANPAY_KEY_VERSION`，舊版本 QR 不會被新版本接受。
