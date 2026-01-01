package common

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/neura-flow/common/llm/types"
	"github.com/neura-flow/common/util"
)

func RunRequest(ctx context.Context, serverUrl, secret string, val interface{}) (*types.CompletionResponse, error) {
	payload := strings.NewReader(util.ToJson(val))
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, serverUrl, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", secret))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var content *types.CompletionResponseContent
	if err = json.Unmarshal(body, &content); err != nil {
		return nil, err
	}
	if content.ID != "" {
		return &types.CompletionResponse{
			Content: content,
		}, nil
	}

	var errorResp *types.CompletionResponseError
	if err := json.Unmarshal(body, &errorResp); err != nil {
		return nil, err
	}
	return &types.CompletionResponse{
		Error: errorResp,
	}, nil
}
