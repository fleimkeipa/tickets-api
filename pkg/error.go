package pkg

type Error struct {
	err        error
	statusCode int
	message    string
}

func NewError(err error, message string, statusCode int) *Error {
	return &Error{
		err:        err,
		message:    message,
		statusCode: statusCode,
	}
}

func (rc *Error) Error() string {
	return rc.err.Error()
}

func (rc *Error) Message() string {
	return rc.message
}

func (rc *Error) StatusCode() int {
	return rc.statusCode
}
