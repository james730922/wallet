package conf

type LogHandler struct{}

// 系統log
func (LogHandler) GetSysLog() string {
	return zqbConf.v.GetString("log.sys_log")
}

// 服務log
func (LogHandler) GetApLog() string {
	return zqbConf.v.GetString("log.ap_log")
}

// 訪問log
func (LogHandler) GetAccessLog() string {
	return zqbConf.v.GetString("log.access_log")
}

// Gin log 啟用狀態
func (LogHandler) GetGinLogEnable() bool {
	return zqbConf.v.GetBool("log.gin_log_enable")
}
