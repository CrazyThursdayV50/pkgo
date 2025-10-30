package sugar

import (
	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/CrazyThursdayV50/pkgo/log/zap"
)

var _ log.FormatterLogger[any] = (*apiLogger)(nil)
var _ log.SimpleLogger[any] = (*apiLogger)(nil)

type apiLogger struct {
	sugar log.Logger
}

func New(cfg *Config) *apiLogger {
	sugar := zap.New(cfg).Sugar()

	// if err := sugar.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
	// 	sugar.Warnf("sugarLogger sync failed: %v", err)
	// }

	return &apiLogger{sugar: sugar}
}

// Logger methods
func (l *apiLogger) Debug(args ...any) {
	l.sugar.Debug(args...)
}

func (l *apiLogger) Debugf(template string, args ...any) {
	l.sugar.Debugf(template, args...)
}

func (l *apiLogger) Info(args ...any) {
	l.sugar.Info(args...)
}

func (l *apiLogger) Infof(template string, args ...any) {
	l.sugar.Infof(template, args...)
}

func (l *apiLogger) Warn(args ...any) {
	l.sugar.Warn(args...)
}

func (l *apiLogger) Warnf(template string, args ...any) {
	l.sugar.Warnf(template, args...)
}

func (l *apiLogger) Error(args ...any) {
	l.sugar.Error(args...)
}

func (l *apiLogger) Errorf(template string, args ...any) {
	l.sugar.Errorf(template, args...)
}

func (l *apiLogger) Printf(template string, args ...any) {
	l.sugar.Infof(template, args...)
}
