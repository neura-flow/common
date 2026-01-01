package event

import (
	"testing"
	"time"

	"github.com/neura-flow/common/exception"
	"github.com/neura-flow/common/metadata"
	"github.com/stretchr/testify/assert"
)

func TestInterceptor(t *testing.T) {
	var count int
	testError := StringError("test")
	i := Interceptor(func(e Event, err Error) (Event, Error) {
		if err.Error() == "test" {
			panic(err)
		}
		return e, err
	}).WithRecover(func(pe exception.PanicException) bool {
		err, _ := pe.Unwrap().(Error)
		if err != nil && err.Error() == "test" {
			return true
		}
		return false
	}).WithErrorFilter(func(err Error) Error {
		if err.Code() == "panic" {
			return err
		}
		return nil
	}).WithReporter(ReportFunc(func(e Event, latency time.Duration, err Error) {
		count++
	}))
	testEvt := NewEvent(metadata.FromKVList(metadata.NewKV("test", "test")), "test")
	e, err := i(testEvt, testError)
	assert.Equal(t, 1, count)
	assert.Equal(t, "panic", string(err.Code()))
	assert.Equal(t, testEvt, e)
	testError = "test1"
	e, err = i(testEvt, testError)
	assert.Equal(t, 2, count)
	assert.Equal(t, nil, err)
	testError = "test"
	h1 := HandlerInterceptor(HandleFunc(func(e Event) Error {
		panic(testError)
		return nil
	})).WithRecover(func(pe exception.PanicException) bool {
		err, _ := pe.Unwrap().(Error)
		if err != nil && err.Error() == "test" {
			return true
		}
		return false
	})
	h2 := ErrorInterceptor(func(err Error) Error {
		if err.Code() == "panic" {
			return err
		}
		return nil
	})
	i = ChainInterceptor(h1, h2)
	e, err = i(testEvt, testError)
	assert.Equal(t, "panic", string(err.Code()))
}
