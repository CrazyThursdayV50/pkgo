package leader

import (
	"context"
	"testing"
	"time"

	"github.com/CrazyThursdayV50/pkgo/worker"
)

func TestLeader(t *testing.T) {
	leader := New[int](context.TODO(), time.Millisecond, time.Millisecond)
	handler1 := func(id int) {
		t.Logf("id[1]: %d\n", id)
	}
	handler2 := func(id int) {
		t.Logf("id[2]: %d\n", id)
	}

	w1, _ := worker.New[int]("worker1", handler1)
	w2, _ := worker.New[int]("worker2", handler2)
	leader.AddWorker(w1)
	leader.AddWorker(w2)
	for id := range make([]int, 100) {
		leader.Do(id)
	}

	time.Sleep(time.Second * 10)
	// <-make(chan struct{})
}
