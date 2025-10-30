package deepseek

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/CrazyThursdayV50/pkgo/file"
	"github.com/CrazyThursdayV50/pkgo/log/sugar"
)

func TestChat(t *testing.T) {
	var config Config
	config.Model = "deepseek-chat"
	config.SystemFile = "../.system"
	ctx := context.TODO()
	logger := sugar.New(sugar.DefaultConfig())

	userPrompt, err := file.ReadFileToString("../.user")
	if err != nil {
		t.Fatalf("read user prompt failed: %v", err)
		return
	}

	q, err := file.ReadFileToString("../question.json")
	if err != nil {
		t.Fatalf("read question failed: %v", err)
		return
	}

	q = strings.ReplaceAll(userPrompt, "{{input_data}}", q)
	client, err := New(&config, logger)
	if err != nil {
		t.Fatalf("new client failed: %v", err)
		return
	}

	start := time.Now()
	text, err := client.Chat(ctx, q)
	cost := time.Since(start)
	if err != nil {
		t.Fatalf("chat failed: %v", err)
		return
	}
	fmt.Printf("[%s]receive: %s\n", cost.String(), text)
}
