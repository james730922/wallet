package conf

type SignHandler struct{}

// 預設過期秒數
func (SignHandler) GetSalt() []byte {
	return []byte(zqbConf.v.GetString("sign.salt"))
}
