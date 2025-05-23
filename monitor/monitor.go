package monitor

import (
	"context"

	"github.com/CrazyThursdayV50/pkgo/goo"
	"github.com/CrazyThursdayV50/pkgo/log"
	defaultlogger "github.com/CrazyThursdayV50/pkgo/log/default"
)

type Monitor struct {
	ctx     context.Context
	cancel  context.CancelFunc
	logger  log.Logger
	onStart func()
	onExit  func()
	run     func(context.Context)
}

func New(ctx context.Context, name string) *Monitor {
	var s Monitor
	s.ctx, s.cancel = context.WithCancel(ctx)
	s.logger = defaultlogger.New(defaultlogger.DefaultConfig())
	s.logger.Init()
	if name == "" {
		name = "Monitor"
	}
	s.onStart = func() { s.logger.Debugf("%s start", name) }
	s.onExit = func() { s.logger.Debugf("%s exit", name) }
	return &s
}

func (s *Monitor) Run(f func(context.Context)) {
	if s.onStart != nil {
		s.onStart()
	}

	goo.Goo(func() {
		<-s.ctx.Done()

		if s.onExit == nil {
			return
		}

		s.onExit()
	}, func(err error) {
		s.logger.Errorf("exit panic: %v", err)
	})

	goo.Goo(func() {
		f(s.ctx)
	}, func(err error) { s.logger.Errorf("monitor panic: %v", err) })
}

func wrap(next, f func()) func() {
	return func() {
		defer func() {
			if next == nil {
				return
			}
			next()
		}()

		if f != nil {
			f()
		}
	}
}

func (s *Monitor) Stop()                       { s.cancel() }
func (s *Monitor) Done() <-chan struct{}       { return s.ctx.Done() }
func (s *Monitor) OnStart(f func())            { s.onStart = wrap(f, s.onStart) }
func (s *Monitor) OnExit(f func())             { s.onExit = wrap(s.onExit, f) }
func (s *Monitor) SetLogger(logger log.Logger) { s.logger = logger }
