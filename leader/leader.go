package leader

import (
	"context"
	"time"

	"github.com/CrazyThursdayV50/pkgo/builtin"
	gchan "github.com/CrazyThursdayV50/pkgo/builtin/chan"
	"github.com/CrazyThursdayV50/pkgo/log"
	defaultlogger "github.com/CrazyThursdayV50/pkgo/log/default"
	"github.com/CrazyThursdayV50/pkgo/worker"
)

type Leader[J any] struct {
	ctx          context.Context
	cancel       context.CancelFunc
	deliveryChan builtin.ChanAPI[J]
	logger       log.Logger
}

func (b *Leader[J]) Do(job J) {
	b.deliveryChan.Send(job)
}

func (b *Leader[J]) AddWorker(worker *worker.Worker[J]) {
	worker.WithContext(b.ctx)
	worker.WithTrigger(b.deliveryChan.Unwrap())
	worker.WithLogger(b.logger)
	worker.Run()
}

func New[J any](ctx context.Context, sendTimeout, recvTimeout time.Duration) *Leader[J] {
	var leader Leader[J]
	leader.ctx, leader.cancel = context.WithCancel(ctx)

	leader.deliveryChan = gchan.Make[J](0)
	leader.deliveryChan.SendTimeout(sendTimeout)
	leader.deliveryChan.RecvTimeout(recvTimeout)
	leader.logger = defaultlogger.New(defaultlogger.DefaultConfig())
	leader.logger.Init()
	return &leader
}

func (b *Leader[J]) SetLogger(logger log.Logger) { b.logger = logger }

func (b *Leader[J]) Stop() {
	b.cancel()
	b.deliveryChan.Close()
}
