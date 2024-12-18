package worker

import (
	"context"

	"github.com/CrazyThursdayV50/pkgo/builtin"
	gchan "github.com/CrazyThursdayV50/pkgo/builtin/chan"
	"github.com/CrazyThursdayV50/pkgo/goo"
	"github.com/CrazyThursdayV50/pkgo/log"
	defaultlogger "github.com/CrazyThursdayV50/pkgo/log/default"
	"github.com/CrazyThursdayV50/pkgo/monitor"
)

type Worker[J any] struct {
	*monitor.Monitor
	name     string
	do       func(J)
	count    int64
	trigger  builtin.Iter[int, J]
	graceful bool
	logger   log.Logger
}

func (w *Worker[J]) run(ctx context.Context) {
	switch w.graceful {
	case true:
		goo.Go(func() {
			w.trigger.Iter(func(_ int, element J) (bool, error) {
				w.do(element)
				return true, nil
			})
		})

	default:
		goo.Go(func() {
			w.trigger.Iter(func(_ int, element J) (bool, error) {
				select {
				case <-ctx.Done():
					return false, nil

				default:
					w.do(element)
					return true, nil
				}
			})
		})
	}
}

func New[J any](name string, do func(job J)) (*Worker[J], func(J)) {
	var w Worker[J]
	w.name = name
	w.do = func(j J) {
		do(j)
		w.count++
	}
	var trigger = gchan.Make[J](0)
	w.trigger = trigger
	w.WithContext(context.TODO())
	w.logger = defaultlogger.New(defaultlogger.DefaultConfig())
	w.logger.Init()
	return &w, func(j J) {
		trigger.Send(j)
	}
}

func (w *Worker[J]) WithContext(ctx context.Context) {
	w.Monitor = monitor.New(ctx, w.name)
}

func (w *Worker[J]) WithGraceful(ok bool) {
	w.graceful = ok
}

func (w *Worker[J]) WithTrigger(trigger <-chan J) {
	w.trigger = gchan.FromRead(trigger)
}

func (w *Worker[J]) WithLogger(logger log.Logger) {
	w.logger = logger
}

func (w *Worker[J]) Run() {
	w.Monitor.SetLogger(w.logger)
	if w.graceful {
		w.Monitor.OnExit(func() {
			w.logger.Info("Worker stop gracefully, wait for all jobs done and trigger closed.")
		})
	}
	w.Monitor.Run(w.run)
}

func (m *Worker[J]) Stop() {
	m.Monitor.Stop()
}

func (m *Worker[J]) Count() int64 { return m.count }
