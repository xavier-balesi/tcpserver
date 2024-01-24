package main

import (
	// "time"

	"tcpserver/config"
	"tcpserver/server"
)

func main() {
	appConfig := GetAppConfig()

	serverInfo := server.ServerInfo{
		Host: appConfig.Host, Port: appConfig.Port, Protocol: appConfig.Protocol,
	}

	config.InitLogger(appConfig.DebugMode, appConfig.TraceMode)
	log := config.GetLogger("main")

	log.
		WithField("serverInfo", serverInfo).
		WithField("debugMode", appConfig.DebugMode).
		WithField("traceMode", appConfig.TraceMode).
		Debug("input flags")

	s := server.New(serverInfo)
	go s.Start()

	for {
	}
}
