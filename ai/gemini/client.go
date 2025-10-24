package gemini

import (
	"context"
	"os"

	"github.com/CrazyThursdayV50/pkgo/ai"
	"github.com/CrazyThursdayV50/pkgo/file"
	"github.com/CrazyThursdayV50/pkgo/log"
	"google.golang.org/genai"
)

var _ ai.Chatter = (*Client)(nil)

type Client struct {
	cfg           *Config
	client        *genai.Client
	systemContent string
}

func New(ctx context.Context, logger log.Logger, cfg *Config) (*Client, error) {
	if cfg.Token != "" {
		os.Setenv("GEMINI_API_KEY", cfg.Token)
	}

	var c Client
	c.cfg = cfg

	if cfg.SystemFile != "" {
		systemContent, err := file.ReadFileToString(cfg.SystemFile)
		if err != nil {
			return nil, err
		}
		c.systemContent = systemContent
	}

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, err
	}

	c.client = client
	return &c, nil
}

func (c *Client) System() string { return c.systemContent }

func (c *Client) Chat(ctx context.Context, q string) (string, error) {
	var content = []*genai.Content{
		genai.NewContentFromText(q, genai.RoleUser),
	}

	if c.systemContent != "" {
		content = append(content, genai.NewContentFromText(c.systemContent, genai.RoleModel))
	}

	result, err := c.client.Models.GenerateContent(
		ctx,
		c.cfg.Model,
		content,
		nil,
	)
	if err != nil {
		return "", err
	}
	return result.Text(), nil
}

func (c *Client) ChatStream(ctx context.Context, q string) (<-chan string, <-chan error) {
	seq := c.client.Models.GenerateContentStream(
		ctx,
		c.cfg.Model,
		[]*genai.Content{
			genai.NewContentFromText(c.systemContent, genai.RoleModel),
			genai.NewContentFromText(q, genai.RoleUser),
		},
		nil,
	)

	var textChan = make(chan string, 100)
	var errChan = make(chan error, 1)
	go func() {
		defer close(errChan)
		defer close(textChan)

		seq(func(resp *genai.GenerateContentResponse, err error) bool {
			if err != nil {
				errChan <- err
				return false
			}

			textChan <- resp.Text()
			return true
		})
	}()

	return textChan, errChan
}
