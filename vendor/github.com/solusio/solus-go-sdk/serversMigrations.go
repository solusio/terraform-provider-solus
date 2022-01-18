package solus

import (
	"context"
)

// ServersMigrationsService handles all available methods with server's migrations.
type ServersMigrationsService service

// ServersMigration represents a server's migration.
type ServersMigration struct {
	ID                         int             `json:"id"`
	DestinationComputeResource ComputeResource `json:"destination_compute_resource"`
	Task                       Task            `json:"task"`
	Children                   []Task          `json:"children"`
}

// ServersMigrationRequest represents available properties for creating a new
// server's migration.
type ServersMigrationRequest struct {
	IsLive                       bool  `json:"is_live"`
	PreserveIPs                  bool  `json:"preserve_ips"`
	DestinationComputeResourceID int   `json:"destination_compute_resource_id"`
	Servers                      []int `json:"servers"`
}

// Create creates new server's migration.
func (s *ServersMigrationsService) Create(ctx context.Context, data ServersMigrationRequest) (ServersMigration, error) {
	var resp struct {
		Data ServersMigration `json:"data"`
	}
	return resp.Data, s.client.create(ctx, "servers_migrations", data, &resp)
}
