# 專案結構

- `cmd/zqbapis`：demo 服務執行入口。
- `conf.d`：本機服務與第三方依賴設定。
- `develop/tools`：MySQL 與 Redis 的本機 Docker Compose。
- `docker`：應用容器範例。
- `docs`：目前保留功能的文件。
- `service/internal/app`：前台 API 入口。
- `service/internal/controller`：前台 HTTP controller。
- `service/internal/core/base`：登入、存款、交易、掃碼與錢包核心邏輯。
- `service/internal/core/usecase`：前台 use case。
- `service/internal/models`：資料庫、請求及回應模型。
- `service/internal/thirdparty`：MySQL、Redis 與檔案服務介接。
- `service/internal/utils`：設定、錯誤、簽章及共用工具。
- `develop/tools/sql/DDL_schema.sql`：前台功能所需的本機資料表。
- `develop/tools/sql/DML_default_value.sql`：本機 demo 初始資料。
- `vendor`：固定版本的第三方依賴。
