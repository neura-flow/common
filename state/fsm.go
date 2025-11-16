package state

import (
	"context"
	"errors"
	"sync"
)

const (
	Begin   = State("begin")
	Running = State("running")
	Stopped = State("stopped")
	End     = State("end")
)

type State string

func (s State) String() string {
	return string(s)
}

type Context struct {
	Ctx  context.Context
	From State
	To   State
	Err  error
}

type Handler interface {
	OnStateChange(m FSM, ctx Context)
}

type HandlerFunc func(m FSM, ctx Context)

func (f HandlerFunc) OnStateChange(m FSM, ctx Context) {
	f(m, ctx)
}

type Map struct {
	Begin State
	Maps  map[State][]State
	End   State
}

func (m *Map) Check() error {
	if m.Begin == "" || m.End == "" {
		return errors.New("begin state or end state is invalid")
	}
	for k, list := range m.Maps {
		if k == "" {
			return errors.New("invalid state")
		}
		if k == m.End {
			return errors.New("end state can not be set in state maps")
		}
		for _, v := range list {
			if v == "" {
				return errors.New("invalid state in state map of '" + k.String() + "'")
			}
		}
	}
	return nil
}

func DefaultStateMap() Map {
	return Map{
		Begin: Begin,
		Maps: map[State][]State{
			Begin:   {Running},
			Stopped: {Running, End},
			Running: {Stopped, End},
		},
		End: End,
	}
}

type FSM interface {
	Is(st State) bool
	Next(ctx context.Context, next State, err error) bool
	End(ctx context.Context, err error)
}

type fsm struct {
	sync.RWMutex
	st   State
	sm   Map
	h    Handler
	next chan Context
}

func NewFSM(sm Map, handler Handler) (FSM, error) {
	if err := sm.Check(); err != nil {
		return nil, err
	}
	m := &fsm{
		sm:   sm,
		st:   sm.Begin,
		next: make(chan Context, 1),
		h:    handler,
	}
	go m.doNext()
	return m, nil
}

func (m *fsm) Is(st State) bool {
	m.RLock()
	defer m.RUnlock()
	return m.st == st
}

func (m *fsm) Next(ctx context.Context, next State, err error) bool {
	old, ok := m.update(next)
	if !ok {
		return false
	}
	m.next <- Context{
		Ctx:  ctx,
		From: old,
		To:   next,
		Err:  err,
	}
	m.Lock()
	defer m.Unlock()
	if next != m.st && next == m.sm.End {
		close(m.next)
	}
	return true
}

func (m *fsm) doNext() {
	for ctx := range m.next {
		m.h.OnStateChange(m, ctx)
	}
}

func (m *fsm) End(ctx context.Context, err error) {
	m.Next(ctx, m.sm.End, err)
}

func (m *fsm) update(next State) (old State, ok bool) {
	m.Lock()
	defer m.Unlock()
	if next == m.st {
		return next, false
	}
	old = m.st
	if next == m.sm.End {
		m.st = next
		return old, true
	}

	for _, v := range m.sm.Maps[m.st] {
		if v == next {
			m.st = next
			return old, true
		}
	}
	return old, false
}
