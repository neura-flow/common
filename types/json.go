package types

import (
	"bytes"
	"strconv"
	"time"
)

type Int int

func (i *Int) UnmarshalJSON(data []byte) error {
	v, err := strconv.ParseInt(string(bytes.Trim(data, "\"")), 10, 64)
	if err != nil {
		return err
	}
	*i = Int(v)
	return nil
}

func (i *Int) Val() int {
	return int(*i)
}

func NewInt(i int) *Int {
	return (*Int)(&i)
}

type Float float64

func (f *Float) UnmarshalJSON(data []byte) error {
	v, err := strconv.ParseFloat(string(bytes.Trim(data, "\"")), 64)
	if err != nil {
		return err
	}
	*f = Float(v)
	return nil
}

func (f *Float) Val() float64 {
	return float64(*f)
}

type Bool bool

func (b *Bool) UnmarshalJSON(data []byte) error {
	v, err := strconv.ParseBool(string(bytes.Trim(data, "\"")))
	if err != nil {
		return err
	}
	*b = Bool(v)
	return nil
}

func (b *Bool) Val() bool {
	return bool(*b)
}

func NewBool(b bool) *Bool {
	return (*Bool)(&b)
}

// TimeMS 毫秒时间
type TimeMS int64

func (t *TimeMS) Duration() time.Duration {
	return time.Duration(*t) * time.Millisecond
}

func (t *TimeMS) UnmarshalJSON(data []byte) error {
	v, err := strconv.ParseInt(string(bytes.Trim(data, "\"")), 10, 64)
	if err != nil {
		return err
	}
	*t = TimeMS(v)

	return nil
}

func (t *TimeMS) Val() int64 {
	return int64(*t)
}
