package adzuna

import (
	jobvacancies "job_vacancies/internal/job_vacancies"
	"job_vacancies/internal/keywordextractor"
	"log"
	"net/http"
)

type adzunaClient struct {
	Adzuna_application_key string
	Adzuna_application_id  string
	Http_Client            *http.Client
}

func NewAdzunaClient(key string, id string, httpClient *http.Client) *adzunaClient {
	return &adzunaClient{Adzuna_application_key: key, Adzuna_application_id: id}
}

func (a *adzunaClient) FindVacancies(keywords keywordextractor.KeyWordFormat, opt jobvacancies.SearchOptions) ([]jobvacancies.Job, error) {
	if !IsWhitelisted(Iso2CountryCode(opt.Country)) {
		return nil, jobvacancies.ErrInvalidCountry
	}

	pathParams := toPathParams(Iso2CountryCode(opt.Country), int64(opt.Page))

	queryParams := toQueryParams(a.Adzuna_application_id, a.Adzuna_application_key, keywords)

	req, err := buildRequest(pathParams, queryParams)
	if err != nil {
		return nil, err
	}
	resp, err := a.Http_Client.Do(req)
	if err != nil {
		log.Printf("FindVacancies Adzuna error: %v", err)
		return nil, err
	}

	//todo read resp into jobs

}
