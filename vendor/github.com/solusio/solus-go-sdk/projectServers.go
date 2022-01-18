package solus

import (
	"context"
	"fmt"
)

// ProjectServersCreateRequest represents available properties for creating a new
// server on a project.
type ProjectServersCreateRequest struct {
	Name             string                `json:"name"`
	PlanID           int                   `json:"plan_id"`
	LocationID       int                   `json:"location_id"`
	OsImageVersionID int                   `json:"os_image_version_id,omitempty"`
	ApplicationID    int                   `json:"application_id,omitempty"`
	ApplicationData  string                `json:"application_data,omitempty"`
	SSHKeys          []int                 `json:"ssh_keys,omitempty"`
	UserData         string                `json:"user_data,omitempty"`
	FQDNs            []string              `json:"fqdns,omitempty"`
	BackupSettings   *ServerBackupSettings `json:"backup_settings,omitempty"`
}

// ProjectServersResponse represents paginated list of project's servers.
// This cursor can be used for iterating over all available project's servers.
type ProjectServersResponse struct {
	paginatedResponse

	Data []Server `json:"data"`
}

// ServersCreate creates a server on the specified project.
func (s *ProjectsService) ServersCreate(
	ctx context.Context,
	projectID int,
	data ProjectServersCreateRequest,
) (Server, error) {
	var resp struct {
		Data Server `json:"data"`
	}
	return resp.Data, s.client.create(ctx, fmt.Sprintf("projects/%d/servers", projectID), data, &resp)
}

// ServersListAll lists all servers on the specified project.
// Deprecated: use Servers instead.
func (s *ProjectsService) ServersListAll(ctx context.Context, id int) ([]Server, error) {
	resp, err := s.Servers(ctx, id)
	if err != nil {
		return nil, err
	}

	servers := make([]Server, len(resp.Data))
	copy(servers, resp.Data)
	for resp.Next(ctx) {
		servers = append(servers, resp.Data...)
	}
	return servers, resp.Err()
}

// Servers lists all servers on the specified project.
func (s *ProjectsService) Servers(ctx context.Context, id int) (ProjectServersResponse, error) {
	resp := ProjectServersResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, fmt.Sprintf("projects/%d/servers", id), &resp)
}
