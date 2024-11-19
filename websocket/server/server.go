package server

import (
	"context"
	"net/http"
	"sync/atomic"

	"github.com/CrazyThursdayV50/gotils/pkg/async/goo"
	"github.com/gorilla/websocket"
	"github.com/opentracing/opentracing-go/log"
)

func (s *Server) Broadcast(ctx context.Context, messageType int, data []byte) {
	s.conns.IterFully(func(k int64, v *conn) error {
		select {
		case <-ctx.Done():
			return nil

		default:
			_ = v.send(ctx, s.tracer, messageType, data)
			return nil
		}
	})
}

func (s *Server) newConn(c *websocket.Conn, cancel func()) *conn {
	id := atomic.AddInt64(&s.connID, 1)
	return &conn{id: id, conn: c, cancel: cancel}
}

func (s *Server) uncompress(ctx context.Context, data []byte) ([]byte, error) {
	span, ctx := s.tracer.NewSpan(ctx)
	defer span.Finish()

	if s.c == nil {
		return data, nil
	}

	data, err := s.c.Uncompress(data)
	if err != nil {
		span.LogFields(log.Event("Uncompress"), log.Error(err))
		return nil, err
	}
	return data, nil
}

func (s *Server) compress(ctx context.Context, data []byte) ([]byte, error) {
	span, ctx := s.tracer.NewSpan(ctx)
	defer span.Finish()

	if s.c == nil {
		return data, nil
	}

	data, err := s.c.Compress(data)
	if err != nil {
		span.LogFields(log.Event("Compress"), log.Error(err))
		return nil, err
	}
	return data, nil
}

func (s *Server) handle(ctx context.Context, messageType int, data []byte, err error, cancel func()) (int, []byte, error) {
	span, ctx := s.tracer.NewSpan(ctx)
	defer span.Finish()
	messageType, data, err = s.handler(messageType, data, err)
	if err != nil {
		span.LogFields(log.Event("handle"), log.Error(err))
		cancel()
	}

	return messageType, data, err
}

func (s *Server) runConn(ctx context.Context, c *conn) error {
	span, ctx := s.tracer.NewSpanWithName(ctx, "loop")
	defer span.Finish()

	messageType, data, err := c.recv(ctx, s.tracer)
	if err != nil {
		return err
	}

	data, err = s.uncompress(ctx, data)
	if err != nil {
		return err
	}

	messageType, data, err = s.handle(ctx, messageType, data, err, c.cancel)
	if err != nil {
		return err
	}

	if data == nil {
		return nil
	}

	data, err = s.compress(ctx, data)
	if err != nil {
		return err
	}

	return c.send(ctx, s.tracer, messageType, data)
}

func (s *Server) Run(ctx context.Context, w http.ResponseWriter, r *http.Request, h http.Header) error {
	var upgrader = websocket.Upgrader{ReadBufferSize: s.readBufferSize, WriteBufferSize: s.writeBufferSize}
	conn, err := upgrader.Upgrade(w, r, h)
	if err != nil {
		s.logger.Errorf("upgrade failed: %v", err)
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	c := s.newConn(conn, cancel)
	// s.conns.AddSoft(c.id, c)
	goo.Go(func() {
		for {
			err := s.runConn(ctx, c)
			if err != nil {
				return
			}
		}
	})

	goo.Go(func() {
		<-ctx.Done()
		conn.Close()
		s.conns.Del(c.id)
		s.logger.Warnf("conn %d closed", c.id)
	})

	return nil
}
