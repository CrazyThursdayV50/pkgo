package cron

import (
	"context"
	"time"

	"github.com/CrazyThursdayV50/pkgo/builtin/slice"
	"github.com/CrazyThursdayV50/pkgo/goo"
	"github.com/CrazyThursdayV50/pkgo/log"
	defaultlogger "github.com/CrazyThursdayV50/pkgo/log/default"
	"github.com/CrazyThursdayV50/pkgo/worker"
)

type Cron struct {
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}

	job          func()
	runAfter     time.Duration
	waitAfterRun bool
	tick         time.Duration
	trigger      func(context.Context)
	worker       *worker.Worker[context.Context]
	runOnStart   func()
	logger       log.Logger
	pauseChan    *chan struct{}
	restoreChan  *chan struct{}
}

func defaultOptions() []Option {
	logger := defaultlogger.New(defaultlogger.DefaultConfig())
	logger.Init()
	return []Option{
		WithJob(func() {}, time.Minute),
		WithRunAfterStart(-1),
		WithWaitAfterRun(false),
		WithLogger(logger),
	}
}

func timerDo(duration time.Duration, done <-chan struct{}, do func()) {
	timer := time.NewTimer(duration)
	select {
	case <-done:
		return
	case <-timer.C:
		do()
	}
}

func (c *Cron) restore() {
	select {
	case <-*c.pauseChan:
		pause := make(chan struct{})
		c.pauseChan = &pause
		close(*c.restoreChan)

	case <-*c.restoreChan:
	}
}

func (c *Cron) pause() {
	select {
	case <-*c.pauseChan:
	case <-*c.restoreChan:
		restore := make(chan struct{})
		c.restoreChan = &restore
		close(*c.pauseChan)
	}
}

func (c *Cron) init() {
	if c.tick < 1 {
		c.tick = time.Second
	}

	pause := make(chan struct{})
	c.pauseChan = &pause

	restore := make(chan struct{})
	c.restoreChan = &restore
	close(*c.restoreChan)

	c.done = make(chan struct{})
	c.worker, c.trigger = worker.New("Cron", func(context.Context) { c.job() })
	c.worker.WithGraceful(false)
	c.worker.WithLogger(c.logger)

	if c.runAfter < 0 {
		return
	}

	do := func() { c.trigger(c.ctx) }
	if c.runAfter == 0 {
		c.runOnStart = do
		return
	}

	c.runOnStart = func() { timerDo(c.runAfter, c.done, do) }

	goo.Go(func() {
		<-c.ctx.Done()
		close(c.done)
	})
}

func (c *Cron) tickerRun() {
	ticker := time.NewTicker(c.tick)
	goo.Go(func() {
		defer c.worker.Stop()
		for {
			select {
			case <-*c.pauseChan:
				select {
				case <-*c.restoreChan:
				case <-c.done:
					return
				}

			case <-c.done:
				return

			case _ = <-ticker.C:
				select {
				case <-c.done:
				default:
					c.trigger(c.ctx)
				}
			}
		}
	})
}

func (c *Cron) timerRun() {
	timer := time.NewTimer(c.tick)
	goo.Go(func() {
		defer c.worker.Stop()
		for {
			select {
			case <-*c.pauseChan:
				select {
				case <-*c.restoreChan:
				case <-c.done:
					return
				}

			case <-c.done:
				return

			case _ = <-timer.C:
				select {
				case <-c.done:
				default:
					c.trigger(c.ctx)
					timer.Stop()
					timer.Reset(c.tick)
				}
			}
		}
	})
}

func New(opts ...Option) *Cron {
	var c Cron
	opts = append(defaultOptions(), opts...)
	_, _ = slice.From(opts...).Iter(func(_ int, opt Option) (bool, error) { opt(&c); return true, nil })
	c.init()
	return &c
}

func (c *Cron) Run(ctx context.Context) {
	c.ctx, c.cancel = context.WithCancel(ctx)
	c.worker.Run(c.ctx)
	c.runOnStart()
	if c.waitAfterRun {
		c.timerRun()
	} else {
		c.tickerRun()
	}
}

func (c *Cron) Stop()                 { c.cancel() }
func (c *Cron) Pause()                { c.pause() }
func (c *Cron) Restore()              { c.restore() }
func (c *Cron) OnStart(f func())      { c.worker.OnStart(f) }
func (c *Cron) OnExit(f func())       { c.worker.OnExit(f) }
func (c *Cron) Done() <-chan struct{} { return c.done }
