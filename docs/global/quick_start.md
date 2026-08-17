# 快速啟動

從專案根目錄執行：

```bash
docker compose -f develop/tools/docker-compose.yaml up -d
GOFLAGS=-mod=vendor go run ./cmd/zqbapis
```

Docker Compose 會啟動 MySQL 與 Redis，並由 MySQL 自動載入 demo schema 和 seed。前台 API 啟動後監聽 `127.0.0.1:17801`。

預設會員可使用國碼 `86`、手機 `13800000000`、登入密碼 `demo1234` 登入；需要付款密碼的流程使用 `123456`。

若曾用舊版 schema 建立過 MySQL volume，請先執行以下指令清除 demo 資料並重新初始化：

```bash
docker compose -f develop/tools/docker-compose.yaml down -v
docker compose -f develop/tools/docker-compose.yaml up -d
```

停止依賴服務：

```bash
docker compose -f develop/tools/docker-compose.yaml down
```
