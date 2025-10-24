package gemini

import (
	"context"
	"testing"

	defaultlogger "github.com/CrazyThursdayV50/pkgo/log/default"
)

func TestChat(t *testing.T) {
	var config Config
	config.Model = "gemini-2.5-flash"
	ctx := context.TODO()
	logger := defaultlogger.New(defaultlogger.DefaultConfig())
	logger.Init()

	client, err := New(ctx, logger, &config)
	if err != nil {
		t.Fatalf("new client failed: %v", err)
		return
	}

	text, err := client.Chat(ctx, "hello, i'm Alex, who are you?")
	if err != nil {
		t.Fatalf("chat failed: %v", err)
		return
	}
	t.Logf("receive: %s", text)
}
