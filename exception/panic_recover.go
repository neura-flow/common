package exception

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/neura-flow/common/debug"
)

type PanicException interface {
	Exception
	Unwrap() interface{}
	Stack() string
}

func IsPanicException(e error) bool {
	_, ok := e.(PanicException)
	return ok
}

type panicException struct {
	BaseException
	e     interface{}
	stack []byte
}

func NewPanicException(e interface{}, stack []byte) PanicException {
	return &panicException{
		e:     e,
		stack: stack,
	}
}

func (pe *panicException) Stack() string {
	return string(pe.stack)
}

func (pe *panicException) Error() string {
	e := fmt.Sprintf("%v", pe.e)
	prefix := "panic: "
	if strings.HasPrefix(e, "runtime") || strings.HasPrefix(e, "panic") {
		prefix = ""
	}
	return fmt.Sprintf("%s%s\n\n%s", prefix, e, pe.Stack())
}

func (pe *panicException) HTTPStatus() int {
	return http.StatusExpectationFailed
}

func (pe *panicException) Unwrap() interface{} {
	return pe.e
}

type PanicHandler func(pe PanicException) bool

func Recover(f PanicHandler) {
	if e := recover(); e != nil {
		var pe PanicException
		switch v := e.(type) {
		case PanicException:
			pe = v
		default:
			// 跳过 runtime 中 defer 和 recover 的调用栈
			s := debug.GetStack(2, true)
			pe = NewPanicException(e, s)
		}
		if !f(pe) {
			panic(pe)
		}
	}
}
