package config

import (
	"github.com/sirupsen/logrus"
)

type loggerConfig struct {
	formatter logrus.Formatter
	level     logrus.Level
}

var logConfig *loggerConfig

func InitLogger(debugMode bool, traceMode bool) {
	if logConfig != nil {
		return
	}

	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FullTimestamp = true

	logConfig = new(loggerConfig)
	logConfig.formatter = formatter
	logConfig.level = logrus.InfoLevel

	if debugMode {
		logConfig.level = logrus.DebugLevel
	}
	if traceMode {
		logConfig.level = logrus.TraceLevel
	}
}

func GetLogger(name string) *logrus.Entry {
	if logConfig == nil {
		panic("logger not initialized")
	}
	log := logrus.New()
	log.SetLevel(logConfig.level)
	log.SetFormatter(logConfig.formatter)
	return log.WithField("_logname", name)
}
