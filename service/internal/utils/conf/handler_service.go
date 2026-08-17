package conf

import "net"

type IService interface {
	GetApisHTTPAddress() string
}

type ServiceHandler struct{}

func (ServiceHandler) GetApisHTTPAddress() string {
	host := zqbConf.v.GetString("service.host")
	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		panic("service.host must be a loopback IP address")
	}
	return net.JoinHostPort(host, zqbConf.v.GetString("service.port.apis_http"))
}

func (ServiceHandler) GetTrustedProxies() []string {
	return zqbConf.v.GetStringSlice("service.trusted_proxies")
}
