package log

type Logger interface {
	SimpleLogger[any]
	FormatterLogger[any]
}

type FormatterLogger[Arg any] interface {
	Debugf(fmt string, args ...Arg)
	Infof(fmt string, args ...Arg)
	Warnf(fmt string, args ...Arg)
	Errorf(fmt string, args ...Arg)
}

type SimpleLogger[Arg any] interface {
	Debug(args ...Arg)
	Info(args ...Arg)
	Warn(args ...Arg)
	Error(args ...Arg)
}

type DescLogger[Arg any] interface {
	Debug(desc string, args ...Arg)
	Info(desc string, args ...Arg)
	Warn(desc string, args ...Arg)
	Error(desc string, args ...Arg)
}
