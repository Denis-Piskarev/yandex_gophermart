package customErrors

type CustomError struct {
	Err        string
	StatusCode int
}

func NewCustomError(err string, statusCode int) CustomError {
	return CustomError{
		Err:        err,
		StatusCode: statusCode,
	}
}

func (e CustomError) Error() string {
	return e.Err
}
