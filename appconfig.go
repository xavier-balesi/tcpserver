package main

import "flag"

// Application Config struct
type AppConfig struct {
	Host      string
	Port      int
	Protocol  string
	DebugMode bool
	TraceMode bool
}

// Default values for Application Config
var DefaultConfig = &AppConfig{
	Host:      "127.0.0.1",
	Port:      1234,
	Protocol:  "tcp4",
	DebugMode: false,
	TraceMode: false,
}

// Parse Application Config merged with default values
func GetAppConfig() *AppConfig {
	var commandLineConfig *AppConfig = new(AppConfig)

	flag.StringVar(&commandLineConfig.Host, "host", DefaultConfig.Host, "host to expose")
	flag.IntVar(&commandLineConfig.Port, "port", DefaultConfig.Port, "port to expose")
	flag.StringVar(&commandLineConfig.Protocol, "protocol", DefaultConfig.Protocol, "network protocol")
	flag.BoolVar(&commandLineConfig.DebugMode, "debug", DefaultConfig.DebugMode, "debug mode")
	flag.BoolVar(&commandLineConfig.TraceMode, "trace", DefaultConfig.TraceMode, "trace mode")
	flag.Parse()

	return commandLineConfig
}
