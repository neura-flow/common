package exception

import (
	"encoding/json"
)

type Parameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type BizException struct {
	BaseException `json:"-"`
	Status        int          `json:"status"`
	Code          string       `json:"code"`
	Msg           string       `json:"msg,omitempty"`
	TraceId       string       `json:"traceId,omitempty"`
	ErrCode       string       `json:"errorCode,omitempty"`
	Parameters    []*Parameter `json:"parameters,omitempty"`
}

func (e *BizException) Error() string {
	data, _ := json.Marshal(e)
	return string(data)
}

func (e *BizException) HTTPStatus() int {
	return e.Status
}

type Option func(e *BizException)

func WithErrCode(errCode string) Option {
	return func(o *BizException) {
		o.ErrCode = errCode
	}
}

func WithParams(params ...*Parameter) Option {
	return func(o *BizException) {
		o.Parameters = params
	}
}

func WithMsg(msg string) Option {
	return func(e *BizException) {
		e.Msg = msg
	}
}

func WithTraceId(traceId string) Option {
	return func(e *BizException) {
		e.TraceId = traceId
	}
}

// NewBizException new a BizException
func NewBizException(status int, code string, opts ...Option) Exception {
	e := &BizException{
		Status: status,
		Code:   code,
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}
