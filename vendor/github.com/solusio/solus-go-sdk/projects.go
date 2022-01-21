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

// ProjectRequest represents available properties for creating or updating a project.
type ProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Create creates new project.
func (s *ProjectsService) Create(ctx context.Context, data ProjectRequest) (Project, error) {
	var resp projectResponse
	return resp.Data, s.client.create(ctx, "projects", data, &resp)
}

// List lists projects.
func (s *ProjectsService) List(ctx context.Context, filter *FilterProjects) (ProjectsResponse, error) {
	resp := ProjectsResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "projects", &resp, withFilter(filter.data))
}

// Get gets specified project.
func (s *ProjectsService) Get(ctx context.Context, id int) (Project, error) {
	var resp projectResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("projects/%d", id), &resp)
}

// Update updates specified project.
func (s *ProjectsService) Update(ctx context.Context, id int, data ProjectRequest) (Project, error) {
	var resp projectResponse
	return resp.Data, s.client.update(ctx, fmt.Sprintf("projects/%d", id), data, &resp)
}

// Delete deletes specified project.
func (s *ProjectsService) Delete(ctx context.Context, id int) error {
	return s.client.syncDelete(ctx, fmt.Sprintf("projects/%d", id))
}

type projectResponse struct {
	Data Project `json:"data"`
}
