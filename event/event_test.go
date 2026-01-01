package event

import (
	"errors"
	"testing"

	"github.com/neura-flow/common/metadata"
	"github.com/stretchr/testify/assert"
)

func TestEvent(t *testing.T) {
	md := metadata.FromMap(map[string]interface{}{metadata.KeyTopic: "test"})
	e := NewEvent(md, "hello")
	assert.NotNil(t, e)
	assert.Equal(t, "test", e.Metadata().Value(metadata.KeyTopic))
	assert.Equal(t, "hello", e.Data())

	var h Handler = HandleFunc(func(e Event) Error {
		if e.Data() == "foo" {
			return nil
		}
		return NewError(errors.New("error"), "error")
	})
	foo := NewEvent(nil, "foo")
	bar := NewEvent(nil, "bar")
	assert.NoError(t, h.OnEvent(foo))
	assert.Error(t, h.OnEvent(bar))
	count := 0
	var h1 Handler = HandleFunc(func(e Event) Error {
		count++
		return nil
	})
	h2 := Pipeline(false, h, h1)
	assert.NoError(t, h2.OnEvent(foo))
	assert.Equal(t, 1, count)
	assert.NoError(t, h2.OnEvent(bar))
	assert.Equal(t, 2, count)
	h2 = Pipeline(true, h, h1)
	assert.NoError(t, h2.OnEvent(foo))
	assert.Equal(t, 3, count)
	assert.Error(t, h2.OnEvent(bar))
	assert.Equal(t, 3, count)
}
