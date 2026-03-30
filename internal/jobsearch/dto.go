package jobsearch

type UserInputDTO struct {
	Input     string       `json:"input"`
	SearchOpt SearchOptDTO `json:"search_options,omitempty"`
}

type SearchOptDTO struct {
	Country string `json:"country,omitempty"`
	Limit   int    `json:"limit,omitempty"`
	Page    int    `json:"page,omitempty"`
}

type JobDTO struct {
}

type SearchResponse struct {
	Jobs     []JobDTO `json:"jobs"`
	Warnings []string `json:"warnings,omitempty"` // optional
}
