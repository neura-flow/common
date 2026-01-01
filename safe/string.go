package safe

import (
	"encoding/json"
)

type String struct {
	val string
}

func NewString(val string) *String {
	return &String{val: val}
}

func (s String) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}

func (s *String) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &s.val)
}

func (s String) String() string {
	return ""
}

func (s String) Value() string {
	return s.val
}

func (s *String) SetValue(val string) {
	s.val = val
}
