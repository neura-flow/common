package worker

import (
	"testing"
	"time"

	"github.com/neura-flow/common/event"
	"github.com/stretchr/testify/assert"
)

func TestTickWorker(t *testing.T) {
	w := NewCronWorker()
	ch := make(chan struct{}, 1)
	var j Job
	var err error
	j = FuncJob(func() {
		select {
		case ch <- struct{}{}:
		default:
		}
	})
	err = w.Start()
	assert.NoError(t, err)
	err = w.Push(j)
	assert.NoError(t, err)
	timer := time.NewTimer(5 * DefaultTickInterval)
	select {
	case <-ch:
	case <-timer.C:
		t.Error("timeout")
	}
	err = w.Push(DelayJob{Job: j, Delay: DefaultTickInterval})
	assert.NoError(t, err)
	timer.Reset(5 * DefaultTickInterval)
	select {
	case <-ch:
	case <-timer.C:
		t.Error("timeout")
	}
	err = w.Push(CronJob{Job: j, Interval: 2 * DefaultTickInterval})
	assert.NoError(t, err)
	timer.Reset(10 * DefaultTickInterval)
	count := 0
	expCount := 4
	running := true
	for running {
		select {
		case <-ch:
			count++
			t.Log("tick")
			if count == expCount {
				running = false
			}
		case <-timer.C:
			running = false
		}
	}
	assert.Equal(t, expCount, count)
	err = w.Stop()
	assert.NoError(t, err)
}

func TestWorker(t *testing.T) {
	w := NewWorker(event.NewSinkSource(10))
	ch := make(chan struct{}, 1)
	var j Job
	j = FuncJob(func() {
		select {
		case ch <- struct{}{}:
		default:
		}
	})
	err := w.Push(j)
	assert.NoError(t, err)
	err = w.Start()
	assert.NoError(t, err)
	timer := time.NewTimer(5 * DefaultTickInterval)
	select {
	case <-ch:
	case <-timer.C:
		t.Error("timeout")
	}
}
