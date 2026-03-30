package jobvacancies

import (
	"job_vacancies/internal/keywordextractor"
	"slices"
)

type JobVacancyService struct {
	providers []VacancyGetter
}

func (s *JobVacancyService) Subscribe(p VacancyGetter) {
	if slices.Contains(s.providers, p) {
		return
	}
	s.providers = append(s.providers, p)
}

// todo for later: run jobs per service concurrent
func (s *JobVacancyService) FindVacancies(keywords keywordextractor.KeyWordFormat, opt SearchOptions) ([]Job, error) {
	var allJobs []Job

	for _, provider := range s.providers {
		result, err := provider.FindVacancies(keywords, opt)
		if err != nil {
			continue
		}
		allJobs = append(allJobs, result...)
	}
	return allJobs, nil
}
