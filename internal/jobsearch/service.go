package jobsearch

import (
	"context"
	"errors"
	jobvacancies "job_vacancies/internal/job_vacancies"
	"job_vacancies/internal/keywordextractor"
)

const (
	MIN_INPUT_LEN int = 8
	MAX_INPUT_LEN int = 300
)

type JobSearch struct {
	keywordExtractor keywordextractor.KeywordsExtractor
	jobFinder        jobvacancies.VacancyGetter
	//jobranker        ranking.SimpleRanker
	//cache            *cache.Cache
}

func NewJobSearch(extractor keywordextractor.KeywordsExtractor, finder jobvacancies.VacancyGetter) *JobSearch {
	return &JobSearch{keywordExtractor: extractor, jobFinder: finder}
}

func (j *JobSearch) Search(ctx context.Context, input string, opt jobvacancies.SearchOptions) ([]jobvacancies.Job, error) {
	if err := validateInput(input); err != nil {
		return nil, err
	}

	//validate opt
	if err := opt.Validate(); err != nil {
		return nil, err
	}

	keywords, err := j.keywordExtractor.Translate(ctx, input)
	if err != nil {
		return nil, err
	}
	foundJobs, err := j.jobFinder.FindVacancies(ctx, keywords, opt)
	if err != nil {
		return nil, err
	}

	//do ranking

	return foundJobs, nil
}

func validateInput(input string) error {
	if input == "" {
		return errors.New("no input found")
	}
	if len(input) < MIN_INPUT_LEN {
		return errors.New("input too short")
	}
	if len(input) > MAX_INPUT_LEN {
		return errors.New("input too long")
	}

	return nil
}
