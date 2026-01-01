package event

import (
	"testing"

	"github.com/neura-flow/common/metadata"
	"github.com/stretchr/testify/assert"
)

func TestSinkSource(t *testing.T) {
	ss := NewSinkSource(10)
	assert.NotNil(t, ss)
	e := NewEvent(metadata.FromMap(map[string]interface{}{metadata.KeyTopic: "test"}), "test")
	ss.Push(e)
	e1 := ss.Poll()
	assert.Equal(t, e1, e)
	assert.NoError(t, ss.Close())
}
