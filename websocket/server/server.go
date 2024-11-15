package server

import (
	"context"
	"net/http"

	"github.com/CrazyThursdayV50/gotils/pkg/async/goo"
	"github.com/gorilla/websocket"
	"github.com/opentracing/opentracing-go/log"
)

func (s *Server) Run(ctx context.Context, w http.ResponseWriter, r *http.Request, h http.Header) error {
	var upgrader = websocket.Upgrader{ReadBufferSize: s.readBufferSize, WriteBufferSize: s.writeBufferSize}
	conn, err := upgrader.Upgrade(w, r, h)
	if err != nil {
		return err
	}

	var run = func(ctx context.Context) bool {
		span, ctx := s.tracer.NewSpan(ctx)
		defer span.Finish()

		messageType, message, err := s.handler(conn.ReadMessage())
		if err != nil {
			span.LogFields(log.Event("ReadMessage"), log.Error(err))
			return false
		}

		if message == nil {
			return true
		}

		err = conn.WriteMessage(messageType, message)
		if err != nil {
			span.LogFields(log.Event("WriteMessage"), log.Error(err))
			return false
		}

		return true
	}

	ctx, cancel := context.WithCancel(ctx)

	goo.Go(func() {
		for {
			ok := run(ctx)
			if !ok {
				cancel()
				return
			}
		}
	})

	goo.Go(func() {
		<-ctx.Done()
		close(s.done)
	})

	goo.Go(func() {
		<-s.done
		conn.Close()
	})

	return nil
}
