package types

import "context"

type Embedding interface {
	Call(ctx context.Context, request *Request) (*Response, error)
}

type Request struct {
	Model      string             `json:"model,omitempty"`
	Input      *RequestInput      `json:"input,omitempty"`
	Parameters *RequestParameters `json:"parameters,omitempty"`
}

type RequestInput struct {
	Texts []string `json:"texts"`
}

type RequestParameters struct {
	Dimension int `json:"dimension"`
}

type Response struct {
	StatusCode int             `json:"status_code"`
	RequestId  string          `json:"request_id"`
	Code       string          `json:"code"`
	Message    string          `json:"message"`
	Output     *ResponseOutput `json:"output"`
	Usage      *Usage          `json:"usage"`
}

type ResponseOutput struct {
	Embeddings []*ResponseEmbedding `json:"embeddings"`
}

type ResponseEmbedding struct {
	Embedding []float64 `json:"embedding"`
	TextIndex int       `json:"text_index"`
}

type Usage struct {
	TotalTokens int `json:"total_tokens"`
}
