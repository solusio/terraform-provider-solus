package solus

import (
	"context"
	"fmt"
)

// BackupsService handles all available methods with backups.
type BackupsService service

// Backup represent a backup.
type Backup struct {
	ID                int                  `json:"id"`
	Type              BackupType           `json:"type"`
	CreationMethod    BackupCreationMethod `json:"creation_method"`
	Status            BackupStatus         `json:"status"`
	Size              float32              `json:"size"`
	ComputeResourceVM VirtualServer        `json:"compute_resource_vm"`
	BackupNode        BackupNode           `json:"backup_node"`
	Creator           User                 `json:"creator"`
	CreatedAt         string               `json:"created_at"`
	BackupProgress    float32              `json:"backup_progress"`
	BackupFailReason  string               `json:"backup_fail_reason"`
	Disk              int                  `json:"disk"`
}

// BackupType a backup type.
type BackupType string

const (
	// BackupTypeFull a full server backup.
	// Include all disk's data.
	BackupTypeFull BackupType = "full"

	// BackupTypeIncremental an incremental server backup.
	// Include only changed disk's data since last backup.
	BackupTypeIncremental BackupType = "incremental"
)

// BackupCreationMethod represents a way how a backup was created.
type BackupCreationMethod string

const (
	// BackupCreationMethodAuto backup was created automatically by scheduler.
	BackupCreationMethodAuto BackupCreationMethod = "auto"

	// BackupCreationMethodManual backup was created manually by user request.
	BackupCreationMethodManual BackupCreationMethod = "manual"
)

// BackupStatus represents available backup's statuses.
type BackupStatus string

const (
	// BackupStatusPending indicates backup is still not processing and waits until
	// it will be dispatched.
	BackupStatusPending BackupStatus = "pending"

	// BackupStatusInProgress indicates backup is processing right now.
	BackupStatusInProgress BackupStatus = "in_progress"

	// BackupStatusCreated indicates backup was successfully created.
	BackupStatusCreated BackupStatus = "created"

	// BackupStatusFailed indicates backup wasn't created due to some reason.
	BackupStatusFailed BackupStatus = "failed"
)

// IsFinished returns true if the backup is finished, successfully or not.
func (b Backup) IsFinished() bool {
	return b.Status == BackupStatusCreated ||
		b.Status == BackupStatusFailed
}

type backupResponse struct {
	Data Backup `json:"data"`
}

// Get gets specified backup.
func (s *BackupsService) Get(ctx context.Context, id int) (Backup, error) {
	var resp backupResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("backups/%d", id), &resp)
}

// Delete deletes specified backup.
func (s *BackupsService) Delete(ctx context.Context, id int) error {
	return s.client.syncDelete(ctx, fmt.Sprintf("backups/%d", id))
}

// Restore restores a related server from a specific backup.
func (s *BackupsService) Restore(ctx context.Context, id int) (Task, error) {
	return s.client.asyncPost(ctx, fmt.Sprintf("backups/%d/restore", id))
}
