package event

import (
	"sync"
	"testing"
	"time"

	"github.com/neura-flow/common/metadata"
	"github.com/stretchr/testify/assert"
)

func TestConsumer(t *testing.T) {
	const eventTopic = "test"
	ss := NewSinkSource(10)
	c := NewConsumer(ss)
	ch := make(chan struct{})
	h := func(e Event) Error {
		topic := e.Metadata().Value(metadata.KeyTopic)
		assert.Equal(t, eventTopic, topic)
		assert.Equal(t, "hello", e.Data().(string))
		close(ch)
		return nil
	}
	c.Subscribe(HandleFunc(h))
	chStop := make(chan struct{})
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		wg.Done()
		err := c.Start()
		assert.NoError(t, err)
		close(chStop)
	}()
	wg.Wait()
	timer := time.NewTimer(time.Second)
	md := metadata.FromMap(map[string]interface{}{metadata.KeyTopic: eventTopic})
	e := NewEvent(md, "hello")
	assert.NoError(t, ss.Push(e))
	timer.Reset(time.Second)
	select {
	case <-ch:
	case <-chStop:
		t.Error("stopped")
	case <-timer.C:
		t.Error("timeout")
	}
	assert.NoError(t, ss.Close())
	timer.Reset(time.Second)
	select {
	case <-chStop:
	case <-timer.C:
		t.Error("timeout")
	}
}
