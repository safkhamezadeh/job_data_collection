package jobsearch

import (
	"context"
	"errors"
	jobvacancies "job_vacancies/internal/job_vacancies"
	"job_vacancies/internal/keywordextractor"

	"github.com/google/uuid"
)

const (
	MIN_INPUT_LEN int = 8
	MAX_INPUT_LEN int = 300
)

type CacheID string

type Ranker interface {
	RankJobs(keywords keywordextractor.KeyWordFormat, jobs []jobvacancies.Job) []jobvacancies.Job
}

type Cache interface {
	Set(id CacheID, data []jobvacancies.Job)
	Get(id CacheID) ([]jobvacancies.Job, error)
}

type JobSearch struct {
	keywordExtractor keywordextractor.KeywordsExtractor
	jobFinder        jobvacancies.VacancyGetter
	jobRanker        Ranker
	cache            Cache
}

func NewJobSearch(extractor keywordextractor.KeywordsExtractor, finder jobvacancies.VacancyGetter) *JobSearch {
	return &JobSearch{keywordExtractor: extractor, jobFinder: finder}
}

func (j *JobSearch) Search(ctx context.Context, input string, sessionID CacheID, opt jobvacancies.SearchOptions) ([]jobvacancies.Job, CacheID, error) {
	if err := validateInput(input); err != nil {
		return nil, "", err
	}
	//validate opt
	if err := opt.Validate(); err != nil {
		return nil, "", err
	}

	keywords, err := j.keywordExtractor.Translate(ctx, input)
	if err != nil {
		return nil, "", err
	}
	foundJobs, err := j.jobFinder.FindVacancies(ctx, keywords, opt)
	if err != nil {
		return nil, "", err
	}

	ranked := j.jobRanker.RankJobs(keywords, foundJobs)
	cacheID := uuid.NewString()
	j.cache.Set(CacheID(cacheID), ranked)

	return ranked, CacheID(cacheID), nil
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
