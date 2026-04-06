package adzuna

import (
	"fmt"
	jobvacancies "job_vacancies/internal/job_vacancies"
	"job_vacancies/internal/keywordextractor"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type PathParams struct {
	country Iso2CountryCode
	page    int64
}

func toPathParams(ctry Iso2CountryCode, page int64) PathParams {
	return PathParams{country: ctry, page: page}
}

type QueryParams struct {
	app_id           string
	app_key          string
	results_per_page int64
	what             string //keywords to search for, space seperated
}

func toQueryParams(id string, key string, keywords keywordextractor.KeyWordFormat) QueryParams {
	keywordstring := keywords.ToString(" ")
	return QueryParams{app_id: id, app_key: key, results_per_page: 10, what: keywordstring}
}

func buildPath(base, template string, params PathParams) string {
	path := template
	path = strings.ReplaceAll(path, "{country}", string(params.country))
	path = strings.ReplaceAll(path, "{page}", fmt.Sprintf("%d", params.page))
	return base + path
}

func buildQuery(params QueryParams) string {
	q := url.Values{}
	q.Set("app_id", params.app_id)
	q.Set("app_key", params.app_key)
	q.Set("results_per_page", fmt.Sprintf("%d", params.results_per_page))
	q.Set("what", params.what)
	return q.Encode() // produces: app_id=...&app_key=...&...
}

func buildRequest(pathParams PathParams, queryParams QueryParams) (*http.Request, error) {
	baseURL := "https://api.adzuna.com/v1/api"
	pathTemplate := "/jobs/{country}/search/{page}"

	fullPath := buildPath(baseURL, pathTemplate, pathParams)
	queryString := buildQuery(queryParams)

	finalURL := fullPath + "?" + queryString

	req, err := http.NewRequest("GET", finalURL, nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

type AdzunaResponse struct {
	Results []AdzunaJob `json:"results"`
}

type AdzunaJob struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	RedirectURL string         `json:"redirect_url"`
	Created     time.Time      `json:"created"`
	SalaryMin   float64        `json:"salary_min"`
	SalaryMax   float64        `json:"salary_max"`
	Company     AdzunaCompany  `json:"company"`
	Location    AdzunaLocation `json:"location"`
}

type AdzunaCompany struct {
	DisplayName string `json:"display_name"`
}

type AdzunaLocation struct {
	DisplayName string `json:"display_name"`
}

func mapToJobs(resp AdzunaResponse) []jobvacancies.Job {
	var jobs []jobvacancies.Job

	for _, j := range resp.Results {
		jobs = append(jobs, jobvacancies.Job{
			Id:                  j.ID,
			Title:               j.Title,
			Company:             j.Company.DisplayName,
			Description:         j.Description,
			LocationDisplayName: j.Location.DisplayName,
			Date_posted:         j.Created,
			External_url:        j.RedirectURL,
			Source:              "Adzuna",
			Salary_Min:          j.SalaryMin,
			Salary_Max:          j.SalaryMax,
		})
	}

	return jobs
}
