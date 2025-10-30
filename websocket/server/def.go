package server

import (
	"context"

	"github.com/CrazyThursdayV50/pkgo/builtin"
	gmap "github.com/CrazyThursdayV50/pkgo/builtin/map"
	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/CrazyThursdayV50/pkgo/log/sugar"
	"github.com/CrazyThursdayV50/pkgo/trace"
	"github.com/CrazyThursdayV50/pkgo/websocket/compressor"
	"github.com/gorilla/websocket"
)

type conn struct {
	id     int64
	conn   *websocket.Conn
	cancel func()
}

type Server struct {
	tracer          trace.Tracer
	logger          log.Logger
	readBufferSize  int
	writeBufferSize int
	handler         func(ctx context.Context, messageType int, data []byte, err error) (int, []byte, error)
	c               compressor.Compressor
	connID          int64
	conns           builtin.MapAPI[int64, *conn]
}

func New(opts ...Option) *Server {
	var s Server
	s.readBufferSize = defaultReadBufferSize
	s.writeBufferSize = defaultWriteBufferSize
	logger := sugar.New(sugar.DefaultConfig())
	s.logger = logger
	s.conns = gmap.Make[int64, *conn](0)

	for _, opt := range opts {
		opt(&s)
	}

	return &s
}
