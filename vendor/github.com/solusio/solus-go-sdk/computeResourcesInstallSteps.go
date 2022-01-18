package solus

import (
	"context"
	"fmt"
)

// ComputeResourceInstallStep represents a compute resource install step.
type ComputeResourceInstallStep struct {
	ID                int                              `json:"id"`
	ComputeResourceID int                              `json:"compute_resource_id"`
	Title             string                           `json:"title"`
	Status            ComputeResourceInstallStepStatus `json:"status"`
	StatusText        string                           `json:"status_text"`
	Progress          float32                          `json:"progress"`
}

// ComputeResourceInstallStepStatus represents available compute resource's install
// step statuses.
type ComputeResourceInstallStepStatus string

const (
	// ComputeResourceInstallStepStatusRunning indicates install step is running.
	ComputeResourceInstallStepStatusRunning ComputeResourceInstallStepStatus = "running"

	// ComputeResourceInstallStepStatusDone indicates install step is successfully
	// done.
	ComputeResourceInstallStepStatusDone ComputeResourceInstallStepStatus = "done"

	// ComputeResourceInstallStepStatusError indicates install step is failed.
	ComputeResourceInstallStepStatusError ComputeResourceInstallStepStatus = "error"
)

// InstallSteps lists specified compute resource's install steps.
func (s *ComputeResourcesService) InstallSteps(ctx context.Context, id int) ([]ComputeResourceInstallStep, error) {
	var resp struct {
		Data []ComputeResourceInstallStep `json:"data"`
	}
	return resp.Data, s.client.get(ctx, fmt.Sprintf("compute_resources/%d/install_steps", id), &resp)
}
