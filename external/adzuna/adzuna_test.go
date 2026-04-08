package adzuna_test

import (
	"context"
	"job_vacancies/config"
	"job_vacancies/external/adzuna"
	jobvacancies "job_vacancies/internal/job_vacancies"
	"job_vacancies/internal/keywordextractor"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestFindVacancies_integration(t *testing.T) {
	if os.Getenv("DEV") != "true" {
		t.Skip("skipping integration test (DEV != true)")
	}

	if err := godotenv.Load("./.env.test"); err != nil {
		t.Fatalf("failed to load .env.test: %v", err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	adzunaKeys := config.AdzunaKeys{
		ApplicationKey: os.Getenv("ADZUNA_API_APPLICATION_KEY"),
		ApplicationID:  os.Getenv("ADZUNA_API_APPLICATION_ID"),
	}
	httpClient := &http.Client{Timeout: 5 * time.Second}
	adzunaClient := adzuna.NewAdzunaClient(adzunaKeys, httpClient)

	searchOpt := jobvacancies.SearchOptions{
		Location: jobvacancies.Location{Country: "NL"},
		Limit:    3,
		Page:     1,
	}

	t.Log("Starting Adzuna integration test with keywords:", makeKeywords())
	jobs, err := adzunaClient.FindVacancies(ctx, makeKeywords(), searchOpt)
	if err != nil {
		t.Fatalf("FindVacancies returned an error: %v", err)
	}

	t.Logf("Received %d jobs from Adzuna", len(jobs))
	if len(jobs) != 3 {
		t.Errorf("Expected 3 jobs, got %d", len(jobs))
	}

	for i, job := range jobs {
		t.Logf("Job %d:\n%s", i+1, job)
	}
}

func makeKeywords() keywordextractor.KeyWordFormat {
	return keywordextractor.KeyWordFormat{
		JobTitles: []string{
			"Software Engineer",
			"Backend Developer",
		},
		Keywords: []string{"golang"},
	}
}
