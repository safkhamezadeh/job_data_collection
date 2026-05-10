package main

import (
	"context"
	"job_vacancies/config"
	"job_vacancies/external/adzuna"
	googlegemini "job_vacancies/external/googleGemini"
	"job_vacancies/internal/infrastructure/cache"
	jobvacancies "job_vacancies/internal/job_vacancies"
	"job_vacancies/internal/jobsearch"
	"job_vacancies/internal/ranker"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatalf("err loading env: %v", err)
	}
	cfg := config.Load()

	ctx := context.Background()
	// -------------------------
	// Infrastructure
	// -------------------------

	cacheImpl := cache.NewInMemoryCache[[]jobvacancies.Job](15 * time.Minute)

	keywordExtractor, err := googlegemini.NewGeminiClient(ctx)

	httpClient := &http.Client{Timeout: 10 * time.Second}
	jobFinder := adzuna.NewAdzunaClient(cfg.APIKeys.JobVacancyKeys.AdzunaKeys, httpClient)

	jobRanker := ranker.NewJobRanker()

	// -------------------------
	// Service layer
	// -------------------------

	jobSearchService := jobsearch.NewJobSearchService(
		keywordExtractor,
		jobFinder,
		jobRanker,
		cacheImpl,
	)

	// -------------------------
	// HTTP layer
	// -------------------------

	h := jobsearch.Handler{
		JobSearchService: jobSearchService,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/jobs/search", h.HandleFindJobs)

	// -------------------------
	// Server
	// -------------------------

	srv := &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server running on :%s", cfg.HTTP.Port)
	log.Fatal(srv.ListenAndServe())
}
