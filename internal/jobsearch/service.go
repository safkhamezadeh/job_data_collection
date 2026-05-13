package jobsearch

import (
	"context"
	"errors"
	"fmt"
	jobvacancies "job_vacancies/internal/job_vacancies"
	"job_vacancies/internal/keywordextractor"
	"log"

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

func (j *JobSearchService) Search(
	ctx context.Context,
	input string,
	sessionID CacheID,
	opt jobvacancies.SearchOptions,
) (SearchResults, error) {

	log.Printf("[Search] start input=%s session=%s page=%d limit=%d",
		input, sessionID, opt.Page, opt.Limit)

	// validate
	if err := validateInput(input); err != nil {
		log.Printf("[Search] input validation failed: %v", err)
		return SearchResults{}, err
	}
	if err := opt.Validate(); err != nil {
		log.Printf("[Search] options validation failed: %v", err)
		return SearchResults{}, err
	}

	// cache hit
	if sessionID != "" {
		if cached, ok := j.cache.Get(string(sessionID)); ok {
			log.Printf("[Search] cache HIT session=%s total_cached=%d",
				sessionID, len(cached))

			page := paginate(cached, opt.Page, opt.Limit)

			return SearchResults{
				Jobs:       page,
				SessionID:  sessionID,
				TotalCount: len(cached),
				Page:       opt.Page,
			}, nil
		}

		log.Printf("[Search] cache MISS session=%s", sessionID)
	}

	// keyword extraction
	keywords, err := j.keywordExtractor.Translate(ctx, input)
	if err != nil {
		log.Printf("[Search] keyword extraction failed: %v", err)
		return SearchResults{}, err
	}

	if ctx.Err() != nil {
		return SearchResults{}, fmt.Errorf("translate aborted: %w", ctx.Err())
	}

	// job fetch
	foundJobs, err := j.jobFinder.FindVacancies(ctx, keywords, opt)
	if err != nil {
		log.Printf("[Search] jobFinder failed: %v", err)
		return SearchResults{}, err
	}

	if ctx.Err() != nil {
		return SearchResults{}, fmt.Errorf("search aborted: %w", ctx.Err())
	}

	// ranking
	ranked := j.jobRanker.RankJobs(keywords, foundJobs)

	// cache set
	newSessionID := CacheID(uuid.NewString())
	j.cache.Set(string(newSessionID), ranked)

	// pagination
	page := paginate(ranked, 1, opt.Limit)

	return SearchResults{
		Jobs:       page,
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
