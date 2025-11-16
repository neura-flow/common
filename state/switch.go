package state

import "sync/atomic"

const (
	On  = State("on")
	Off = State("off")
)

type Switch interface {
	On() (changed bool)
	Off() (changed bool)
	IsOn() bool
}

type switchState struct {
	v        int32
	handlers []SwitchHandler
}

type SwitchHandler func(s Switch, st State)

func NewSwitch(handlers ...SwitchHandler) Switch {
	return &switchState{
		handlers: handlers,
	}
}

func (s *switchState) On() (changed bool) {
	ok := atomic.CompareAndSwapInt32(&s.v, 0, 1)
	if ok {
		s.notify(On)
	}
	return ok
}

func (s *switchState) Off() (changed bool) {
	ok := atomic.CompareAndSwapInt32(&s.v, 1, 0)
	if ok {
		s.notify(Off)
	}
	return ok
}

func (s *switchState) IsOn() bool {
	return atomic.LoadInt32(&s.v) == 1
}

func (s *switchState) notify(st State) {
	for _, h := range s.handlers {
		h(s, st)
	}
}
