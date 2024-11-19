package server

import (
	"context"

	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/CrazyThursdayV50/pkgo/trace"
	"github.com/CrazyThursdayV50/pkgo/websocket/compressor"
)

type Option func(*Server)

func WithLogger(l log.Logger) Option {
	return func(s *Server) {
		s.logger = l
	}
}

const (
	defaultReadBufferSize  = 1 << 10
	defaultWriteBufferSize = 1 << 10
)

func WithBuffer(read, write int) Option {
	return func(s *Server) {
		if read <= 0 {
			read = defaultReadBufferSize
		}

		if write <= 0 {
			write = defaultWriteBufferSize
		}

		s.readBufferSize, s.writeBufferSize = read, write
	}
}

func WithTracer(tracer trace.Tracer) Option {
	return func(s *Server) {
		s.tracer = tracer
	}
}

func WithHandler(handler func(ctx context.Context, messageType int, data []byte, err error) (int, []byte, error)) Option {
	return func(s *Server) {
		s.handler = handler
	}
}

func WithCompressor(compressor compressor.Compressor) Option {
	return func(s *Server) {
		s.c = compressor
	}
}
