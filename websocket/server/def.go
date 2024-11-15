package server

import (
	"github.com/CrazyThursdayV50/pkgo/log"
	defaultlogger "github.com/CrazyThursdayV50/pkgo/log/default"
	"github.com/CrazyThursdayV50/pkgo/trace"
)

type Server struct {
	tracer          trace.Tracer
	logger          log.Logger
	readBufferSize  int
	writeBufferSize int
	done            chan struct{}
	handler         func(messageType int, data []byte, err error) (int, []byte, error)
}

func New(opts ...Option) *Server {
	var s Server
	s.done = make(chan struct{})
	s.readBufferSize = defaultReadBufferSize
	s.writeBufferSize = defaultWriteBufferSize
	s.logger = defaultlogger.New(defaultlogger.DefaultConfig())
	s.logger.Init()

	for _, opt := range opts {
		opt(&s)
	}

	return &s
}
