package event

import (
	"context"
	"reflect"
	"time"

	"github.com/neura-flow/common/metadata"
)

// Selector 用于从一批 Source 中接收事件，返回事件列表以及已关闭的 Source
type Selector interface {
	// Select 用于从一批 Source 中接收事件，返回事件列表以及已关闭的 Source
	Select(sources []Source, opts ...SelectOption) (results []Event, closed []Source, ok bool)
}

// SelectFunc 用于从一批 Source 中接收事件，返回事件列表以及已关闭的 Source
type SelectFunc func(sources []Source, opts ...SelectOption) (results []Event, closed []Source, ok bool)

// Select 用于从一批 Source 中接收事件，返回事件列表以及已关闭的 Source
func (f SelectFunc) Select(sources []Source, opts ...SelectOption) (results []Event, closed []Source, ok bool) {
	return f(sources, opts...)
}

// SelectOptions Select 选项
type SelectOptions struct {
	//Ctx 用于通过 context 结束 Select
	Ctx context.Context
	//Timer 用于通过 timer 结束 Select
	Timer *time.Timer
	//Max 用于通过最大条数限制结束 Select，默认 10 条
	Max int
	//DefaultCase 用于通过 default 分支结束 Select，值为 true 表示使用 default 分支，默认 false
	DefaultCase bool
	//WaitAll 用于指定等到所有 Source 关闭，默认 false 表示一旦有一个 Source 关闭了，则立即结束
	WaitAll bool
}

func (o *SelectOptions) defaultValue() {
	o.Max = 10
}

// SelectOption Select 选项
type SelectOption func(o *SelectOptions)

// SelectWithContext 用于通过 context 结束 Select
func SelectWithContext(ctx context.Context) SelectOption {
	return func(o *SelectOptions) {
		o.Ctx = ctx
	}
}

// SelectWithTimer 用于通过 timer 结束 Select
func SelectWithTimer(timer *time.Timer) SelectOption {
	return func(o *SelectOptions) {
		o.Timer = timer
	}
}

// SelectWithMax 用于通过最大条数限制结束 Select，最小 1 条
func SelectWithMax(max int) SelectOption {
	return func(o *SelectOptions) {
		if max < 1 {
			max = 1
		}
		o.Max = max
	}
}

// SelectWithDefaultCase 用于通过 default 分支结束 Select
func SelectWithDefaultCase() SelectOption {
	return func(o *SelectOptions) {
		o.DefaultCase = true
	}
}

// SelectWithWaitAll 用于指定等到所有 Source 关闭，默认 false 表示一旦有一个 Source 关闭了，则立即结束
func SelectWithWaitAll() SelectOption {
	return func(o *SelectOptions) {
		o.WaitAll = true
	}
}

// Select 用于从一批 Source 中接收事件，返回事件列表以及已关闭的 Source
func Select(sources []Source, opts ...SelectOption) (results []Event, closed []Source, ok bool) {
	if len(sources) == 0 {
		return nil, nil, false
	}
	ss := append([]Source{}, sources...)
	var opt SelectOptions
	opt.defaultValue()
	for _, o := range opts {
		o(&opt)
	}
	cases := make([]reflect.SelectCase, 0)
	for _, s := range ss {
		c := reflect.SelectCase{
			Chan: reflect.ValueOf(s.Chan()),
			Dir:  reflect.SelectRecv,
		}
		cases = append(cases, c)
	}
	ctlCases := make([]reflect.SelectCase, 0)
	if opt.Timer != nil {
		ctlCases = append(ctlCases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(opt.Timer.C)})
	}
	if opt.Ctx != nil && opt.Ctx.Done() != nil {
		ctlCases = append(ctlCases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(opt.Ctx.Done())})
	}
	if opt.DefaultCase {
		ctlCases = append(ctlCases, reflect.SelectCase{Dir: reflect.SelectDefault})
	}
	cases = append(cases, ctlCases...)
	results = make([]Event, 0)
	closed = make([]Source, 0)
	for len(results) < opt.Max {
		n, v, ok1 := reflect.Select(cases)
		//如果是走到了控制类分支，则退出循环
		if n >= len(cases)-len(ctlCases) {
			return results, closed, !cases[n].Chan.IsValid()
		} else if ok1 {
			results = append(results, v.Interface().(Event))
			continue
		}
		//如果通道已经关闭，则移除该通道
		closed = append(closed, ss[n])
		if !opt.WaitAll || len(closed) == len(sources) {
			return results, closed, false
		} else {
			cases[n].Chan = reflect.Value{}
		}
	}
	return results, closed, true
}

type fanOptions struct {
	ctx     context.Context
	timeout time.Duration
	max     int
}

func FanIn(in []Source, out Sink, opts ...SelectOption) (count int, closed []Source, ok bool) {
	if len(in) == 0 || out == nil {
		return
	}
	var events []Event
	events, closed, ok = Select(in, opts...)
	for _, e := range events {
		if err := out.Push(e); err != nil {
			return count, closed, false
		}
		count++
	}
	return count, closed, true
}

func FanOutByMetadata(in Source, out map[Sink][]metadata.KV, opts ...SelectOption) (count int, closed []Sink, ok bool) {
	if len(out) == 0 || in == nil {
		return
	}
	sb := NewSubscriber()
	handlers := make(map[Sink]Handler)
	for s, md := range out {
		sink := s
		h := HandleFunc(func(e Event) Error {
			err := sink.Push(e)
			if err != nil {
				sb.Unsubscribe(handlers[sink])
				closed = append(closed, sink)
				delete(handlers, sink)
			}
			return err
		})
		handlers[sink] = h
		sb.Subscribe(h, md...)
	}
	var events []Event
	events, _, ok = Select([]Source{in}, opts...)
	for _, e := range events {
		if len(handlers) == 0 {
			return count, closed, false
		}
		if sb.OnEvent(e) == nil {
			count++
		}
	}
	return
}

func FanOutBroadcast(in Source, out []Sink, opts ...SelectOption) (count int, closed []Sink, ok bool) {
	if len(out) == 0 || in == nil {
		return
	}
	sinks := append([]Sink{}, out...)
	var events []Event
	events, _, ok = Select([]Source{in}, opts...)
	for _, e := range events {
		i := 0
		for {
			if len(sinks) == 0 {
				return count, closed, false
			}
			if i == len(sinks) {
				break
			}
			s := sinks[i]
			if s.Push(e) != nil {
				out = append(sinks[:i], sinks[i+1:]...)
				closed = append(closed, s)
			} else {
				count++
				i++
			}
		}
	}
	return
}

func FanOutRoundRobin(in Source, out []Sink, opts ...SelectOption) (count int, closed []Sink, ok bool) {
	if len(out) == 0 || in == nil {
		return
	}
	sinks := append([]Sink{}, out...)
	var events []Event
	events, _, ok = Select([]Source{in}, opts...)
	for _, e := range events {
		i := 0
		for {
			if i == len(sinks) {
				i = 0
			}
			s := sinks[i]
			if s.Push(e) != nil {
				sinks = append(sinks[:i], sinks[i+1:]...)
				closed = append(closed, sinks[i])
				if len(sinks) == 0 {
					return count, closed, false
				}
				continue
			} else {
				count++
				i++
			}
			break
		}
	}
	return count, closed, ok
}
