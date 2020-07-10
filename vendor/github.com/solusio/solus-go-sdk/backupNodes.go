package solus

import (
	"context"
	"fmt"
)

type BackupNodesService service

type BackupNodeType string

const (
	BackupNodeTypeSSHRsync          BackupNodeType = "ssh_rsync"
	BackupNodeTypeHetznerStorageBox BackupNodeType = "hetzner_storage_box"
)

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

type BackupNodeRequest struct {
	Name             string                 `json:"name"`
	Type             BackupNodeType         `json:"type"`
	ComputeResources []int                  `json:"compute_resources,omitempty"`
	Credentials      map[string]interface{} `json:"credentials,omitempty"`
}

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

type backupNodeResponse struct {
	Data BackupNode `json:"data"`
}

func (s *BackupNodesService) Create(ctx context.Context, data BackupNodeRequest) (BackupNode, error) {
	var resp backupNodeResponse
	return resp.Data, s.client.create(ctx, "backup_nodes", data, &resp)
}

func (s *BackupNodesService) Update(ctx context.Context, id int, data BackupNodeRequest) (BackupNode, error) {
	var resp backupNodeResponse
	return resp.Data, s.client.update(ctx, fmt.Sprintf("backup_nodes/%d", id), data, &resp)
}

func (s *BackupNodesService) Delete(ctx context.Context, id int) error {
	return s.client.delete(ctx, fmt.Sprintf("backup_nodes/%d", id))
}
