package solus

import (
	"context"
	"fmt"
)

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

type ProjectServersResponse struct {
	paginatedResponse

	Data []Server `json:"data"`
}

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

func (s *ProjectsService) Servers(ctx context.Context, id int) (ProjectServersResponse, error) {
	resp := ProjectServersResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, fmt.Sprintf("projects/%d/servers", id), &resp)
}
