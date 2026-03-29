package keywordextractor

import (
	apperror "job_vacancies/internal/AppError"
)

var (
	ExternalErr      apperror.AppError = apperror.AppError{Code: "EXTERNAL_SERVER_ERROR", Message: "something went wrong on the external server"}
	ErrInvalidOutput apperror.AppError = apperror.AppError{Code: "INVALID_OUTPUT", Message: "something went wrong with receiving output"}
)
