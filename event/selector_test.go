package event

import (
	"context"
	"testing"
	"time"

	"github.com/neura-flow/common/metadata"
	"github.com/stretchr/testify/assert"
)

func TestSelector(t *testing.T) {
	ss := NewSinkSource(1)
	ss1 := NewSinkSource(1)
	e := NewEvent(metadata.FromMap(map[string]interface{}{metadata.KeyTopic: "test"}), "test")
	e1 := NewEvent(metadata.FromMap(map[string]interface{}{metadata.KeyTopic: "test1"}), "test")
	ss.Push(e)
	ss1.Push(e1)
	evts, closed, ok := Select([]Source{ss, ss1}, SelectWithMax(2), SelectWithDefaultCase())
	assert.True(t, ok)
	assert.Equal(t, 0, len(closed))
	assert.Equal(t, len(evts), 2)
	m := eventListToMap(evts)
	assert.Equal(t, m["test"], e)
	assert.Equal(t, m["test1"], e1)
	evts, closed, ok = Select([]Source{ss, ss1}, SelectWithTimer(time.NewTimer(time.Millisecond)))
	assert.Equal(t, len(evts), 0)
	assert.False(t, ok)
	ctx, cancel := context.WithCancel(context.Background())
	go cancel()
	evts, closed, ok = Select([]Source{ss}, SelectWithContext(ctx))
	assert.False(t, ok)
	assert.Equal(t, len(evts), 0)
	assert.Equal(t, 0, len(closed))
	ss.Close()
	ss1.Close()
	evts, closed, ok = Select([]Source{ss, ss1}, SelectWithWaitAll())
	assert.False(t, ok)
	assert.Equal(t, len(evts), 0)
	assert.Equal(t, 2, len(closed))
}

func eventListToMap(evts []Event) map[interface{}]Event {
	m := make(map[interface{}]Event)
	for _, e := range evts {
		m[e.Metadata().Value(metadata.KeyTopic)] = e
	}
	return m
}
