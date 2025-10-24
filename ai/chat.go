package ai

import "context"

type Chatter interface {
	Chat(context.Context, string) (string, error)
	ChatStream(context.Context, string) (<-chan string, <-chan error)
}
