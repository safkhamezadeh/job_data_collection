package jobsearch

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	jobSearchService *JobSearchService
}

func (h *Handler) HandleFindJobs(w http.ResponseWriter, r *http.Request) {
	//parse req
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

	result, err := h.jobSearchService.Search(ctx, body.Input, CacheID(body.SessionID), body.SearchOpt.toSearchOpt())
	if err != nil {
		http.Error(w, "Error finding Jobs", http.StatusInternalServerError)
		log.Printf("Err HandleFindJobs .Search, Err: %v", err)
	}

	print(result)
	//return response
}
