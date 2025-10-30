package sugar

import (
	"errors"
	"syscall"

	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/CrazyThursdayV50/pkgo/log/zap"
)

var _ log.Logger = (*apiLogger)(nil)

type apiLogger struct {
	sugar log.Logger
}

func New(cfg *Config) *apiLogger {
	sugar := zap.New(cfg).Sugar()

	if err := sugar.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		sugar.Error(err)
		return nil
	}

	return &apiLogger{sugar: sugar}
}

// Logger methods
func (l *apiLogger) Debug(args ...interface{}) {
	l.sugar.Debug(args...)
}

func (l *apiLogger) Debugf(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}

func (l *apiLogger) Info(args ...interface{}) {
	l.sugar.Info(args...)
}

func (l *apiLogger) Infof(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}

func (l *apiLogger) Warn(args ...interface{}) {
	l.sugar.Warn(args...)
}

func (l *apiLogger) Warnf(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}

func (l *apiLogger) Error(args ...interface{}) {
	l.sugar.Error(args...)
}

func (l *apiLogger) Errorf(template string, args ...interface{}) {
	l.sugar.Errorf(template, args...)
}

func (l *apiLogger) Printf(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}
