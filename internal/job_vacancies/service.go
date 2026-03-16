package jobvacancies

import apperror "job_vacancies/internal/AppError"

type jobVacancyService struct {
	providers []VacancyGetter
}

func (s *jobVacancyService) Subscribe(p VacancyGetter) {
	s.providers = append(s.providers, p)
}

//todo: run jobs per service concurrent
func (s *jobVacancyService) FindVacancies(keywords []string) ([]Job, apperror.AppError) {
	var jobs []Job

	for _, provider := range s.providers {
		result, err := provider.FindVacancies(keywords)
		if err.Code != "" {
			continue
		}

		jobs = append(jobs, result...)
	}

	return jobs, apperror.AppError{}
}
