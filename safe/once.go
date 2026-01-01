package safe

import (
	"sync"
	"sync/atomic"
)

type Once struct {
	F     func() interface{}
	once  sync.Once
	value atomic.Value
}

func (v *Once) Value() interface{} {
	v.once.Do(func() {
		v.value.Store(v.F())
	})
	return v.value.Load()
}
