package jobsearch

import (
	"errors"
	"time"

	jobvacancies "job_vacancies/internal/job_vacancies"
	"job_vacancies/internal/location"
)

type UserInputDTO struct {
	Input     string       `json:"input"`
	SessionID string       `json:"session_id"`
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
	Id                  string    `json:"id"`
	Title               string    `json:"title"`
	Company             string    `json:"company"`
	Description         string    `json:"description"`
	LocationDisplayName string    `json:"location_display_name"`
	Date_posted         time.Time `json:"date_posted"`
	Source              string    `json:"source"` // like "adzuna, indeed"
	Salary_Min          float64   `json:"salary_min"`
	Salary_Max          float64   `json:"salary_max"`
	External_url        string    `json:"external_url"`
}

func ToJobDTO(job jobvacancies.Job) JobDTO {
	return JobDTO{
		Id:                  job.Id,
		Title:               job.Title,
		Company:             job.Company,
		Description:         job.Description,
		LocationDisplayName: job.LocationDisplayName,
		Date_posted:         job.Date_posted,
		Source:              job.Source,
		Salary_Min:          job.Salary_Min,
		Salary_Max:          job.Salary_Max,
		External_url:        job.External_url,
	}
}

func ToJobDTOs(jobs []jobvacancies.Job) []JobDTO {
	jobDTOs := make([]JobDTO, len(jobs))

	for i, job := range jobs {
		jobDTOs[i] = JobDTO{
			Id:                  job.Id,
			Title:               job.Title,
			Company:             job.Company,
			Description:         job.Description,
			LocationDisplayName: job.LocationDisplayName,
			Date_posted:         job.Date_posted,
			Source:              job.Source,
			Salary_Min:          job.Salary_Min,
			Salary_Max:          job.Salary_Max,
			External_url:        job.External_url,
		}
	}

	return jobDTOs
}

type SearchResponseDTO struct {
	SessionKey CacheID  `json:"session_key"`
	Jobs       []JobDTO `json:"jobs"`
	TotalCount int      `json:"total_count"`
	Page       int      `json:"page"`
}

func toSearchResponseDTO(res SearchResults) SearchResponseDTO {
	return SearchResponseDTO{
		SessionKey: res.SessionID,
		Jobs:       ToJobDTOs(res.Jobs),
		TotalCount: res.TotalCount,
		Page:       res.Page,
	}
}
