# 建置與本機部署

```bash
GOFLAGS=-mod=vendor go test ./...
GOFLAGS=-mod=vendor go build ./...
docker build -t zqb-apis:demo .
```

此 demo 不包含既有公司 registry、Jenkins、DEV、UAT 或正式環境部署資訊。執行服務前請先依[快速啟動](quick_start.md)啟動本機資料庫。
