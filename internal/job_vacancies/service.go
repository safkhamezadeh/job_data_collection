package jobvacancies

import (
	apperror "job_vacancies/internal/AppError"
	"job_vacancies/internal/keywordextractor"
)

type JobVacancyService struct {
	providers []VacancyGetter
}

func (s *JobVacancyService) Subscribe(p VacancyGetter) {
	s.providers = append(s.providers, p)
}

// todo for later: run jobs per service concurrent
func (s *JobVacancyService) FindVacancies(keywords []keywordextractor.KeyWordFormat, opt SearchOptions) ([]Job, apperror.AppError) {

	return nil, apperror.AppError{}
}
