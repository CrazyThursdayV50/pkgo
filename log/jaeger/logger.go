package jaegerlogger

import (
	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/uber/jaeger-client-go"
)

type jaegerLogger struct{ log.Logger }

func (l *jaegerLogger) Error(msg string) {
	l.Logger.Errorf(msg)
}

func (l *jaegerLogger) Infof(msg string, args ...any) {
	l.Logger.Infof(msg, args...)
}

func New(logger log.Logger) jaeger.Logger {
	return &jaegerLogger{logger}
}
