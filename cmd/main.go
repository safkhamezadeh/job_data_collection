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
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatalf("err loading env: %v", err)
	}
	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// -------------------------
	// Infrastructure
	// -------------------------

	cacheImpl := cache.NewInMemoryCache[[]jobvacancies.Job](15 * time.Minute)
	cacheImpl.StartCleanup(ctx, 10*time.Minute)

	keywordExtractor, err := googlegemini.NewGeminiClient(ctx)
	if err != nil {
		log.Fatalf("failed to init Gemini client: %v", err)
	}

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
	h.RegisterRoute(mux)

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

	go func() {
		log.Printf("Server running on :%s", cfg.HTTP.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// wait for interrupt (Ctrl+C)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	srv.Shutdown(shutdownCtx)
	cancel() // stops cache cleanup goroutine
}
