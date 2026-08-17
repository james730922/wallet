# Sandbox

## Introduction
sandbox 保留單元測試樣板；本機 MySQL 與 Redis 統一由 `develop/tools/docker-compose.yaml` 啟動。

### 目錄結構
<pre>
|-- sandbox
    |-- template unit test template
        |-- template_test.go
    |-- sandbox_test.go
</pre>

### Unit Test
- 複製 ./template/template_test.go 到你想做單元測試的檔案同一個位置
- 按照template完成你的測試

### 初始化 SQL

- Schema：`develop/tools/sql/DDL_schema.sql`
- Demo 預設資料：`develop/tools/sql/DML_default_value.sql`
- SQL 異動後須重建 demo MySQL volume，操作方式見 `develop/tools/README.md`。
