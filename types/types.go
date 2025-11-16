package types

import (
	"bytes"
	"errors"
	"time"
)

type Timeout struct {
	Dail  int `json:"dail,omitempty"`
	Read  int `json:"read,omitempty"`
	Write int `json:"write,omitempty"`
}

type Pool struct {
	MaxIdle  int `json:"maxIdle,omitempty"`
	MinIdle  int `json:"minIdle,omitempty"`
	MaxOpen  int `json:"maxOpen,omitempty"`
	LifeTime int `json:"lifeTime,omitempty"`
}

type Selector struct {
	Tags map[string]string `json:"tags,omitempty"`
}

type Duration string

func (d *Duration) Val() *time.Duration {
	v, _ := time.ParseDuration(string(*d))
	return &v
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	*d = Duration(bytes.Trim(data, "\""))
	return nil
}

func NewDuration(s string) *Duration {
	d := Duration(s)
	return &d
}

type BatchPolicy struct {
	Enabled    bool        `json:"enabled,omitempty"`
	Count      int         `json:"count,omitempty"`
	ByteSize   int         `json:"byte_size,omitempty"`
	Period     string      `json:"period,omitempty"`
	Check      string      `json:"check,omitempty"`
	Processors interface{} `json:"processors,omitempty"`
}

type Any map[string]interface{}
type Form map[string]string

var (
	NotFound = errors.New("NOT_FOUND")
)
