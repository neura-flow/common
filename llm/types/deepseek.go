package types

import "encoding/json"

type DeepSeekCompletionRequest struct {
	Base
	Logprobs    bool `json:"logprobs,omitempty"`
	TopLogprobs *int `json:"top_logprobs,omitempty"`
}

func (r *DeepSeekCompletionRequest) Map() (map[string]any, error) {
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

func NewDeepSeekCompletionRequest(model string) *DeepSeekCompletionRequest {
	q := &DeepSeekCompletionRequest{
		Base: Base{
			Model: model,
		},
	}
	return q
}
