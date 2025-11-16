package exception

import (
	"net/http"
)

const (
	DefaultStatus        = http.StatusExpectationFailed
	InvalidRequestStatus = http.StatusBadRequest
	MethodNotSupport     = http.StatusUnsupportedMediaType
	InnerExceptionStatus = http.StatusInternalServerError
)

type BaseException struct{}

func (e *BaseException) isException() bool {
	return true
}

type Exception interface {
	isException() bool
	HTTPStatus() int
	Error() string
}

func Throw(e Exception) {
	panic(e)
}

type Handler func(e Exception) bool

type CatchFunc func(e interface{}) bool

func CatchException(h Handler) CatchFunc {
	return func(e interface{}) bool {
		if e, ok := e.(Exception); ok {
			return h(e)
		}
		return false
	}
}

// Try implemented try-catch
func Try(f func(), catch ...CatchFunc) {
	h := func(pe PanicException) bool {
		caught := false
		for _, c := range catch {
			if ok := c(pe.Unwrap()); ok {
				caught = ok
				break
			}
		}
		return caught
	}
	defer Recover(h)
	f()
}
