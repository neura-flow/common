package deepseek

import (
	"context"
	"errors"

	"github.com/neura-flow/common/llm/common"
	"github.com/neura-flow/common/llm/types"
)

const (
	ServerUrl = "https://api.deepseek.com/chat/completions"
)

type Config struct {
	ServerUrl string `json:"serverUrl"`
	Secret    string `json:"secret"`
}

type Client struct {
	cfg *Config
}

func New(cfg *Config) *Client {
	if cfg.ServerUrl == "" {
		cfg.ServerUrl = ServerUrl
	}
	return &Client{
		cfg: cfg,
	}
}

func (c *Client) Call(ctx context.Context, completion types.CompletionRequest) (*types.CompletionResponse, error) {
	completion0, ok := completion.(*types.DeepSeekCompletionRequest)
	if !ok {
		return nil, errors.New("completion is illegal")
	}
	return common.RunRequest(ctx, c.cfg.ServerUrl, c.cfg.Secret, completion0)
}
