package event

import (
	"sync"
	"sync/atomic"
)

type Sink interface {
	Push(e Event) Error
	Close() error
}

type Source interface {
	Poll() Event
	Chan() <-chan Event
	Close() error
}

func SinkToHandler(s Sink) Handler {
	return HandleFunc(func(e Event) Error {
		return s.Push(e)
	})
}

type SinkSource interface {
	Sink
	Source
}

type Filter func(e Event) Event

type joinSinkSource struct {
	source SinkSource
	sink   SinkSource
	filter Filter
	ch     chan struct{}
	closed int32
}

func Join(source SinkSource, filter Filter, sink SinkSource) SinkSource {
	ss := &joinSinkSource{
		source: source,
		sink:   sink,
		filter: filter,
		ch:     make(chan struct{}),
	}
	go ss.doFilter()
	return ss
}

func (ss *joinSinkSource) doFilter() {
	for {
		select {
		case j, ok := <-ss.source.Chan():
			if !ok {
				return
			}
			if ss.sink.Push(ss.filter(j)) != nil {
				return
			}
		case <-ss.ch:
			atomic.StoreInt32(&ss.closed, 1)
			return
		}
	}
}

func (ss *joinSinkSource) Push(e Event) Error {
	return ss.source.Push(e)
}

func (ss *joinSinkSource) Close() error {
	if !atomic.CompareAndSwapInt32(&ss.closed, 0, 1) {
		return nil
	}
	close(ss.ch)
	return nil
}

func (ss *joinSinkSource) Poll() Event {
	return ss.sink.Poll()
}

func (ss *joinSinkSource) Chan() <-chan Event {
	return ss.sink.Chan()
}

type sinkSource struct {
	closed int32
	wg     sync.WaitGroup
	ch     chan Event
}

func NewSinkSource(size int) SinkSource {
	s := &sinkSource{
		closed: 0,
		wg:     sync.WaitGroup{},
		ch:     make(chan Event, size),
	}
	return s
}

func (s *sinkSource) Push(e Event) Error {
	s.wg.Add(1)
	defer s.wg.Done()
	if atomic.LoadInt32(&s.closed) == 1 {
		return ErrorClosed
	}
	s.ch <- e
	return nil
}

func (s *sinkSource) Poll() Event {
	j := <-s.ch
	return j
}

func (s *sinkSource) Chan() <-chan Event {
	return s.ch
}

func (s *sinkSource) Close() error {
	if !atomic.CompareAndSwapInt32(&s.closed, 0, 1) {
		return nil
	}
	s.wg.Wait()
	close(s.ch)
	return nil
}
