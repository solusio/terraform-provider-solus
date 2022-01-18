package solus

import (
	"context"
	"fmt"
)

// BackupNodesService handles all available methods with backup nodes.
type BackupNodesService service

// BackupNode represents a backup node.
// The backup node is a server or a service where backup can be stored.
type BackupNode struct {
	ID                    int                    `json:"id"`
	Name                  string                 `json:"name"`
	Type                  BackupNodeType         `json:"type"`
	Credentials           map[string]interface{} `json:"credentials"`
	ComputeResourcesCount int                    `json:"compute_resources_count"`
	BackupsCount          int                    `json:"backups_count"`
	TotalBackupsSize      int                    `json:"total_backups_size"`
	ComputeResources      []ComputeResource      `json:"compute_resources"`
}

// BackupNodeType a backup node type.
type BackupNodeType string

const (
	// BackupNodeTypeSSHRsync uses SSH with rsync to manipulate backups.
	BackupNodeTypeSSHRsync BackupNodeType = "ssh_rsync"

	// BackupNodeTypeHetznerStorageBox Hetzner Storage Box specific backup node
	// type.
	// https://docs.hetzner.com/robot/storage-box/general
	BackupNodeTypeHetznerStorageBox BackupNodeType = "hetzner_storage_box"
)

// BackupNodeRequest represents available properties for creating new or updating
// existing backup nodes.
type BackupNodeRequest struct {
	Name             string                 `json:"name"`
	Type             BackupNodeType         `json:"type"`
	ComputeResources []int                  `json:"compute_resources,omitempty"`
	Credentials      map[string]interface{} `json:"credentials,omitempty"`
}

// BackupNodeSSHRsyncCredentials creates SSH+Rsync specific connection credentials.
func BackupNodeSSHRsyncCredentials(
	host string,
	port int,
	login string,
	key string,
	storagePath string,
) map[string]interface{} {
	return map[string]interface{}{
		"host":         host,
		"port":         port,
		"login":        login,
		"key":          key,
		"storage_path": storagePath,
	}
}

// BackupNodeHetznerStorageBoxCredentials creates Hetzner Storage Box specific
// connection credentials.
func BackupNodeHetznerStorageBoxCredentials(
	host string,
	login string,
	key string,
) map[string]interface{} {
	return map[string]interface{}{
		"host":  host,
		"login": login,
		"key":   key,
	}
}

type backupNodeResponse struct {
	Data BackupNode `json:"data"`
}

// Create creates new backup node.
func (s *BackupNodesService) Create(ctx context.Context, data BackupNodeRequest) (BackupNode, error) {
	var resp backupNodeResponse
	return resp.Data, s.client.create(ctx, "backup_nodes", data, &resp)
}

// Update updates specified backup node.
func (s *BackupNodesService) Update(ctx context.Context, id int, data BackupNodeRequest) (BackupNode, error) {
	var resp backupNodeResponse
	return resp.Data, s.client.update(ctx, fmt.Sprintf("backup_nodes/%d", id), data, &resp)
}

// Delete deletes specified backup node.
func (s *BackupNodesService) Delete(ctx context.Context, id int) error {
	return s.client.syncDelete(ctx, fmt.Sprintf("backup_nodes/%d", id))
}
