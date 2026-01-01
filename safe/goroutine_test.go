package safe_test

import (
	"testing"
	"time"

	"github.com/neura-flow/common/exception"
	"github.com/neura-flow/common/safe"

	"github.com/stretchr/testify/assert"
)

func TestGo(t *testing.T) {
	ch := make(chan exception.PanicException, 1)
	safe.Go(func() {
		panic("hello")
	}, func(pe exception.PanicException) bool {
		ch <- pe
		return true
	})
	timer := time.NewTimer(300 * time.Millisecond)
	select {
	case pe := <-ch:
		assert.Equal(t, pe.Unwrap(), "hello")
	case <-timer.C:
		t.Error("timeout")
	}
}
