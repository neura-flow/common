package event

import (
	"github.com/neura-flow/common/metadata"
	"github.com/neura-flow/common/state"
)

type Consumer interface {
	Subscribe(h Handler, md ...metadata.KV)
	Start() error
	Stop() error
}

type consumer struct {
	Source
	Subscriber
	state.Switch
	done chan struct{}
}

func NewConsumer(s Source) Consumer {
	c := &consumer{
		Switch:     state.NewSwitch(),
		Source:     s,
		Subscriber: NewSubscriber(),
		done:       make(chan struct{}),
	}
	return c
}

func (q *consumer) Start() error {
	if !q.On() {
		return nil
	}
	for {
		select {
		case e, ok := <-q.Chan():
			if !ok {
				q.Stop()
				return nil
			}
			q.OnEvent(e)
		case <-q.done:
			return nil
		}
	}
}

func (q *consumer) Stop() error {
	if !q.Off() {
		return nil
	}
	select {
	case <-q.done:
	default:
		close(q.done)
	}
	return nil
}
