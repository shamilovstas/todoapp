package api

type HttpError struct {
	code int
	error
}

func NewHttpError(code int, e error) *HttpError {
	return &HttpError{code: code, error: e}
}

func (e HttpError) HttpStatus() int {
	return e.code
}

func (e HttpError) Unwrap() error {
	return e.error
}
