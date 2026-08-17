package conf

type DBHandler struct{}

// 預設過期秒數
func (DBHandler) GetRawSQLRoot() string {
	return appConf.v.GetString("db.raw_sql_root")
}

func (DBHandler) GetRawSQLUrl() string {
	return appConf.v.GetString("db.sql_url")
}

func (DBHandler) GetSQLSlaveUrl() string {
	return appConf.v.GetString("db.sql_slave_url")
}

func (DBHandler) GetMockSQLUrl() string {
	return appConf.v.GetString("db.mock_sql_url")
}

func (DBHandler) GetGormLogMode() bool {
	return appConf.v.GetBool("db.gorm_log_mode")
}
