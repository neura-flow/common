package tongyi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/neura-flow/common/embedding/types"
	"github.com/neura-flow/common/util"
)

var (
	_ = (types.Embedding)(nil)
)

type Client struct {
	cfg *Config
}

func New(cfg *Config) *Client {
	return &Client{
		cfg,
	}
}

func (c *Client) Call(ctx context.Context, request *types.Request) (*types.Response, error) {
	payload := strings.NewReader(util.ToJson(request))
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, c.cfg.ServerUrl, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.cfg.ApiKey))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var embResp *types.Response
	if err = json.Unmarshal(body, &embResp); err != nil {
		return nil, err
	}

	return embResp, nil
}
