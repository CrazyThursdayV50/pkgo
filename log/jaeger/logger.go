package jaegerlogger

import (
	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/uber/jaeger-client-go"
)

type jaegerLogger struct{ log.FormatterLogger[any] }

func (l *jaegerLogger) Error(msg string) {
	l.FormatterLogger.Errorf(msg)
}

func (l *jaegerLogger) Infof(msg string, args ...any) {
	l.FormatterLogger.Infof(msg, args...)
}

func New(logger log.FormatterLogger[any]) jaeger.Logger {
	return &jaegerLogger{logger}
}
