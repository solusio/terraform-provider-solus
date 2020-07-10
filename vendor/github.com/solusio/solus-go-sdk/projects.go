package solus

import (
	"context"
	"fmt"
)

type ProjectsService service

type ProjectsResponse struct {
	paginatedResponse

	Data []Project `json:"data"`
}

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

type ProjectCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *ProjectsService) Create(ctx context.Context, data ProjectCreateRequest) (Project, error) {
	var resp struct {
		Data Project `json:"data"`
	}
	return resp.Data, s.client.create(ctx, "projects", data, &resp)
}

func (s *ProjectsService) List(ctx context.Context) (ProjectsResponse, error) {
	resp := ProjectsResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "projects", &resp)
}

func (s *ProjectsService) Delete(ctx context.Context, id int) error {
	return s.client.delete(ctx, fmt.Sprintf("projects/%d", id))
}
