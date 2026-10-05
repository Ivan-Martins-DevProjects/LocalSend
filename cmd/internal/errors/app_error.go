package apperror

type AppError struct {
	Code    string
	Message string
	Err     error
}

func (a *AppError) Error() string {
	return a.Message
}

func InternalServerError(message string, err error) *AppError {
	return &AppError{
		Code:    "INTERNAL_SERVER_ERRRO",
		Message: message,
		Err:     err,
	}
}
