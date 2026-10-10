package apperror

const (
	INTERNAL_ERROR = "INTERNAL_SERVER_ERROR"
	FILE_NOT_FOUND = "FILE_NOT_FOUND"
)

type AppError struct {
	Code    string
	Message string
	Err     error
}

func (a *AppError) Error() string {
	return a.Code
}

func InternalServerError(message string, err error) *AppError {
	return &AppError{
		Code:    INTERNAL_ERROR,
		Message: message,
		Err:     err,
	}
}

func FileNotFound(message string, err error) *AppError {
	return &AppError{
		Code:    FILE_NOT_FOUND,
		Message: message,
		Err:     err,
	}
}
