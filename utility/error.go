package utility

type ResError struct {
	StatusCode *int
	Message    *string
	Err        error
	StackTrace *string
}

func (slf *ResError) Error() string {
	if slf.Message != nil {
		return *slf.Message
	}

	return slf.Err.Error()
}

func NewError(err error, statusCode int, errOpts ...*string) *ResError {
	resErr := &ResError{
		Err:        err,
		StatusCode: &statusCode,
	}

	if len(errOpts) > 0 && errOpts[0] != nil {
		resErr.Message = errOpts[0]
	}

	if len(errOpts) > 1 && errOpts[1] != nil {
		resErr.StackTrace = errOpts[1]
	}

	return resErr
}
