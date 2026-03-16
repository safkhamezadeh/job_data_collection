package apperror

import "fmt"

type AppError struct {
	Code    string
	Message string
}

func (a AppError) Error() string {
	return fmt.Sprintf("errorCode: %s, errorMessage: %s", a.Code, a.Message)
}
