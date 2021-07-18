package solus

import (
	"context"
)

type PermissionsService service

type Permission struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PermissionResponse struct {
	paginatedResponse

	Data []Permission `json:"data"`
}

func (s *PermissionsService) List(ctx context.Context) (PermissionResponse, error) {
	resp := PermissionResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "permissions", &resp)
}
