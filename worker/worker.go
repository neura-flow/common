package worker

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/neura-flow/common/event"
	"github.com/neura-flow/common/exception"
	"github.com/neura-flow/common/log"
	"github.com/neura-flow/common/safe"

	"github.com/letian0805/go-timewheel"
)

const (
	DefaultPoolSize     = 10
	DefaultTickInterval = 100 * time.Millisecond
)

var defaultCronWorker = &safe.Once{F: func() interface{} {
	w := NewCronWorker()
	go w.Start()
	return w
}}

var defaultWorker = &safe.Once{F: func() interface{} {
	w := NewWorker(event.NewSinkSource(2 * DefaultPoolSize))
	go w.Start()
	return w
}}

func DefaultWorker() Worker {
	return defaultWorker.Value().(Worker)
}

// DefaultTickWorker return default cron worker
// Deprecated: use DefaultCronWorker
func DefaultTickWorker() CronWorker {
	return defaultCronWorker.Value().(CronWorker)
}

// DefaultCronWorker return default cron worker
func DefaultCronWorker() CronWorker {
	return defaultCronWorker.Value().(CronWorker)
}

type Worker interface {
	Push(j Job) error
	Start() error
	Stop() error
}

type Option func(o *options)

type options struct {
	tickInterval time.Duration
	poolSize     int
	logger       log.Logger
	recover      bool
}

func newOptions() options {
	return options{
		poolSize:     DefaultPoolSize,
		tickInterval: DefaultTickInterval,
	}
}

func WithLogger(l log.Logger) Option {
	return func(o *options) {
		o.logger = l
	}
}

func WithPoolSize(poolSize int) Option {
	return func(o *options) {
		if poolSize < 1 {
			poolSize = DefaultPoolSize
		}
		o.poolSize = poolSize
	}
}

func WithTickInterval(interval time.Duration) Option {
	return func(o *options) {
		if interval < DefaultTickInterval {
			interval = DefaultTickInterval
		}
		o.tickInterval = interval
	}
}

func WithRecover() Option {
	return func(o *options) {
		o.recover = true
	}
}

type worker struct {
	event.SinkSource
	sync.RWMutex
	options
	running int32
	stopCh  chan struct{}
}

func NewWorker(ss event.SinkSource, opts ...Option) Worker {
	w := &worker{
		SinkSource: ss,
		options:    newOptions(),
	}
	for _, o := range opts {
		o(&w.options)
	}
	return w
}

func (w *worker) Start() error {
	if !atomic.CompareAndSwapInt32(&w.running, 0, 1) {
		return nil
	}
	wg := &sync.WaitGroup{}
	wg.Add(w.poolSize)
	stopCh := make(chan struct{})
	w.Lock()
	w.stopCh = stopCh
	w.Unlock()
	for i := 0; i < w.poolSize; i++ {
		go func() {
			wg.Done()
			for {
				select {
				case evt, ok := <-w.Chan():
					if !ok {
						return
					}
					if j, ok := evt.(Job); ok {
						w.doJob(j)
					}
				case <-stopCh:
					return
				}
			}
		}()
	}
	wg.Wait()
	return nil
}

func (w *worker) doJob(j Job) {
	defer exception.Recover(func(pe exception.PanicException) bool {
		if w.logger != nil {
			w.logger.Errorf("%v", pe)
		}
		return w.options.recover
	})
	j.Do()
}

func (w *worker) Stop() (err error) {
	if !atomic.CompareAndSwapInt32(&w.running, 1, 0) {
		return nil
	}
	w.RLock()
	defer exception.Recover(func(pe exception.PanicException) bool {
		err = pe
		if w.logger != nil {
			w.logger.Errorf("%v", pe)
		}
		return w.options.recover
	})
	defer w.RUnlock()
	close(w.stopCh)
	return nil
}

func (w *worker) Push(j Job) error {
	return w.SinkSource.Push(j)
}

type CronWorker interface {
	Worker
	Interval() time.Duration
}

type cronWorker struct {
	tw *timewheel.TimeWheelPool
	options
}

func NewCronWorker(opts ...Option) CronWorker {
	w := &cronWorker{
		options: newOptions(),
	}
	for _, o := range opts {
		o(&w.options)
	}
	tw, _ := timewheel.NewTimeWheelPool(w.poolSize, w.tickInterval, 3600)
	w.tw = tw
	return w
}

func (w *cronWorker) Push(j Job) error {
	tw := w.tw.GetRandom()
	if v, ok := j.(CronJob); ok {
		tw.AddCron(v.Interval, j.Do)
	} else if v, ok := j.(DelayJob); ok {
		tw.Add(v.Delay, v.Do)
	} else {
		_ = DefaultWorker().Push(j)
	}
	return nil
}

func (w *cronWorker) Interval() time.Duration {
	return w.tickInterval
}

func (w *cronWorker) Start() error {
	w.tw.Start()
	return nil
}

func (w *cronWorker) Stop() error {
	w.tw.Stop()
	return nil
}
