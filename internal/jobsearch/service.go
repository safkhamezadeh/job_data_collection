package jobsearch

import (
	"context"
	apperror "job_vacancies/internal/AppError"
	jobvacancies "job_vacancies/internal/job_vacancies"
	"job_vacancies/internal/keywordextractor"
)

type JobSearch struct {
	keywordExtractor keywordextractor.KeywordsExtractor
	jobFinder        jobvacancies.VacancyGetter
}

func NewJobSearch(extractor keywordextractor.KeywordsExtractor, finder jobvacancies.VacancyGetter) *JobSearch {
	return &JobSearch{keywordExtractor: extractor, jobFinder: finder}
}

func (j *JobSearch) Search(ctx context.Context, input string, opt jobvacancies.SearchOptions) ([]jobvacancies.Job, error) {
	if err := validateInput(input); err != nil {
		return nil, err
	}

	keywords, err := j.keywordExtractor.Translate(ctx, input)
	if err != nil {
		return nil, err
	}
	foundJobs, err := j.jobFinder.FindVacancies(keywords, opt)
	if err != nil {
		return nil, err
	}

	return foundJobs, nil
}

func validateInput(input string) error {
	if input == "" {
		return apperror.New("INVALID_INPUT", "no input found")
	}
	if len(input) < MIN_INPUT_LEN {
		return apperror.New("INVALID_INPUT", "input too short, please try to describe what you want to do in the job")
	}
	if len(input) > MAX_INPUT_LEN {
		return apperror.New("INVALID_INPUT", "input too long")
	}

	return nil
}
