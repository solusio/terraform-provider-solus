package solus

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type ServersService service

type ServerStatus string

const (
	ServerStatusNotExists   ServerStatus = "not exists"
	ServerStatusProcessing  ServerStatus = "processing"
	ServerStatusStarted     ServerStatus = "started"
	ServerStatusStopped     ServerStatus = "stopped"
	ServerStatusPaused      ServerStatus = "paused"
	ServerStatusUnavailable ServerStatus = "unavailable"
)

type BootMode string

const (
	BootModeDisk   BootMode = "disk"
	BootModeRescue BootMode = "rescue"
)

type Server struct {
	ID                    int                  `json:"id"`
	Name                  string               `json:"name"`
	Description           string               `json:"description"`
	UUID                  string               `json:"uuid"`
	Specifications        ServerSpecifications `json:"specifications"`
	Status                ServerStatus         `json:"status"`
	IPs                   []IPBlockIPAddress   `json:"ips"`
	Location              Location             `json:"location"`
	Plan                  Plan                 `json:"plan"`
	FQDNs                 []string             `json:"fqdns"`
	BootMode              BootMode             `json:"boot_mode"`
	IsSuspended           bool                 `json:"is_suspended"`
	IsProcessing          bool                 `json:"is_processing"`
	User                  User                 `json:"user"`
	Project               Project              `json:"project"`
	Usage                 ServerUsage          `json:"usage"`
	BackupSettings        ServerBackupSettings `json:"backup_settings"`
	NextScheduledBackupAt string               `json:"next_scheduled_backup_at"`
	SSHKeys               []SSHKey             `json:"ssh_keys"`
	CreatedAt             string               `json:"created_at"`
}

type ServerSpecifications struct {
	Disk int `json:"disk"`
	RAM  int `json:"ram"`
	VCPU int `json:"vcpu"`
}

type ServerUsage struct {
	CPU float64 `json:"cpu"`
}

type ServerUpdateRequest struct {
	Name           string                `json:"name,omitempty"`
	BootMode       BootMode              `json:"boot_mode,omitempty"`
	Description    string                `json:"description,omitempty"`
	UserData       string                `json:"user_data,omitempty"`
	FQDNs          []string              `json:"fqdns,omitempty"`
	BackupSettings *ServerBackupSettings `json:"backup_settings,omitempty"`
}

type serverResponse struct {
	Data Server `json:"data"`
}

type ServerBackupSettingsScheduleType string

const (
	ServerBackupSettingsScheduleTypeMonthly ServerBackupSettingsScheduleType = "monthly"
	ServerBackupSettingsScheduleTypeWeekly  ServerBackupSettingsScheduleType = "weekly"
	ServerBackupSettingsScheduleTypeDaily   ServerBackupSettingsScheduleType = "daily"
)

type ServerBackupSettingsScheduleTime struct {
	Hour    int `json:"hour"`
	Minutes int `json:"minutes"`
}

type ServerBackupSettingsSchedule struct {
	Type ServerBackupSettingsScheduleType `json:"type"`
	Time ServerBackupSettingsScheduleTime `json:"time"`
	Days []int                            `json:"days,omitempty"`
}

type ServerBackupSettings struct {
	Enabled  bool                         `json:"enabled,omitempty"`
	Schedule ServerBackupSettingsSchedule `json:"schedule,omitempty"`
	Limit    UnitPlanLimit                `json:"limit,omitempty"`
}

type ServersResponse struct {
	paginatedResponse

	Data []Server `json:"data"`
}

// List return list of server, filter can be nil.
func (s *ServersService) List(ctx context.Context, filter *FilterServers) (ServersResponse, error) {
	resp := ServersResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "servers", &resp, withFilter(filter.data))
}

func (s *ServersService) Get(ctx context.Context, id int) (Server, error) {
	var resp serverResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("servers/%d", id), &resp)
}

func (s *ServersService) Patch(ctx context.Context, id int, data ServerUpdateRequest) (Server, error) {
	var resp serverResponse
	return resp.Data, s.client.patch(ctx, fmt.Sprintf("servers/%d", id), data, &resp)
}

func (s *ServersService) Start(ctx context.Context, id int) (Task, error) {
	return s.action(ctx, id, "start")
}

func (s *ServersService) Stop(ctx context.Context, id int) (Task, error) {
	return s.action(ctx, id, "stop")
}

func (s *ServersService) Restart(ctx context.Context, id int) (Task, error) {
	return s.action(ctx, id, "restart")
}

func (s *ServersService) Backup(ctx context.Context, id int) (Backup, error) {
	path := fmt.Sprintf("servers/%d/backups", id)
	body, code, err := s.client.request(ctx, http.MethodPost, path)
	if err != nil {
		return Backup{}, err
	}

	if code != http.StatusCreated {
		return Backup{}, newHTTPError(http.MethodPost, path, code, body)
	}

	var resp backupResponse
	return resp.Data, unmarshal(body, &resp)
}

type ServerResizeRequest struct {
	PreserveDisk   bool                  `json:"preserve_disk"`
	PlanID         int                   `json:"plan_id"`
	BackupSettings *ServerBackupSettings `json:"backup_settings,omitempty"`
}

func (s *ServersService) Resize(ctx context.Context, id int, data ServerResizeRequest) (Task, error) {
	path := fmt.Sprintf("servers/%d/resize", id)
	body, code, err := s.client.request(ctx, http.MethodPost, path, withBody(data))
	if err != nil {
		return Task{}, err
	}

	if code != http.StatusOK {
		return Task{}, newHTTPError(http.MethodPost, path, code, body)
	}

	var resp taskResponse
	return resp.Data, unmarshal(body, &resp)
}

func (s *ServersService) Delete(ctx context.Context, id int) (Task, error) {
	path := fmt.Sprintf("servers/%d", id)
	body, code, err := s.client.request(ctx, http.MethodDelete, path)
	if err != nil {
		return Task{}, err
	}

	if code != http.StatusOK {
		return Task{}, newHTTPError(http.MethodDelete, path, code, body)
	}

	var resp taskResponse
	if err := unmarshal(body, &resp); err != nil {
		return Task{}, err
	}

	if resp.Data.ID == 0 {
		return Task{}, errors.New("task doesn't have an id")
	}

	return resp.Data, nil
}

func (s *ServersService) action(ctx context.Context, id int, action string) (Task, error) {
	path := fmt.Sprintf("servers/%d/%s", id, action)
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
