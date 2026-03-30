package jobsearch

import "net/http"

type Handler struct {
	jobSearchService *JobSearch
}

func (h *Handler) HandleFindJobs(w http.ResponseWriter, r *http.Request) {
	//parse req
	//call service
	//return response
}
