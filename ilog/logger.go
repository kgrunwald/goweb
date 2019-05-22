package ilog

import "github.com/sirupsen/logrus"

type Logger interface {
	Debug(...interface{})
	Info(...interface{})
}

func NewLogger() Logger {
	return logrus.WithFields(logrus.Fields{"channel": "app"})
}
