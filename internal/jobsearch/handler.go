package jobsearch

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type Handler struct {
	jobSearchService *JobSearch
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

	jobs, err := h.jobSearchService.Search(ctx, body.Input, body.SearchOpt.toSearchOpt())
	if err != nil {

	}
	//call service
	//return response
}
