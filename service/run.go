package service

import (
	"os"
	"os/signal"

	"github.com/james730922/wallet/service/internal/thirdparty/logger"
)

type Serve int

const (
	ServeDemo Serve = iota
)

var (
	serverStrategy = map[Serve]IServer{
		ServeDemo: newServer(),
	}
)

type IServer interface {
	Run()
}

func Run(serv Serve) {
	serverInterrupt()

	obj, ok := serverStrategy[serv]
	if !ok {
		panic("server strategy not found")
	}
	obj.Run()

	select {}
}

func serverInterrupt() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill)

	go func() {
		select {
		case c := <-interrupt:
			logger.SysLog().Warnf("Server Shutdown, osSignal: %v", c)
			os.Exit(0)
		}
	}()
}
