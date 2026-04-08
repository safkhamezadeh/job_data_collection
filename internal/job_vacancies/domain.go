package jobvacancies

import (
	"context"
	"fmt"
	"job_vacancies/internal/keywordextractor"
	"time"
)

type CountryISO2 string

type VacancyGetter interface {
	FindVacancies(ctx context.Context, keywords keywordextractor.KeyWordFormat, opt SearchOptions) ([]Job, error)
}

type Job struct {
	Id                  string
	Title               string
	Company             string
	Description         string
	LocationDisplayName string
	Date_posted         time.Time
	Source              string //like "adzuna, indeed"
	Salary_Min          float64
	Salary_Max          float64
	External_url        string
}

func (j Job) String() string {
	return fmt.Sprintf(
		`Job:
  ID:       %s
  Title:    %s
  Company:  %s
  Location: %s
  Posted:   %s
  Source:   %s
  Salary:   %.2f - %.2f
  URL:      %s
  Description: %s`,
		j.Id,
		j.Title,
		j.Company,
		j.LocationDisplayName,
		j.Date_posted.Format("2006-01-02"),
		j.Source,
		j.Salary_Min,
		j.Salary_Max,
		j.External_url,
		j.Description,
	)
}

type Location struct {
	Country CountryISO2
	City    string
	Region  string // state/province
	// Address    string // street + number
	// PostalCode string
}

type SearchOptions struct {
	Location Location
	Limit    int

	// Pagination (best effort)
	Page int
}
