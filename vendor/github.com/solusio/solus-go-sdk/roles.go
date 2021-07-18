package solus

import (
	"context"
	"fmt"
)

type RolesService service

type Role struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	IsDefault  bool   `json:"is_default"`
	UsersCount int    `json:"users_count"`
}

type RoleCreateRequest struct {
	Name        string `json:"name"`
	Permissions []int  `json:"permissions,omitempty"`
}

type RolesResponse struct {
	paginatedResponse

	Data []Role `json:"data"`
}

type roleResponse struct {
	Data Role `json:"data"`
}

func (s *RolesService) Create(ctx context.Context, data RoleCreateRequest) (Role, error) {
	var resp roleResponse
	return resp.Data, s.client.create(ctx, "roles", data, &resp)
}

func (s *RolesService) List(ctx context.Context) (RolesResponse, error) {
	resp := RolesResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "roles", &resp)
}

func (s *RolesService) Get(ctx context.Context, id int) (Role, error) {
	var resp roleResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("roles/%d", id), &resp)
}

func (s *RolesService) GetByName(ctx context.Context, name string) (Role, error) {
	roles, err := s.List(ctx)
	if err != nil {
		return Role{}, err
	}

	for _, role := range roles.Data {
		if role.Name == name {
			return role, nil
		}
	}

	return Role{}, fmt.Errorf("failed to get role by name %q: role not found", name)
}
