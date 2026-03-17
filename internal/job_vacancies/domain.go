package jobvacancies

import (
	"job_vacancies/internal/keywordextractor"
	"time"
)

type VacancyGetter interface {
	FindVacancies(keywords []keywordextractor.KeyWordFormat, opt SearchOptions) ([]Job, error)
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

type SearchOptions struct {
	Country string
	Limit   int

	// Pagination (best effort)
	Page int
}
