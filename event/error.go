package event

type ErrorCode string

type Error interface {
	error
	Code() ErrorCode
	Unwrap() error
}

type StringError string

func (e StringError) Error() string {
	return string(e)
}

func (e StringError) Code() ErrorCode {
	return ErrorCode(e)
}

func (e StringError) Unwrap() error {
	return e
}

const (
	ErrorClosed       = StringError("closed")
	ErrorUnknownEvent = StringError("unknownEvent")

	ErrorCodeClosed  = ErrorCode("closed")
	ErrorCodePanic   = ErrorCode("panic")
	ErrorCodeUnknown = ErrorCode("unknown")
)

type evtError struct {
	error
	code ErrorCode
}

func NewError(err error, code ErrorCode) Error {
	return &evtError{
		error: err,
		code:  code,
	}
}

func (e *evtError) Code() ErrorCode {
	return e.code
}

func (e *evtError) Unwrap() error {
	return e.error
}
