package server

import (
	"github.com/CrazyThursdayV50/gotils/pkg/builtin/api"
	gmap "github.com/CrazyThursdayV50/gotils/pkg/builtin/api/map"
	"github.com/CrazyThursdayV50/pkgo/log"
	defaultlogger "github.com/CrazyThursdayV50/pkgo/log/default"
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
	done            chan struct{}
	handler         func(messageType int, data []byte, err error) (int, []byte, error)
	c               compressor.Compressor
	connID          int64
	// conns           cmap.ConcurrentMap[int64, *conn]
	conns api.MapAPI[int64, *conn]
}

func New(opts ...Option) *Server {
	var s Server
	s.done = make(chan struct{})
	s.readBufferSize = defaultReadBufferSize
	s.writeBufferSize = defaultWriteBufferSize
	s.logger = defaultlogger.New(defaultlogger.DefaultConfig())
	s.logger.Init()
	s.conns = gmap.Make[int64, *conn](0)

	for _, opt := range opts {
		opt(&s)
	}

	return &s
}
