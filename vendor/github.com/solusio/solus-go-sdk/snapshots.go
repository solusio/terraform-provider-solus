// Copyright 1999-2021. Plesk International GmbH. All rights reserved.

package solus

import (
	"context"
	"fmt"
)

// StorageService handles all available methods with storages.
type SnapshotsService service

// Snapshot represents a snapshot.
type Snapshot struct {
	ID   int    `json:"id"`
	Name string `json:"name"`

	// Size a size of snapshot in Gb.
	Size      float64        `json:"size"`
	Status    SnapshotStatus `json:"status"`
	CreatedAt string         `json:"created_at"`
}

// SnapshotStatus represents available snapshot statuses.
type SnapshotStatus string

const (
	// SnapshotStatusAvailable indicates snapshot is available for reverting.
	SnapshotStatusAvailable SnapshotStatus = "available"

	// SnapshotStatusProcessing indicates snapshot still creating.
	SnapshotStatusProcessing SnapshotStatus = "processing"

	// SnapshotStatusFailed indicates snapshot is failed.
	SnapshotStatusFailed SnapshotStatus = "failed"
)

// SnapshotRequest represents available properties for creating a new snapshot.
type SnapshotRequest struct {
	Name string `json:"name"`
}

type snapshotResponse struct {
	Data Snapshot `json:"data"`
}

// Get gets specified snapshot.
func (s *SnapshotsService) Get(ctx context.Context, id int) (Snapshot, error) {
	var resp snapshotResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("snapshots/%d", id), &resp)
}

// Revert reverts VM from specified snapshot.
func (s *SnapshotsService) Revert(ctx context.Context, id int) (Task, error) {
	return s.client.asyncPost(ctx, fmt.Sprintf("snapshots/%d/revert", id))
}

// Delete deletes specified snapshot.
func (s *SnapshotsService) Delete(ctx context.Context, id int) (Task, error) {
	return s.client.asyncDelete(ctx, fmt.Sprintf("snapshots/%d", id))
}
