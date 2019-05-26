package ilog

import (
	"github.com/kgrunwald/goweb/di"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})
}

func NewLogger() Logger {
	logrus.SetLevel(logrus.DebugLevel)
	return logrus.WithFields(logrus.Fields{"channel": "app"})
}

func init() {
	c := di.GetContainer()
	c.Register(NewLogger)
}
