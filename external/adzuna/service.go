package adzuna

import (
	"context"
	"encoding/json"
	"fmt"
	"job_vacancies/config"
	jobvacancies "job_vacancies/internal/job_vacancies"
	"job_vacancies/internal/keywordextractor"
	"log"
	"net/http"
)

type adzunaClient struct {
	applicationKey string
	applicationID  string
	httpClient     *http.Client
}

func NewAdzunaClient(keys config.AdzunaKeys, httpClient *http.Client) *adzunaClient {
	return &adzunaClient{
		applicationKey: keys.ApplicationKey,
		applicationID:  keys.ApplicationID,
		httpClient:     httpClient,
	}
}

func (a *adzunaClient) FindVacancies(ctx context.Context, keywords keywordextractor.KeyWordFormat, opt jobvacancies.SearchOptions) ([]jobvacancies.Job, error) {
	if !IsWhitelisted(Iso2CountryCode(opt.Location.Country)) {
		return nil, jobvacancies.ErrInvalidCountry
	}

	pathParams := toPathParams(Iso2CountryCode(opt.Location.Country), int64(opt.Page))

	queryParams := toQueryParams(a.applicationID, a.applicationKey, int64(opt.Limit), keywords)

	req, err := buildRequest(ctx, pathParams, queryParams)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		log.Printf("FindVacancies Adzuna error: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Adzuna returned status: %d", resp.StatusCode)
		return nil, fmt.Errorf("adzuna: failed to fetch vacancies")

	}

	var adzunaResp AdzunaResponse
	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&adzunaResp); err != nil {
		return nil, fmt.Errorf("failed to parse response")

	}

	jobs := mapToJobs(adzunaResp)

	return jobs, nil
}
