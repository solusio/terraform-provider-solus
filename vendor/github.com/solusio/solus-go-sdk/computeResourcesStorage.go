package solus

import (
	"context"
	"fmt"
)

type ComputeResourceStorageCreateRequest struct {
	TypeID                  int    `json:"type_id"`
	Path                    string `json:"path"`
	ThinPool                string `json:"thin_pool,omitempty"`
	IsAvailableForBalancing bool   `json:"is_available_for_balancing"`
}

func (s *ComputeResourcesService) StorageCreate(
	ctx context.Context,
	id int,
	data ComputeResourceStorageCreateRequest,
) (Storage, error) {
	var resp storageResponse
	return resp.Data, s.client.create(ctx, fmt.Sprintf("compute_resources/%d/storages", id), data, &resp)
}

func (s *ComputeResourcesService) StorageList(ctx context.Context, id int) ([]Storage, error) {
	resp := struct {
		Data []Storage `json:"data"`
	}{}
	return resp.Data, s.client.get(ctx, fmt.Sprintf("compute_resources/%d/storages", id), &resp)
}
