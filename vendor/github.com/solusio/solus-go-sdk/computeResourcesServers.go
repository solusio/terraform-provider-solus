package solus

import (
	"context"
	"fmt"
)

// ComputeResourceServerCreateRequest represents available properties for creating
// a new servers on a compute resource.
type ComputeResourceServerCreateRequest struct {
	Name             string                       `json:"name"`
	Description      string                       `json:"description"`
	Password         string                       `json:"password"`
	PlanID           int                          `json:"plan_id"`
	OSImageVersionID int                          `json:"os_image_version_id,omitempty"`
	ApplicationID    int                          `json:"application_id,omitempty"`
	ApplicationData  map[string]string            `json:"application_data,omitempty"`
	SSHKeys          []int                        `json:"ssh_keys,omitempty"`
	UserData         string                       `json:"user_data,omitempty"`
	FQDNs            []string                     `json:"fqdns,omitempty"`
	UserID           int                          `json:"user_id"`
	ProjectID        int                          `json:"project_id"`
	BackupSettings   *VirtualServerBackupSettings `json:"backup_settings,omitempty"`
}

// ServersCreate creates a new server on the specified compute resource.
func (s *ComputeResourcesService) ServersCreate(
	ctx context.Context,
	id int,
	data ComputeResourceServerCreateRequest,
) (VirtualServer, error) {
	var resp virtualServerResponse
	return resp.Data, s.client.create(ctx, fmt.Sprintf("compute_resources/%d/servers", id), data, &resp)
}
