package event

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/neura-flow/common/metadata"
)

type Subscriber interface {
	OnEvent(e Event) error
	Subscribe(h Handler, md ...metadata.KV)
	Unsubscribe(h Handler, md ...metadata.KV)
	Handlers(md ...metadata.KV) []Handler
}

type subscriber struct {
	sync.RWMutex
	handlers map[string]map[interface{}]Handler
}

func NewSubscriber() Subscriber {
	s := &subscriber{
		handlers: make(map[string]map[interface{}]Handler),
	}
	return s
}

func kvToString(kv metadata.KV) string {
	if s, ok := kv.(fmt.Stringer); ok {
		return s.String()
	}
	return fmt.Sprintf("%s:%v", kv.Key(), kv.Value())
}

func (s *subscriber) Subscribe(h Handler, md ...metadata.KV) {
	s.Lock()
	defer s.Unlock()
	if h == nil {
		return
	}
	if len(md) == 0 {
		s.doSubscribe(true, h, "")
	}
	for _, kv := range md {
		s.doSubscribe(true, h, kvToString(kv))
	}
}

func (s *subscriber) doSubscribe(sub bool, h Handler, keys ...string) {
	for _, key := range keys {
		handlers := s.handlers[key]
		if sub {
			if handlers == nil {
				handlers = make(map[interface{}]Handler)
				s.handlers[key] = handlers
			}
			handlers[reflect.ValueOf(h)] = h
		} else {
			if h == nil {
				delete(s.handlers, key)
			} else if handlers != nil {
				delete(handlers, reflect.ValueOf(h))
			}
		}
	}
}

func (s *subscriber) Unsubscribe(h Handler, md ...metadata.KV) {
	s.Lock()
	defer s.Unlock()
	if len(md) == 0 {
		s.doSubscribe(false, h, "")
		if h != nil {
			for _, hs := range s.handlers {
				delete(hs, h)
			}
		}
	}
	for _, kv := range md {
		s.doSubscribe(false, h, kvToString(kv))
	}
}

func (s *subscriber) OnEvent(e Event) error {
	handlers := s.Handlers(e.Metadata().List()...)
	var err error
	for _, h := range handlers {
		if err1 := h.OnEvent(e); err1 != nil {
			err = err1
		}
	}
	return err
}

// Handlers 获取已订阅的 Handler。如果 md 为空，返回订阅全量事件的 Handler；否则按照 md 获取相应的 Handler，含订阅全量事件的 Handler
func (s *subscriber) Handlers(md ...metadata.KV) []Handler {
	var keys = make(map[string]struct{})
	handlers := make(map[interface{}]Handler)
	var res []Handler
	for _, kv := range md {
		if kv == nil {
			continue
		}
		keys[kv.String()] = struct{}{}
		keys[metadata.NewKV(kv.Key(), nil).String()] = struct{}{}
	}
	keys[""] = struct{}{}
	s.RLock()
	defer s.RUnlock()
	for key, _ := range keys {
		for k, h := range s.handlers[key] {
			if _, ok := handlers[k]; !ok {
				handlers[k] = h
				res = append(res, h)
			}
		}
	}
	return res
}
