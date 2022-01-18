package solus

import (
	"context"
)

// PermissionsService handles all available methods with permissions.
type PermissionsService service

// Permission represents a permission.
type Permission struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// PermissionResponse represents paginated list of permissions.
// This cursor can be used for iterating over all available permissions.
type PermissionResponse struct {
	paginatedResponse

	Data []Permission `json:"data"`
}

// List lists permissions.
func (s *PermissionsService) List(ctx context.Context) (PermissionResponse, error) {
	resp := PermissionResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "permissions", &resp)
}
