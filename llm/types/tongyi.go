package types

import "encoding/json"

type TongYiCompletionRequest struct {
	Base
	Seed               *int           `json:"seed,omitempty"`
	Modalities         []string       `json:"modalities,omitempty"`
	N                  *int           `json:"n,omitempty"`
	ParallelToolCalls  bool           `json:"parallel_tool_calls,omitempty"`
	TranslationOptions map[string]any `json:"translation_options,omitempty"`
	EnableSearch       bool           `json:"enable_search,omitempty"`
}

func (r *TongYiCompletionRequest) Map() (map[string]any, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	kvm := map[string]any{}
	if err = json.Unmarshal(data, &kvm); err != nil {
		return nil, err
	}
	return kvm, nil
}

func NewTongYiCompletionRequest(model string) *TongYiCompletionRequest {
	q := &TongYiCompletionRequest{
		Base: Base{
			Model: model,
		},
	}
	return q
}
