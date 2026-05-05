package jobsearch

import (
	"errors"

	jobvacancies "job_vacancies/internal/job_vacancies"
	"job_vacancies/internal/location"
)

type UserInputDTO struct {
	Input     string       `json:"input"`
	CacheId   CacheID      `json:"session_id"`
	SearchOpt SearchOptDTO `json:"search_options,omitempty"`
}

func (u UserInputDTO) Validate() error {
	if u.Input == "" {
		return errors.New("input is required")
	}
	return u.SearchOpt.validate()
}

type SearchOptDTO struct {
	CountryISO2 string `json:"country,omitempty"`
	Limit       int    `json:"limit,omitempty"`
	Page        int    `json:"page,omitempty"`
}

func (s SearchOptDTO) validate() error {
	if s.CountryISO2 == "" {
		return errors.New("country is mandatory")
	}
	if s.Limit == 0 {
		return errors.New("limit is mandatory")
	}
	if s.Page == 0 {
		return errors.New("page is mandatory")
	}
	return nil
}

func (s SearchOptDTO) toSearchOpt() jobvacancies.SearchOptions {
	loc := jobvacancies.Location{Country: location.CountryISO2(s.CountryISO2)}

	return jobvacancies.SearchOptions{Location: loc,
		Limit: s.Limit,
		Page:  s.Page}
}

type JobDTO struct {
}

type SearchResponse struct {
	SessionKey CacheID  `json:"session_key"`
	Jobs       []JobDTO `json:"jobs"`
	TotalCount int      `json:"total_count"`
	Page       int      `json:"page"`
}
