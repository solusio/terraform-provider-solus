package solus

import (
	"context"
	"fmt"
)

// RolesService handles all available methods with roles.
type RolesService service

// Role represents a role.
type Role struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	IsDefault  bool   `json:"is_default"`
	UsersCount int    `json:"users_count"`
}

// RoleCreateRequest represents available properties for creating a new role.
type RoleCreateRequest struct {
	Name        string `json:"name"`
	Permissions []int  `json:"permissions,omitempty"`
}

// RolesResponse represents paginated list of roles.
// This cursor can be used for iterating over all available roles.
type RolesResponse struct {
	paginatedResponse

	Data []Role `json:"data"`
}

type roleResponse struct {
	Data Role `json:"data"`
}

// Create creates new role.
func (s *RolesService) Create(ctx context.Context, data RoleCreateRequest) (Role, error) {
	var resp roleResponse
	return resp.Data, s.client.create(ctx, "roles", data, &resp)
}

// List lists roles.
func (s *RolesService) List(ctx context.Context) (RolesResponse, error) {
	resp := RolesResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "roles", &resp)
}

// Get gets specified role.
func (s *RolesService) Get(ctx context.Context, id int) (Role, error) {
	var resp roleResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("roles/%d", id), &resp)
}

// GetByName gets specified role by name.
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
