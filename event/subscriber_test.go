package event

import (
	"testing"

	"github.com/neura-flow/common/metadata"
	"github.com/stretchr/testify/assert"
)

func TestSubscriber(t *testing.T) {
	const eventTopic = "test"
	d := NewSubscriber()
	h := func(e Event) Error {
		topic := e.Metadata().Value(metadata.KeyTopic)
		assert.Equal(t, eventTopic, topic)
		assert.Equal(t, "hello", e.Data().(string))
		return nil
	}
	d.Subscribe(HandleFunc(h))
	md := metadata.FromMap(map[string]interface{}{metadata.KeyTopic: eventTopic})
	e := NewEvent(md, "hello")
	assert.NoError(t, d.OnEvent(e))
}
