package types

import (
	"context"
)

type Client interface {
	Call(ctx context.Context, completion CompletionRequest) (*CompletionResponse, error)
}

type CompletionRequest interface {
	Map() (map[string]any, error)
}

type Base struct {
	Model             string           `json:"model"`
	Messages          []any            `json:"messages"`
	Stream            bool             `json:"stream,omitempty"`
	StreamOptions     *StreamOptions   `json:"stream_options,omitempty"`
	Temperature       *float64         `json:"temperature,omitempty"`
	TopP              *float64         `json:"top_p,omitempty"`
	PresencePenalty   *float64         `json:"presence_penalty,omitempty"`
	FrequencyPenalty  *float64         `json:"frequency_penalty,omitempty"`
	ResponseFormat    *ResponseFormat  `json:"response_format,omitempty"`
	MaxTokens         *int             `json:"max_tokens,omitempty"`
	StopWords         []string         `json:"stop_words,omitempty"`
	Tools             []map[string]any `json:"tools,omitempty"`
	RepetitionPenalty *float64         `json:"repetition_penalty,omitempty"`
	ToolChoice        any              `json:"tool_choice,omitempty"`
}

type StreamOptions struct {
	IncludeUsage bool `json:"include_usage"`
}

const (
	ResponseFormat_JSON = "json_object"
	ResponseFormat_Text = "text"
)

type ResponseFormat struct {
	// Must be one of text or json_object.
	Type string `json:"type"`
}

type SystemMessage struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type UserMessage struct {
	Role    string `json:"role,omitempty"`
	Content any    `json:"content,omitempty"`
}

type AssistantMessage struct {
	Role      string           `json:"role,omitempty"`
	Content   string           `json:"content,omitempty"`
	Partial   bool             `json:"partial,omitempty"`
	ToolCalls []map[string]any `json:"tool_calls,omitempty"`
}

type ToolMessage struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type CompletionResponse struct {
	Content *CompletionResponseContent `json:"content,omitempty"`
	Error   *CompletionResponseError   `json:"error,omitempty"`
}

type CompletionResponseError struct {
	RequestId string                         `json:"requestId"`
	Error     *CompletionResponseErrorDetail `json:"error"`
}

type CompletionResponseErrorDetail struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   any    `json:"param"`
	Code    any    `json:"code"`
}

type CompletionResponseContent struct {
	ID                string                      `json:"id,omitempty"`
	Choices           []*CompletionResponseChoice `json:"choices,omitempty"`
	Created           int64                       `json:"created,omitempty"`
	Model             string                      `json:"model,omitempty"`
	SystemFingerprint string                      `json:"system_fingerprint,omitempty"`
	Object            string                      `json:"object,omitempty"`
	Usage             *CompletionResponseUsage    `json:"usage,omitempty"`
}

type CompletionResponseUsage struct {
	CompletionTokens        int                                    `json:"completion_tokens,omitempty"`
	PromptTokens            int                                    `json:"prompt_tokens,omitempty"`
	PromptTokensDetails     *CompletionResponsePromptTokensDetails `json:"prompt_tokens_details,omitempty"`
	PromptCacheHitTokens    int                                    `json:"prompt_cache_hit_tokens,omitempty"`
	PromptCacheMissTokens   int                                    `json:"prompt_cache_miss_tokens,omitempty"`
	TotalTokens             int                                    `json:"total_tokens,omitempty"`
	CompletionTokensDetails *CompletionResponseTokensDetails       `json:"completion_tokens_details,omitempty"`
}

type CompletionResponsePromptTokensDetails struct {
	CachedTokens int `json:"cached_tokens,omitempty"`
}

type CompletionResponseTokensDetails struct {
	ReasoningTokens int `json:"reasoning_tokens,omitempty"`
}

type CompletionResponseChoice struct {
	FinishReason string                     `json:"finish_reason,omitempty"`
	Index        int                        `json:"index,omitempty"`
	Message      *CompletionResponseMessage `json:"message,omitempty"`
	Logprob      *CompletionResponseLogprob `json:"logprob,omitempty"`
}

type CompletionResponseMessage struct {
	Content          string                        `json:"content,omitempty"`
	ReasoningContent string                        `json:"reasoning_content,omitempty"`
	Role             string                        `json:"role,omitempty"`
	ToolCalls        []*CompletionResponseToolCall `json:"tool_calls,omitempty"`
}

type CompletionResponseToolCall struct {
	ID       string                          `json:"id,omitempty"`
	Type     string                          `json:"type,omitempty"`
	Function *CompletionResponseToolCallFunc `json:"function,omitempty"`
}

type CompletionResponseToolCallFunc struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}

type CompletionResponseLogprob struct {
	Content []*CompletionResponseLogprobContent `json:"content,omitempty"`
}

type CompletionResponseLogprobContent struct {
	Token       string                           `json:"token,omitempty"`
	Logprob     float64                          `json:"logprob,omitempty"`
	Bytes       []int64                          `json:"bytes,omitempty"`
	TopLogprobs []*CompletionResponseTopLogprobs `json:"top_logprobs,omitempty"`
}

type CompletionResponseTopLogprobs struct {
	Token   string  `json:"token,omitempty"`
	Logprob float64 `json:"logprob,omitempty"`
	Bytes   []int64 `json:"bytes,omitempty"`
}
