package conf

type BlacklistHandler struct{}

func (BlacklistHandler) GetLoginNames() []string {
	return zqbConf.v.GetStringSlice("blacklist.login_name")
}

func (BlacklistHandler) GetSimulatorEnable() bool {
	return zqbConf.v.GetBool("blacklist.simulator_enable")
}
