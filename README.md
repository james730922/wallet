# 品牌交易 API Demo

這是本機展示版本，只保留前台會員註冊／登入、存款、交易紀錄與會員掃碼支付流程。


## 本機啟動

需求：Go、Docker 與 Docker Compose。

```bash
docker compose -f develop/tools/docker-compose.yaml up -d
GOFLAGS=-mod=vendor go run ./cmd/zqbapis
```

本機連接埠：

- 前台 API：`http://127.0.0.1:17801`
- MySQL：`127.0.0.1:3306`
- Redis：`127.0.0.1:6379`

健康與可觀測性端點：

- Liveness：`GET /livez`
- Readiness：`GET /readyz`（實際檢查 primary DB 與 Redis）
- Prometheus metrics：`GET /metrics`

每個 HTTP response 會回傳 `X-Trace-ID`；如需匯出 OpenTelemetry traces，請設定
`OTEL_EXPORTER_OTLP_TRACES_ENDPOINT`（例如 `http://otel-collector:4318/v1/traces`）。

預設值都只供本機 demo 使用，設定集中於 `conf.d/app.conf` 與 `conf.d/zqb.yaml`。外部檔案服務預設指向 localhost；未啟動時，只有呼叫相依功能才會受影響。

MySQL 初始化後可直接使用以下會員資料測試前台：

- 國碼：`86`
- 手機：`13800000000`
- 登入密碼：`demo1234`
- 安全密碼：`123456`

存款方式與一筆 100 元掃碼付款資料也會一併建立，詳細測試值請見[本機環境](docs/global/environment.md)。

## 驗證

```bash
GOFLAGS=-mod=vendor go test ./...
GOFLAGS=-mod=vendor go build ./...
```

## 文件

- [前台 API](docs/app_apis_doc.md)
- [快速啟動](docs/global/quick_start.md)
- [本機環境](docs/global/environment.md)
- [專案結構](docs/global/project_structure.md)
- [掃碼支付](docs/scanpay/scanpay.md)
