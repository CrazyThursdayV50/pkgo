package impl

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/CrazyThursdayV50/gotils/pkg/wrapper"
	"github.com/CrazyThursdayV50/gotils/pkg/wrapper/wrap"
)

type (
	ChanRead[E any] interface {
		chan E | <-chan E
	}

	ChanR[E any] struct {
		l           *sync.Mutex
		done        chan struct{}
		recvTimeout time.Duration
		sendTimeout time.Duration
		c           <-chan E
		count       int64
	}
)

func From[E any, C ChanRead[E]](c C) *ChanR[E] {
	return &ChanR[E]{
		l:    &sync.Mutex{},
		done: make(chan struct{}),
		c:    c,
	}
}

func (c *ChanR[E]) Len() int {
	if c == nil {
		return 0
	}
	return len(c.c)
}

func (c *ChanR[E]) IsEmpty() bool { return c.Len() == 0 }

func (c *ChanR[E]) Closed() bool {
	if c == nil {
		return true
	}

	_, ok := <-c.c
	if !ok {
		return true
	}
	return false
}

func (c *ChanR[E]) RecvTimeout(recv time.Duration) {
	if c == nil {
		return
	}
	c.recvTimeout = recv
}

func (c *ChanR[E]) Receive() (wrapper.UnWrapper[E], bool) {
	if c == nil {
		return wrap.Nil[E](), false
	}

	if c.recvTimeout <= 0 {
		element, ok := <-c.c
		if ok {
			atomic.AddInt64(&c.count, 1)
		}
		return wrap.Wrap(element), ok
	}

	timer := time.NewTimer(c.recvTimeout)
	select {
	case element, ok := <-c.c:
		if ok {
			atomic.AddInt64(&c.count, 1)
		}
		return wrap.Wrap(element), ok

	case <-timer.C:
		return wrap.Nil[E](), false
	}
}
