package app

import (
	"github.com/pyroscope-io/pyroscope/pkg/agent/profiler"

	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/conf"
)

func NewPyroscopeServe() IPyroscopeServe {
	return &pyroscopeServe{}
}

type IPyroscopeServe interface {
	Run()
}

// 收集服務資訊，將類似 pprof 的資訊傳送到 pyroscope
type pyroscopeServe struct {
}

func (c *pyroscopeServe) Run() {
	if !conf.Pyroscope().GetEnable() {
		return
	}

	logger.SysLog().Info("serve start [pyroscope]")

	if _, err := profiler.Start(profiler.Config{
		ApplicationName: conf.Pyroscope().GetApplicationName(),
		ServerAddress:   conf.Pyroscope().GetHost(),
	}); err != nil {
		logger.SysLog().Errorf("[pyroscope] start err: %s", err)
	}
}
