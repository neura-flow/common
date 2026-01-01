package worker

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJob(t *testing.T) {
	const delay = 123 * time.Millisecond
	var j Job
	j = FuncJob(func() {})
	dj := DelayJob{Job: j, Delay: delay}
	assert.NotNil(t, dj)
	assert.Equal(t, delay, dj.Delay)
	cj := CronJob{Job: j, Interval: delay}
	assert.NotNil(t, cj)
	assert.Equal(t, delay, cj.Interval)
	jg := NewJobGroup(JobGroupWithContext(context.Background()), JobGroupWithRecover(), JobGroupWithCancelOnError())
	count := int32(0)
	for i := 0; i < 5; i++ {
		j = jg.NewJob(func() error {
			time.Sleep(100 * time.Millisecond)
			atomic.AddInt32(&count, 1)
			return nil
		}, i)
		DefaultWorker().Push(j)
	}
	err := jg.Wait()
	assert.NoError(t, err)
	assert.Equal(t, atomic.LoadInt32(&count), int32(5))
}
