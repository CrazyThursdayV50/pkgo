package impl

import (
	"context"
	"sync"

	"github.com/CrazyThursdayV50/pkgo/builtin/slice"
	"github.com/CrazyThursdayV50/pkgo/goo"
	"github.com/CrazyThursdayV50/pkgo/worker"
)

type Fan[T any] struct {
	worker *worker.Worker[T]
}

func New[element any](handler func(t element), from ...<-chan element) *Fan[element] {
	var fan Fan[element]
	worker, delivery := worker.New("worker", handler)
	worker.Run(context.TODO())
	worker.WithGraceful(true)
	fan.worker = worker

	var wg sync.WaitGroup
	_, _ = slice.From(from...).Iter(func(_ int, ch <-chan element) (bool, error) {
		wg.Add(1)
		goo.Go(func() {
			defer wg.Done()
			for e := range ch {
				delivery(e)
			}
		})
		return true, nil
	})

	// 所有  from 通道关闭之后，worker 才自动结束
	goo.Go(func() {
		wg.Wait()
		fan.worker.Stop()
	})

	return &fan
}

func (f *Fan[T]) Close() { f.worker.Stop() }
