package unimpl

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/CrazyThursdayV50/pkgo/log"
)

type UnimplementLogger struct{}

var _ log.FormatterLogger[any] = (*UnimplementLogger)(nil)
var _ log.SimpleLogger[any] = (*UnimplementLogger)(nil)
var _ log.Logger = (*UnimplementLogger)(nil)

func New() *UnimplementLogger {
	return &UnimplementLogger{}
}

func prefixArgs(args ...any) []any {
	pc, _, _, _ := runtime.Caller(1)
	callerName := runtime.FuncForPC(pc).Name()
	index := strings.LastIndex(callerName, ".")
	return append([]any{fmt.Sprintf("[UNIMPLEMENT %s]", callerName[index+1:])}, args...)
}

func prefixFormat(format string) string {
	pc, _, _, _ := runtime.Caller(1)
	callerName := runtime.FuncForPC(pc).Name()
	index := strings.LastIndex(callerName, ".")

	return fmt.Sprintf("[UNIMPLEMENT %s]%s", callerName[index+1:], format)
}

// implement SimpleLogger
func (l *UnimplementLogger) Debug(args ...any) {
	fmt.Println(prefixArgs(args...)...)
}

func (l *UnimplementLogger) Info(args ...any) {
	fmt.Println(prefixArgs(args...)...)
}

func (l *UnimplementLogger) Warn(args ...any) {
	fmt.Println(prefixArgs(args...)...)
}

func (l *UnimplementLogger) Error(args ...any) {
	fmt.Println(prefixArgs(args...)...)
}

// implement FormatterLogger
func (l *UnimplementLogger) Debugf(format string, args ...any) {
	fmt.Printf(prefixFormat(format), args...)
}

func (l *UnimplementLogger) Infof(format string, args ...any) {
	fmt.Printf(prefixFormat(format), args...)
}

func (l *UnimplementLogger) Warnf(format string, args ...any) {
	fmt.Printf(prefixFormat(format), args...)
}

func (l *UnimplementLogger) Errorf(format string, args ...any) {
	fmt.Printf(prefixFormat(format), args...)
}
