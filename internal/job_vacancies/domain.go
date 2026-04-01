package jobvacancies

import (
	"job_vacancies/internal/keywordextractor"
	"time"
)

type CountryISO2 string

type VacancyGetter interface {
	FindVacancies(keywords keywordextractor.KeyWordFormat, opt SearchOptions) ([]Job, error)
}

type Job struct {
	Id           string
	Title        string
	Company      string
	Description  string
	Location     Location
	Date_posted  time.Time
	Source       string //like "adzuna, indeed"
	Salary_Min   float64
	Salary_Max   float64
	External_url string
}

type Location struct {
	Country    CountryISO2
	City       string
	Region     string // state/province
	Address    string // street + number
	PostalCode string
}

type SearchOptions struct {
	Country CountryISO2
	Limit   int

	// Pagination (best effort)
	Page int
}
