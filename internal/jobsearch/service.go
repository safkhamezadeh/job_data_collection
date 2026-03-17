package jobsearch

import (
	jobvacancies "job_vacancies/internal/job_vacancies"
	"job_vacancies/internal/keywordextractor"
)

type JobSearch struct {
	keywordextracter keywordextractor.KeywordsExtractor
	jobFinder        jobvacancies.VacancyGetter
}
