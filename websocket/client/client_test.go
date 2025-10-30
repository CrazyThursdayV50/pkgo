package client

import (
	"context"
	"testing"
	"time"

	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/CrazyThursdayV50/pkgo/log/sugar"
	"github.com/gorilla/websocket"
)

func TestClient(t *testing.T) {
	cfg := sugar.DefaultConfig()
	cfg.Level = "debug"
	var logger = sugar.New(cfg)

	ctx := context.TODO()

	wsclient := New(
		WithReconnectOnStartup(true),
		WithURL("ws://localhost:18080"),
		WithLogger(logger),

		WithMessageHandler(func(ctx context.Context, l log.Logger, typ int, data []byte, f func(error)) (int, []byte) {
			l.Infof("client receive: %s", data)
			return websocket.BinaryMessage, nil
		}),

		WithSendOnConnect(func() (int, []byte) {
			return TextMessage, []byte("hello!")
		}),

		WithPongHandler(0, func(msg string) error {
			logger.Infof("receive pong: %s", msg)
			return nil
		}),

		WithPingLoop(func(done <-chan struct{}, conn *websocket.Conn) {
			if conn == nil {
				return
			}

			for {
				select {
				case <-done:
					return

				default:
					err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second*30))
					if err != nil {
						return
					}
					time.Sleep(time.Second * 3)
				}
			}
		}),
	)

	err := wsclient.Run(ctx)
	if err != nil {
		logger.Errorf("run ws client failed: %v", err)
	}

	time.Sleep(time.Second * 3600)
}
