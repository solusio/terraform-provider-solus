package solus

// ResponseLinks represent useful links which is returned with paginated response.
type ResponseLinks struct {
	First string `json:"first"`
	Last  string `json:"last"`
	Prev  string `json:"prev"`
	Next  string `json:"next"`
}

// ResponseMeta represent response metadata which is returned with paginated response.
type ResponseMeta struct {
	CurrentPage int    `json:"current_page"`
	From        int    `json:"from"`
	LastPage    int    `json:"last_page"`
	Path        string `json:"path"`
	PerPage     int    `json:"per_page"`
	To          int    `json:"to"`
	Total       int    `json:"total"`
}

type paginatedResponse struct {
	Links ResponseLinks `json:"links"`
	Meta  ResponseMeta  `json:"meta"`

	err     error
	service *service
}

// Run `go generate` to add `Next()` method to all required structs.

// Err return an error which is occurred during fetching next page data.
func (r *paginatedResponse) Err() error {
	return r.err
}
