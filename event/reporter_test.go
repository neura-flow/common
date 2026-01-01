package event

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/neura-flow/common/log"
	"github.com/neura-flow/common/metadata"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
)

type testAlarm struct {
	code string
}

func (a *testAlarm) SetLevel(l log.Level) {}

func (a *testAlarm) Alarm(level log.Level, component string, errCode string, err string, traceId string) {
	a.code = errCode
}

func TestReporter(t *testing.T) {
	addr := "127.0.0.1:12345"
	lis, err := net.Listen("tcp", addr)
	assert.NoError(t, err)
	m := http.NewServeMux()
	m.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{
		Handler: m,
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		wg.Done()
		srv.Serve(lis)
	}()
	wg.Wait()
	time.Sleep(200 * time.Millisecond)
	defer srv.Close()
	opts := MetricsReporterOpts{
		Namespace: "http",
		Subsystem: "request",
		LabelKeys: []string{"foo"},
		ConstLabels: map[string]string{
			metadata.KeyTopic: "test",
		},
	}
	mr, err := MetricsReporter(opts)
	assert.NoError(t, err)
	evt := NewEvent(metadata.FromKVList(metadata.NewKV("foo", "bar")), "hello")
	mr.Report(evt, 1000*time.Millisecond, ErrorClosed)
	resp, err := http.Get(fmt.Sprintf("http://%s/metrics", addr))
	assert.NoError(t, err)
	data, _ := io.ReadAll(resp.Body)
	fmt.Println(string(data))
	assert.True(t, bytes.Contains(data, []byte(`http_request_duration_bucket{foo="bar",topic="test",le="1000"} 1`)))
	assert.True(t, bytes.Contains(data, []byte(`http_request_duration_bucket{foo="bar",topic="test",le="500"} 0`)))
	assert.True(t, bytes.Contains(data, []byte(`http_request_duration_bucket{foo="bar",topic="test",le="500"} 0`)))
	lr := LogReporter(log.DefaultLogger(), func(e Event) bool {
		return true
	})
	lr.Report(evt, time.Millisecond, ErrorClosed)
	a := &testAlarm{}
	ar := AlarmReporter(a, func(err Error) log.Level {
		return log.LevelError
	})
	ar.Report(evt, time.Millisecond, ErrorClosed)
	assert.Equal(t, string(ErrorClosed.Code()), a.code)
}
