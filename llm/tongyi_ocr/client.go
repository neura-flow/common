package tongyi_ocr

import (
	"context"
	"errors"

	"github.com/neura-flow/common/llm/common"
	"github.com/neura-flow/common/llm/types"
)

const (
	serverUrl = "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"
)

type Config struct {
	Secret    string `json:"secret"`
	ServerUrl string `json:"serverUrl"`
}

type Client struct {
	cfg *Config
}

func New(cfg *Config) *Client {
	if cfg.ServerUrl == "" {
		cfg.ServerUrl = serverUrl
	}
	return &Client{
		cfg: cfg,
	}
}

func (c *Client) Call(ctx context.Context, completion types.CompletionRequest) (*types.CompletionResponse, error) {
	completion0, ok := completion.(*types.TongYiCompletionRequest)
	if !ok {
		return nil, errors.New("completion is illegal")
	}
	return common.RunRequest(ctx, c.cfg.ServerUrl, c.cfg.Secret, completion0)
}

func (c *Client) newImageRequest(imageUrl string) *Request {
	return &Request{
		Model: "qwen-vl-ocr",
		Messages: []*RequestMessage{
			{
				Role: "user",
				Content: []*RequestMessageContent{
					{
						Type:      "image_url",
						ImageUrl:  imageUrl,
						MinPixels: 3136,
						MaxPixels: 1003520,
					},
					{
						Type: "text",
						Text: "Read all the text in the image.",
					},
				},
			},
		},
	}
}

type Request struct {
	Model    string            `json:"model"`
	Messages []*RequestMessage `json:"messages"`
}

type RequestMessage struct {
	Role    string                   `json:"role"`
	Content []*RequestMessageContent `json:"content"`
}

type RequestMessageContent struct {
	Type      string `json:"type,omitempty"`
	ImageUrl  string `json:"image_url,omitempty"`
	MinPixels int    `json:"min_pixels,omitempty"`
	MaxPixels int    `json:"max_pixels,omitempty"`
	Text      string `json:"text,omitempty"`
}
