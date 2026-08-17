# 前台 API（Demo）

Base URL：`http://127.0.0.1:17801/api`

除註冊、登入、忘記密碼與驗證碼註冊外，其餘 API 都需要會員 token。實際 request／response 結構以 controller 與 `service/internal/models` 為準。

## 匿名 API

| Method | Path | 用途 |
| --- | --- | --- |
| POST | `/v1/registration` | 會員註冊 |
| POST | `/v1/login` | 密碼登入 |
| POST | `/v1/captcha/register` | 註冊圖形驗證 |

## 會員 API

| Method | Path | 用途 |
| --- | --- | --- |
| POST | `/v1/logout` | 登出 |
| POST | `/v1/update-pwd` | 修改登入密碼 |
| GET | `/v1/security-passwd/first` | 查詢是否首次設定安全密碼 |
| POST | `/v1/security-passwd/forget` | 忘記安全密碼 |
| POST | `/v1/security-passwd/update` | 修改安全密碼 |
| POST | `/v1/security-passwd/identifier` | 驗證安全密碼 |
| GET | `/v1/banks` | 取得存款流程可用銀行 |
| GET | `/v1/deposit` | 存款方式列表 |
| GET | `/v1/deposit/methods` | 存款方式選項 |
| POST | `/v1/deposit/order` | 建立存款單 |
| GET | `/v1/transaction` | 交易紀錄列表 |
| GET | `/v1/transaction/detail/:id` | 交易紀錄詳情 |
| POST | `/v1/scan-pay/scan` | 掃描支付碼 |
| POST | `/v1/scan-pay/pay` | 完成掃碼支付 |
| POST | `/v1/scan-pay/qrcode` | 產生掃碼支付碼 |
