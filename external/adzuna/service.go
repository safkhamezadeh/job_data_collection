package adzuna

import (
	"encoding/json"
	apperror "job_vacancies/internal/AppError"
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
	if !IsWhitelisted(Iso2CountryCode(opt.Location.Country)) {
		return nil, jobvacancies.ErrInvalidCountry
	}

	pathParams := toPathParams(Iso2CountryCode(opt.Location.Country), int64(opt.Page))

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

	if resp.StatusCode != http.StatusOK {
		log.Printf("Adzuna returned status: %d", resp.StatusCode)
		return nil, apperror.AppError{
			Code:    "adzuna_error",
			Message: "failed to fetch vacancies",
		}
	}

	var adzunaResp AdzunaResponse
	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&adzunaResp); err != nil {
		log.Printf("Error decoding Adzuna response: %v", err)
		return nil, apperror.AppError{
			Code:    "parse_error",
			Message: "failed to parse response",
		}
	}

	jobs := mapToJobs(adzunaResp)

	return jobs, nil
}
