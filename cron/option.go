package cron

import (
	"time"

	"github.com/CrazyThursdayV50/pkgo/log"
)

type Option func(*Cron)

func WithJob(job func(), tick time.Duration) Option {
	return func(c *Cron) {
		c.job = job
		c.tick = tick
	}
}

func WithRunAfterStart(duration time.Duration) Option {
	return func(c *Cron) {
		c.runAfter = duration
	}
}

func WithWaitAfterRun(ok bool) Option {
	return func(c *Cron) {
		c.waitAfterRun = ok
	}
}

func WithLogger(logger log.Logger) Option {
	return func(c *Cron) {
		c.logger = logger
	}
}
