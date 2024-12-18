package impl

import (
	"context"
	"fmt"

	"github.com/CrazyThursdayV50/pkgo/builtin/slice"
	"github.com/CrazyThursdayV50/pkgo/leader"
	"github.com/CrazyThursdayV50/pkgo/worker"
)

type Fan[T any] struct {
	leader *leader.Leader[T]
	close  func()
}

func New[T any](count, buffer int, handler func(t T)) *Fan[T] {
	var fan Fan[T]
	ctx, cancel := context.WithCancel(context.Background())
	fan.leader = leader.New[T](ctx, 0, 0)
	slice.Make[struct{}](count, count).Iter(func(i int, _ struct{}) (bool, error) {
		worker, _ := worker.New(fmt.Sprintf("Worker_%d", i), handler)
		worker.WithGraceful(true)
		fan.leader.AddWorker(worker)
		return true, nil
	})

	fan.close = func() {
		cancel()
		fan.leader.Stop()
	}

	return &fan
}

func (f *Fan[T]) Close()   { f.close() }
func (f *Fan[T]) Send(t T) { f.leader.Do(t) }
