package jobsearch

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	JobSearchService *JobSearchService
}

func (h *Handler) RegisterRoute(mux *http.ServeMux) {
	mux.HandleFunc("/jobs/search", http.HandlerFunc(h.HandleFindJobs))
}

func (h *Handler) HandleFindJobs(w http.ResponseWriter, r *http.Request) {
	// parse request
	var body UserInputDTO
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = body.Validate()
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	result, err := h.JobSearchService.Search(
		ctx,
		body.Input,
		CacheID(body.SessionID),
		body.SearchOpt.toSearchOpt(),
	)
	if err != nil {
		log.Printf("Err HandleFindJobs .Search, Err: %v", err)
		http.Error(w, "Error finding jobs", http.StatusInternalServerError)
		return
	}

	resDTO := toSearchResponseDTO(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(resDTO)
	if err != nil {
		log.Printf("Err encoding response: %v", err)
	}
}
