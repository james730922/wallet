package conf

type DataBaseHandler struct{}

// 預設過期秒數
func (DataBaseHandler) GetSlaveSwitch() bool {
	return zqbConf.v.GetBool("database.slave_switch")
}
