package solus

import (
	"context"
)

type ServersMigrationsService service

type ServersMigration struct {
	ID                         int             `json:"id"`
	DestinationComputeResource ComputeResource `json:"destination_compute_resource"`
	Task                       Task            `json:"task"`
	Children                   []Task          `json:"children"`
}

type ServersMigrationRequest struct {
	IsLive                       bool  `json:"is_live"`
	PreserveIPs                  bool  `json:"preserve_ips"`
	DestinationComputeResourceID int   `json:"destination_compute_resource_id"`
	Servers                      []int `json:"servers"`
}

func (s *ServersMigrationsService) Create(ctx context.Context, data ServersMigrationRequest) (ServersMigration, error) {
	var resp struct {
		Data ServersMigration `json:"data"`
	}
	return resp.Data, s.client.create(ctx, "servers_migrations", data, &resp)
}
