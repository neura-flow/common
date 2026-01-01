package safe

import "github.com/neura-flow/common/exception"

func Go(f func(), handler exception.PanicHandler) {
	go func() {
		defer exception.Recover(handler)
		f()
	}()
}
