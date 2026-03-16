package jobvacancies

import (
	apperror "job_vacancies/internal/AppError"
	"time"
)

type VacancyGetter interface {
	FindVacancies(keywords []string) ([]Job, apperror.AppError)
}

type Job struct {
	Id           string
	Title        string
	Company      string
	Description  string
	Location     string
	Date_posted  time.Time
	Source       string //like "adzuna, indeed"
	Salary       string
	External_url string
}
