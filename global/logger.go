package global

import (
	"github.com/Sirupsen/logrus"
	"os"
)

var Runlogger *logrus.Logger

func InitLogger() {
	runlog, _ := os.OpenFile(GlobalConfig.LogCfg.RunLog, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

	level, _ := logrus.ParseLevel("info")
	Runlogger = &logrus.Logger{
		Out:       runlog,
		Level:     level,
		Formatter: new(logrus.JSONFormatter),
	}
	return
}
