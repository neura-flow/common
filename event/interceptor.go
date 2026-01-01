package event

import (
	"time"

	"github.com/neura-flow/common/exception"
)

// Interceptor 拦截器
type Interceptor func(e Event, err Error) (Event, Error)

// ErrorFilter 错误过滤器
type ErrorFilter func(err Error) Error

// Handler 将拦截器转换成处理器
func (i Interceptor) Handler() Handler {
	return HandleFunc(func(e Event) Error {
		_, err := i(e, nil)
		return err
	})
}

// WithRecover 从panic中恢复
func (i Interceptor) WithRecover(h exception.PanicHandler) Interceptor {
	return func(e Event, err Error) (e1 Event, err1 Error) {
		e1 = e
		defer exception.Recover(func(pe exception.PanicException) bool {
			err1 = NewError(pe, "panic")
			return h(pe)
		})
		return i(e, err)
	}
}

// WithErrorFilter 添加错误过滤器，可以将不需要处理的错误过滤掉
func (i Interceptor) WithErrorFilter(filter ErrorFilter) Interceptor {
	return func(e Event, err Error) (Event, Error) {
		e, err = i(e, err)
		err = filter(err)
		return e, err
	}
}

// WithReporter 添加报告器，可以将监控、日志、告警等功能抽象成报告器加进来
func (i Interceptor) WithReporter(reporters ...Reporter) Interceptor {
	return func(e Event, err Error) (Event, Error) {
		begin := time.Now()
		e, err = i(e, err)
		latency := time.Since(begin)
		for _, r := range reporters {
			r.Report(e, latency, err)
		}
		return e, err
	}
}

// HandlerInterceptor 将一个处理器转换成拦截器
func HandlerInterceptor(h Handler) Interceptor {
	return func(e Event, err Error) (Event, Error) {
		return e, h.OnEvent(e)
	}
}

// ChainInterceptor 按顺序组合拦截器，
func ChainInterceptor(interceptors ...Interceptor) Interceptor {
	return func(e Event, err Error) (Event, Error) {
		for _, i := range interceptors {
			e, err = i(e, err)
		}
		return e, err
	}
}

// ErrorInterceptor 将一个错误过滤器转换成拦截器
func ErrorInterceptor(filter ErrorFilter) Interceptor {
	return func(e Event, err Error) (Event, Error) {
		err = filter(err)
		return e, err
	}
}
