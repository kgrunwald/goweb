package ilog

import (
	"encoding/json"
	"os"
	"sync"

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

var globalLogger = &logger{log.WithFields(log.Fields{})}

func Debug(msg string) {
	globalLogger.Debug(msg)
}

func Info(msg string) {
	globalLogger.Info(msg)
}

func Warn(msg string) {
	globalLogger.Warn(msg)
}

func Error(msg string) {
	globalLogger.Error(msg)
}

func Fatal(msg string) {
	globalLogger.Fatal(msg)
}

func WithError(err error) Logger {
	return globalLogger.WithError(err)
}

func WithField(key string, value interface{}) Logger {
	return globalLogger.WithField(key, value)
}

func WithFields(fields ...interface{}) Logger {
	return globalLogger.WithFields(fields...)
}

type Handler struct {
	*json.Encoder
	mu sync.Mutex
}

func (h *Handler) HandleLog(e *log.Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	m := map[string]interface{}{
		"message": e.Message,
		"level":   e.Level,
	}
	for _, f := range e.Fields.Names() {
		m[f] = e.Fields.Get(f)
	}
	return h.Encoder.Encode(m)
}

func init() {
	if os.Getenv("LOG_CLI") != "" {
		log.SetHandler(cli.New(os.Stdout))
	} else {
		log.SetHandler(&Handler{
			Encoder: json.NewEncoder(os.Stderr),
		})
	}

	switch os.Getenv("LOG_LEVEL") {
	case "FATAL":
		log.SetLevel(log.FatalLevel)
		break
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
		break
	case "WARN":
		log.SetLevel(log.WarnLevel)
		break
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
		break
	default:
		log.SetLevel(log.InfoLevel)
	}

	c := di.GetContainer()
	c.Register(NewLogger)
}

func NewLogger() Logger {
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

func (l *logger) WithError(err error) Logger {
	return &logger{l.log.WithError(err)}
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
