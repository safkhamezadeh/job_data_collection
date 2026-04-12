package jobsearch

import jobvacancies "job_vacancies/internal/job_vacancies"

type UserInputDTO struct {
	Input     string       `json:"input"`
	SearchOpt SearchOptDTO `json:"search_options,omitempty"`
}

type SearchOptDTO struct {
	Country string `json:"country,omitempty"`
	Limit   int    `json:"limit,omitempty"`
	Page    int    `json:"page,omitempty"`
}

func (s SearchOptDTO) toSearchOpt() jobvacancies.SearchOptions {
	loc := jobvacancies.Location{Country: jobvacancies.CountryISO2(s.Country)}

	return jobvacancies.SearchOptions{Location: loc,
		Limit: s.Limit,
		Page:  s.Page}
}

type JobDTO struct {
}

type SearchResponse struct {
	Jobs     []JobDTO `json:"jobs"`
	Warnings []string `json:"warnings,omitempty"` // optional
}
