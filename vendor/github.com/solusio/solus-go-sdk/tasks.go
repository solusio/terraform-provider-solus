package solus

import (
	"context"
	"fmt"
)

// TasksService handles all available methods with tasks.
type TasksService service

// TaskStatus represents available task's statuses.
type TaskStatus string

const (
	// TaskStatusPending indicates task is just created but not dispatched.
	TaskStatusPending TaskStatus = "pending"

	// TaskStatusQueued indicates task is already dispatched but didn't started.
	TaskStatusQueued TaskStatus = "queued"

	// TaskStatusRunning indicates task is running.
	TaskStatusRunning TaskStatus = "running"

	// TaskStatusDone indicates task successfully finished.
	TaskStatusDone TaskStatus = "done"

	// TaskStatusDoneWithErrors indicates task group finished but at least one of
	// children task is failed.
	TaskStatusDoneWithErrors TaskStatus = "done_with_errors"

	// TaskStatusFailed indicates task is failed.
	TaskStatusFailed TaskStatus = "failed"

	// TaskStatusCanceled indicates task is canceled.
	TaskStatusCanceled TaskStatus = "canceled"
)

// TaskAction represents available task actions.
type TaskAction string

const (
	// TaskActionServerCreate indicates server create task.
	TaskActionServerCreate TaskAction = "vm-create"

	// TaskActionServerReinstall indicates server reinstall task.
	TaskActionServerReinstall TaskAction = "vm-reinstall"

	// TaskActionServerDelete indicates server delete task.
	TaskActionServerDelete TaskAction = "vm-delete"

	// TaskActionServerUpdate indicates server update task.
	TaskActionServerUpdate TaskAction = "vm-update"

	// TaskActionServerPasswordChange indicates server password change task.
	TaskActionServerPasswordChange TaskAction = "vm-password-change"

	// TaskActionServerStart indicates server create task.
	TaskActionServerStart TaskAction = "vm-start"

	// TaskActionServerStop indicates server stop task.
	TaskActionServerStop TaskAction = "vm-stop"

	// TaskActionServerRestart indicates server restart task.
	TaskActionServerRestart TaskAction = "vm-restart"

	// TaskActionServerSuspend indicates server suspend task.
	TaskActionServerSuspend TaskAction = "vm-suspend"

	// TaskActionServerResume indicates server resume task.
	TaskActionServerResume TaskAction = "vm-resume"

	// TaskActionServerResize indicates server resize task.
	TaskActionServerResize TaskAction = "vm-resize"

	// TaskActionServersMigrate indicates servers migrate task.
	// This is the main task which is a group of TaskActionServerMigrate tasks.
	TaskActionServersMigrate TaskAction = "vms-migrate"

	// TaskActionServerMigrate indicates server migrate task.
	TaskActionServerMigrate TaskAction = "vm-migrate"

	// TaskActionServerUpdateNetwork indicates server update network task.
	TaskActionServerUpdateNetwork TaskAction = "vm-update-network"

	// TaskActionServerUpdateLimits indicates server update limits task.
	TaskActionServerUpdateLimits TaskAction = "vm-update-limits"

	// TaskActionServersUpdateLimits indicates servers update limits task.
	// Task for batch limits updates.
	TaskActionServersUpdateLimits TaskAction = "vms-update-limits"

	// TaskActionDNSRecordRegister indicates DNS record register task.
	TaskActionDNSRecordRegister TaskAction = "dns-record-register"

	// TaskActionDNSRecordsUnregister indicates DNS record unregister task.
	TaskActionDNSRecordsUnregister TaskAction = "dns-records-unregister"

	// TaskActionDNSRecordsUpdate indicates DNS record update task.
	TaskActionDNSRecordsUpdate TaskAction = "dns-record-update"

	// TaskActionReverseDNSRecordRegister indicates reverse DNS record register
	// task.
	TaskActionReverseDNSRecordRegister TaskAction = "reverse-dns-record-register"

	// TaskActionSnapshotCreate indicates snapshot create task.
	TaskActionSnapshotCreate TaskAction = "snapshot-create"

	// TaskActionSnapshotDelete indicates snapshot delete task.
	TaskActionSnapshotDelete TaskAction = "snapshot-delete"

	// TaskActionSnapshotRevert indicates snapshot revert task.
	TaskActionSnapshotRevert TaskAction = "snapshot-revert"

	// TaskActionPrepareInstallerForUpdate indicates prepare installer before update
	// task.
	TaskActionPrepareInstallerForUpdate TaskAction = "prepare installer for version update"

	// TaskActionRunVersionUpdate indicates application update task.
	TaskActionRunVersionUpdate TaskAction = "run version update"

	// TaskActionBackupCreate indicates backup create task.
	TaskActionBackupCreate TaskAction = "backup-create"

	// TaskActionBackupRestore indicates restore from backup task.
	TaskActionBackupRestore TaskAction = "backup-restore"

	// TaskActionBackupDelete indicates backup delete task.
	TaskActionBackupDelete TaskAction = "backup-delete"

	// TaskActionBackupRotate indicates backups rotate task.
	TaskActionBackupRotate TaskAction = "backup-rotate"

	// TaskActionPurgeComputeResourceVM indicates purging server backups task.
	TaskActionPurgeComputeResourceVM TaskAction = "backup-purge-compute-resource-vm"

	// TaskActionConfigureNetwork indicates compute resource network configuring
	// task.
	TaskActionConfigureNetwork TaskAction = "configure network"

	// TaskActionUpdateNetworkRules indicates compute resource network rules update
	// task.
	TaskActionUpdateNetworkRules TaskAction = "update network rules"

	// TaskActionUpgradeComputeResource indicates compute resource upgrade task.
	TaskActionUpgradeComputeResource TaskAction = "upgrade compute resource"

	// TaskActionClearImageCache indicates compute resource OS image cache clear
	// task.
	TaskActionClearImageCache TaskAction = "clear image cache"

	// TaskActionChangeHostname indicates server hostname change task.
	TaskActionChangeHostname TaskAction = "change hostname"
)

// Task represents a task.
type Task struct {
	ID                int        `json:"id"`
	ComputeResourceID int        `json:"compute_resource_id"`
	Queue             string     `json:"queue"`
	Action            TaskAction `json:"action"`
	Status            TaskStatus `json:"status"`
	Output            string     `json:"output"`
	Progress          int        `json:"progress"`
	Duration          int        `json:"duration"`
}

// IsFinished returns true if the task is finished, successfully or not.
func (t Task) IsFinished() bool {
	return t.Status != TaskStatusPending &&
		t.Status != TaskStatusQueued &&
		t.Status != TaskStatusRunning
}

// TasksResponse represents paginated list of tasks.
// This cursor can be used for iterating over all available users.
type TasksResponse struct {
	paginatedResponse

	Data []Task `json:"data"`
}

type taskResponse struct {
	Data Task `json:"data"`
}

// List lists tasks.
func (s *TasksService) List(ctx context.Context, filter *FilterTasks) (TasksResponse, error) {
	resp := TasksResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "tasks", &resp, withFilter(filter.data))
}

// Get gets specified task.
func (s *TasksService) Get(ctx context.Context, id int) (Task, error) {
	var resp taskResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("tasks/%d", id), &resp)
}
