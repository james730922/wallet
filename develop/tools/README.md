# 本機資料服務

從專案根目錄啟動：

```bash
docker compose -f develop/tools/docker-compose.yaml up -d
```

此 Compose 只使用 localhost 對外映射：MySQL `3306`、Redis `6379`。MySQL 首次建立 volume 時會依序載入 `develop/tools/sql/DDL_schema.sql` 與 `develop/tools/sql/DML_default_value.sql`。

若初始化 SQL 有變更，請移除專案專用 volume 後重建；這會清除該 demo volume 中的本機資料：

```bash
docker compose -f develop/tools/docker-compose.yaml down -v
docker compose -f develop/tools/docker-compose.yaml up -d
```
