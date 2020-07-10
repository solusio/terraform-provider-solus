package solus

import (
	"context"
	"fmt"
	"net/http"
)

type BackupsService service

type BackupType string

const (
	BackupTypeAuto   BackupType = "auto"
	BackupTypeManual BackupType = "manual"
)

type BackupStatus string

const (
	BackupStatusPending    BackupStatus = "pending"
	BackupStatusInProgress BackupStatus = "in_progress"
	BackupStatusCreated    BackupStatus = "created"
	BackupStatusFailed     BackupStatus = "failed"
)

type Backup struct {
	ID                int          `json:"id"`
	Type              BackupType   `json:"type"`
	Status            BackupStatus `json:"status"`
	Size              float32      `json:"size"`
	ComputeResourceVM Server       `json:"compute_resource_vm"`
	BackupNode        BackupNode   `json:"backup_node"`
	Creator           User         `json:"creator"`
	CreatedAt         string       `json:"created_at"`
	BackupProgress    float32      `json:"backup_progress"`
	BackupFailReason  string       `json:"backup_fail_reason"`
	Disk              int          `json:"disk"`
}

func (b Backup) IsFinished() bool {
	return b.Status == BackupStatusCreated ||
		b.Status == BackupStatusFailed
}

type backupResponse struct {
	Data Backup `json:"data"`
}

func (s *BackupsService) Get(ctx context.Context, id int) (Backup, error) {
	var resp backupResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("backups/%d", id), &resp)
}

func (s *BackupsService) Delete(ctx context.Context, id int) error {
	return s.client.delete(ctx, fmt.Sprintf("backups/%d", id))
}

func (s *BackupsService) Restore(ctx context.Context, id int) (Task, error) {
	path := fmt.Sprintf("backups/%d/restore", id)
	body, code, err := s.client.request(ctx, http.MethodPost, path)
	if err != nil {
		return Task{}, err
	}

	if code != http.StatusOK {
		return Task{}, newHTTPError(http.MethodPost, path, code, body)
	}

	var resp taskResponse
	return resp.Data, unmarshal(body, &resp)
}
