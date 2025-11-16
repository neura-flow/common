package state

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFSM(t *testing.T) {
	var ch = make(chan State, 1)
	h := func(m FSM, ctx Context) {
		ch <- ctx.To
		t.Log(ctx.From, ctx.To)
		m.End(ctx.Ctx, nil)
	}
	stateMap := DefaultStateMap()
	m, err := NewFSM(stateMap, HandlerFunc(h))
	assert.NoError(t, err)
	m.Next(context.Background(), Running, nil)
	st := <-ch
	assert.Equal(t, st, Running)
	st = <-ch
	assert.Equal(t, st, stateMap.End)
}

func ExampleNewFSM() {
	var ch = make(chan struct{})
	ctx := context.Background()
	stateMap := DefaultStateMap()
	h := func(m FSM, ctx Context) {
		switch ctx.To {
		case Running:
			fmt.Printf("%s->", ctx.To)
			m.Next(context.WithValue(ctx.Ctx, "data", "hello"), Stopped, nil)
		case Stopped:
			fmt.Printf("%s->", ctx.To)
			m.End(context.WithValue(ctx.Ctx, "data", "end"), nil)
		case End:
			fmt.Printf("%s", ctx.To)
			close(ch)
		}
	}
	m, err := NewFSM(stateMap, HandlerFunc(h))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s->", Begin)
	m.Next(context.WithValue(ctx, "data", "hello"), Running, nil)
	<-ch

	// Output: begin->running->stopped->end
}
