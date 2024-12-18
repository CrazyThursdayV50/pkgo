package builtin

import (
	"time"

	"github.com/CrazyThursdayV50/gotils/pkg/wrapper"
)

type (
	chanUnwrapper[C <-chan E | chan<- E | chan E, E any] interface {
		UnWrapper[C]
	}

	baseCommonChanAPI[E any] interface {
		Len() int
		IsEmpty() bool
		Closed() bool
	}

	baseReadWriteChanAPI[E any] interface {
		Renew(buffer int)
		RenewForce(buffer int)
	}

	baseWriteChanAPI[E any] interface {
		Close()
		Send(element E)
		SendTimeout(send time.Duration)
	}

	baseReadChanAPI[E any] interface {
		Receive() (wrapper.UnWrapper[E], bool)
		RecvTimeout(recv time.Duration)
		Iter[int, E]
	}

	ChanAPIR[E any] interface {
		chanUnwrapper[<-chan E, E]
		baseCommonChanAPI[E]
		baseReadChanAPI[E]
	}

	ChanAPIW[E any] interface {
		chanUnwrapper[chan<- E, E]
		baseCommonChanAPI[E]
		baseWriteChanAPI[E]
	}

	ChanAPI[E any] interface {
		chanUnwrapper[chan E, E]
		baseCommonChanAPI[E]
		baseWriteChanAPI[E]
		baseReadChanAPI[E]
		baseReadWriteChanAPI[E]
	}
)

// var (
// 	_ ChanAPI[any]  = (*models.ChanRW[any])(nil)
// 	_ ChanAPIR[any] = (*models.ChanR[any])(nil)
// 	_ ChanAPIR[any] = (*models.ChanRW[any])(nil)
// 	_ ChanAPIW[any] = (*models.ChanW[any])(nil)
// 	_ ChanAPIW[any] = (*models.ChanRW[any])(nil)
// )
