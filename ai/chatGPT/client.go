package chatgpt

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/CrazyThursdayV50/pkgo/ai"
	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/sashabaranov/go-openai"
)

var _ ai.Chatter = (*Client)(nil)

type Client struct {
	cfg           *Config
	client        *openai.Client
	systemContent string
}

func New(cfg *Config, logger log.Logger) (*Client, error) {
	if cfg.Token == "" {
		cfg.Token = os.Getenv("OPENAI_API_KEY")
	}

	c := openai.NewClient(cfg.Token)

	client := Client{cfg: cfg, client: c, systemContent: cfg.SystemContent}
	return &client, nil
}

func (c *Client) Chat(ctx context.Context, q string) (string, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: q,
		},
	}

	if c.systemContent != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: c.systemContent,
		})
	}

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       c.cfg.Model,
		Messages:    messages,
		Temperature: c.cfg.Temperature,
	})
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func (c *Client) ChatStream(ctx context.Context, q string) (<-chan string, <-chan error) {
	var textChan = make(chan string, 100)
	var errChan = make(chan error, 1)

	go func() {
		defer close(errChan)
		defer close(textChan)

		messages := []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: q,
			},
		}

		if c.systemContent != "" {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: c.systemContent,
			})
		}

		stream, err := c.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
			Model:       c.cfg.Model,
			Messages:    messages,
			Temperature: c.cfg.Temperature,
		})

		if err != nil {
			errChan <- err
			return
		}

		defer stream.Close()

		for {
			resp, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return
			}

			if err != nil {
				errChan <- err
				return
			}

			textChan <- resp.Choices[0].Delta.Content
		}
	}()

	return textChan, errChan
}

func (c *Client) System() string { return c.systemContent }
