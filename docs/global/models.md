# Models

- `condition`：HTTP request、內部條件與資料庫 query／update 條件。
- `model`：資料表與核心業務資料結構。
- `resp`：對外 API response 與匯出資料結構。

常見資料流：`Req -> Cond -> Query/Update -> model -> Resp`。
