package xerror

type HttpStatus int

type ErrorMessage string

type Error struct {
	status  HttpStatus
	message ErrorMessage
	err     error
}

func New(status HttpStatus, message ErrorMessage, err error) *Error {
	return &Error{status: status, message: message, err: err}
}

func (e *Error) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

func (e *Error) Status() *HttpStatus {
	return &e.status
}

func (e *Error) Message() *ErrorMessage {
	return &e.message
}

func (e *Error) Err() error {
	return e.err
}
