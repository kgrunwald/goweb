package ilog

import (
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

type logger struct {
	Channel string
}

func (l *logger) Debug(args ...interface{}) {
	logrus.WithField("channel", l.Channel).Debug(args...)
}

func (l *logger) Info(args ...interface{}) {
	logrus.WithField("channel", l.Channel).Info(args...)
}

func (l *logger) Warn(args ...interface{}) {
	logrus.WithField("channel", l.Channel).Warn(args...)
}

func (l *logger) Error(args ...interface{}) {
	logrus.WithField("channel", l.Channel).Error(args...)
}

func (l *logger) Fatal(args ...interface{}) {
	logrus.WithField("channel", l.Channel).Fatal(args...)
}
