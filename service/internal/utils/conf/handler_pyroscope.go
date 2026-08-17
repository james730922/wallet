package conf

type PyroscopeHandler struct{}

// pprof 資料收集：啟用/禁用
func (PyroscopeHandler) GetEnable() bool {
	return zqbConf.v.GetBool("pyroscope.enable")
}

// pprof 資料收集：位置
func (PyroscopeHandler) GetHost() string {
	return zqbConf.v.GetString("pyroscope.host")
}

// pprof 資料收集：本服務名稱
func (PyroscopeHandler) GetApplicationName() string {
	return zqbConf.v.GetString("pyroscope.application_name")
}
