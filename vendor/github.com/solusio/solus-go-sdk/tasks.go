package solus

import (
	"context"
	"fmt"
)

type TasksService service

type TaskStatus string

const (
	// status
	TaskStatusPending        TaskStatus = "pending"
	TaskStatusQueued         TaskStatus = "queued"
	TaskStatusRunning        TaskStatus = "running"
	TaskStatusDone           TaskStatus = "done"
	TaskStatusDoneWithErrors TaskStatus = "done_with_errors"
	TaskStatusFailed         TaskStatus = "failed"
	TaskStatusCanceled       TaskStatus = "canceled"
)

type TaskAction string

const (
	// actions
	ServerActionCreate              TaskAction = "vm-create"
	ServerActionReinstall           TaskAction = "vm-reinstall"
	ServerActionDelete              TaskAction = "vm-delete"
	ServerActionUpdate              TaskAction = "vm-update"
	ServerActionPasswordChange      TaskAction = "vm-password-change"
	ServerActionStart               TaskAction = "vm-start"
	ServerActionStop                TaskAction = "vm-stop"
	ServerActionRestart             TaskAction = "vm-restart"
	ServerActionSuspend             TaskAction = "vm-suspend"
	ServerActionResume              TaskAction = "vm-resume"
	ComputeResourceConfigureNetwork TaskAction = "configure network"
)

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

func (t Task) IsFinished() bool {
	return t.Status != TaskStatusPending &&
		t.Status != TaskStatusQueued &&
		t.Status != TaskStatusRunning
}

type TasksResponse struct {
	paginatedResponse

	Data []Task `json:"data"`
}

type Date struct {
	Date         string `json:"date"`
	TimezoneType int    `json:"timezone_type"`
	Timezone     string `json:"timezone"`
}

type taskResponse struct {
	Data Task `json:"data"`
}

// Tasks return list of Task, filter can be nil
func (s *TasksService) List(ctx context.Context, filter *FilterTasks) (TasksResponse, error) {
	resp := TasksResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "tasks", &resp, withFilter(filter.data))
}

func (s *TasksService) Get(ctx context.Context, id int) (Task, error) {
	var resp taskResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("tasks/%d", id), &resp)
}
