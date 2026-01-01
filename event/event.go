package event

import "github.com/neura-flow/common/metadata"

type Event interface {
	Metadata() metadata.Metadata
	Data() interface{}
}

type Handler interface {
	OnEvent(e Event) Error
}

func Pipeline(stopOnError bool, handlers ...Handler) Handler {
	return &pipeline{
		stopOnError: stopOnError,
		handlers:    handlers,
	}
}

type pipeline struct {
	stopOnError bool
	handlers    []Handler
}

func (p *pipeline) OnEvent(e Event) Error {
	var err Error
	for _, h := range p.handlers {
		err = h.OnEvent(e)
		if err != nil && p.stopOnError {
			return err
		}
	}
	return err
}

type HandleFunc func(e Event) Error

func (f HandleFunc) OnEvent(e Event) Error {
	return f(e)
}

type HashFunc func(e Event) int64

type event struct {
	md   metadata.Metadata
	data interface{}
}

func NewEvent(md metadata.Metadata, data interface{}) Event {
	return event{
		md:   md,
		data: data,
	}
}

func (e event) Metadata() metadata.Metadata {
	return e.md
}

func (e event) Data() interface{} {
	return e.data
}
