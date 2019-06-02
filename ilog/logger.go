package ilog

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/kgrunwald/goweb/di"
)

type Fields log.Fields

type Logger interface {
	Debug(string)
	Info(string)
	Warn(string)
	Error(string)
	Fatal(string)
	WithField(key string, value interface{}) Logger
	WithFields(values ...interface{}) Logger
}

func init() {
	c := di.GetContainer()
	c.Register(NewLogger)
}

func NewLogger() Logger {
	log.SetHandler(cli.New(os.Stdout))
	log.SetLevel(log.DebugLevel)
	return &logger{log.WithFields(log.Fields{})}
}

type logger struct {
	log *log.Entry
}

func (l *logger) Debug(msg string) {
	l.log.Debug(msg)
}

func (l *logger) Info(msg string) {
	l.log.Info(msg)
}

func (l *logger) Warn(msg string) {
	l.log.Warn(msg)
}

func (l *logger) Error(msg string) {
	l.log.Error(msg)
}

func (l *logger) Fatal(msg string) {
	l.log.Fatal(msg)
}

func (l *logger) WithField(key string, value interface{}) Logger {
	return &logger{l.log.WithField(key, value)}
}

func (l *logger) WithFields(fields ...interface{}) Logger {
	f := log.Fields{}
	for i := 0; i < len(fields); i += 2 {
		f[fields[i].(string)] = fields[i+1]
	}

	return &logger{l.log.WithFields(f)}
}
