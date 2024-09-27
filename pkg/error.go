package pkg

// Error struct defines a custom error type with an error, status code, and message.
type Error struct {
	err        error
	statusCode int
	message    string
}

// NewError creates a new instance of the Error struct.
func NewError(err error, message string, statusCode int) *Error {
	return &Error{
		err:        err,
		message:    message,
		statusCode: statusCode,
	}
}

// Error implements the error interface by returning the original error message.
func (rc *Error) Error() string {
	return rc.err.Error()
}

// Message returns the custom error message.
func (rc *Error) Message() string {
	return rc.message
}

// StatusCode returns the associated HTTP status code.
func (rc *Error) StatusCode() int {
	return rc.statusCode
}
