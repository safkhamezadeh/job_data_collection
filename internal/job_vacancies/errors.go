package jobvacancies

import (
	apperror "job_vacancies/internal/AppError"
)

var (
	ErrInvalidCountry apperror.AppError = apperror.AppError{Code: "INVALID_COUNTRY_ERROR", Message: "selected invalid country"}
)
