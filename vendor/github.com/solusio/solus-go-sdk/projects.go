package solus

import (
	"context"
	"fmt"
)

// ProjectsService handles all available methods with projects.
type ProjectsService service

// Project represents a project.
type Project struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Members     int    `json:"members"`
	IsOwner     bool   `json:"is_owner"`
	IsDefault   bool   `json:"is_default"`
	Owner       User   `json:"owner"`
	Servers     int    `json:"servers"`
}

// ProjectsResponse represents paginated list of projects.
// This cursor can be used for iterating over all available projects.
type ProjectsResponse struct {
	paginatedResponse

	Data []Project `json:"data"`
}

// ProjectCreateRequest represents available properties for creating a new project.
type ProjectCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Create creates new project.
func (s *ProjectsService) Create(ctx context.Context, data ProjectCreateRequest) (Project, error) {
	var resp struct {
		Data Project `json:"data"`
	}
	return resp.Data, s.client.create(ctx, "projects", data, &resp)
}

// List lists projects.
func (s *ProjectsService) List(ctx context.Context) (ProjectsResponse, error) {
	resp := ProjectsResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "projects", &resp)
}

// Delete deletes specified project.
func (s *ProjectsService) Delete(ctx context.Context, id int) error {
	return s.client.syncDelete(ctx, fmt.Sprintf("projects/%d", id))
}
