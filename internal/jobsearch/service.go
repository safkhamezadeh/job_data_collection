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
	Set(id string, data []jobvacancies.Job)
	Get(id string) ([]jobvacancies.Job, bool)
}

type JobSearchService struct {
	keywordExtractor keywordextractor.KeywordsExtractor
	jobFinder        jobvacancies.VacancyGetter
	jobRanker        Ranker
	cache            Cache
}

func NewJobSearchService(
	extractor keywordextractor.KeywordsExtractor,
	finder jobvacancies.VacancyGetter,
	ranker Ranker,
	cache Cache,
) *JobSearchService {
	return &JobSearchService{
		keywordExtractor: extractor,
		jobFinder:        finder,
		jobRanker:        ranker,
		cache:            cache,
	}
}

type SearchResults struct {
	Jobs       []jobvacancies.Job
	SessionID  CacheID
	TotalCount int
	Page       int
}

func (j *JobSearchService) Search(ctx context.Context, input string, sessionID CacheID, opt jobvacancies.SearchOptions) (SearchResults, error) {
	// validate first, before any cache or network calls
	if err := validateInput(input); err != nil {
		return SearchResults{}, err
	}
	if err := opt.Validate(); err != nil {
		return SearchResults{}, err
	}

	// cache hit — user is paginating an existing session
	if sessionID != "" {
		if cached, ok := j.cache.Get(string(sessionID)); ok {
			return SearchResults{
				Jobs:       paginate(cached, opt.Page, opt.Limit),
				SessionID:  sessionID,
				TotalCount: len(cached),
				Page:       opt.Page,
			}, nil
		}
	}

	// cache miss — new search
	keywords, err := j.keywordExtractor.Translate(ctx, input)
	if err != nil {
		return SearchResults{}, err
	}
	foundJobs, err := j.jobFinder.FindVacancies(ctx, keywords, opt)
	if err != nil {
		return SearchResults{}, err
	}
	ranked := j.jobRanker.RankJobs(keywords, foundJobs)

	newSessionID := CacheID(uuid.NewString())
	j.cache.Set(string(newSessionID), ranked)

	return SearchResults{
		Jobs:       paginate(ranked, 1, opt.Limit),
		SessionID:  newSessionID,
		TotalCount: len(ranked),
		Page:       1,
	}, nil
}

func paginate(jobs []jobvacancies.Job, page, limit int) []jobvacancies.Job {
	if page < 1 {
		page = 1
	}
	start := (page - 1) * limit
	if start >= len(jobs) {
		return nil
	}
	end := min(start+limit, len(jobs))
	return jobs[start:end]
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
