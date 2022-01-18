package solus

import (
	"context"
	"fmt"
	"net/http"
)

// ServersService handles all available methods with servers.
type ServersService service

// Server represents a server.
type Server struct {
	ID                    int                  `json:"id"`
	Name                  string               `json:"name"`
	Description           string               `json:"description"`
	VirtualizationType    VirtualizationType   `json:"virtualization_type"`
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

// ServerStatus represents available server statuses.
type ServerStatus string

//goland:noinspection GoUnusedConst
const (
	// ServerStatusNotExists indicates server didn't exists.
	ServerStatusNotExists ServerStatus = "not exists"

	// ServerStatusProcessing indicates some action is performed right now on a
	// server.
	ServerStatusProcessing ServerStatus = "processing"

	// ServerStatusStarted indicates server is started.
	ServerStatusStarted ServerStatus = "started"

	// ServerStatusStopped indicates server is stopped.
	ServerStatusStopped ServerStatus = "stopped"

	// ServerStatusPaused indicates server is paused.
	ServerStatusPaused ServerStatus = "paused"

	// ServerStatusUnavailable indicates server is unavailable.
	// This may happen due to some problem with a compute resource or a server.
	ServerStatusUnavailable ServerStatus = "unavailable"
)

// BootMode represents available server boot modes.
type BootMode string

//goland:noinspection GoUnusedConst
const (
	// BootModeDisk indicates booting from original disk.
	BootModeDisk BootMode = "disk"

	// BootModeRescue indicates booting from rescue ISO image.
	BootModeRescue BootMode = "rescue"
)

// ServerSpecifications represent server specification.
type ServerSpecifications struct {
	Disk int `json:"disk"`
	RAM  int `json:"ram"`
	VCPU int `json:"vcpu"`
}

// ServerUsage represent server usage.
type ServerUsage struct {
	CPU float64 `json:"cpu"`
}

// ServerUpdateRequest represents available properties for updating an existing
// server.
type ServerUpdateRequest struct {
	Name           string                `json:"name,omitempty"`
	BootMode       BootMode              `json:"boot_mode,omitempty"`
	Description    string                `json:"description,omitempty"`
	UserData       string                `json:"user_data,omitempty"`
	FQDNs          []string              `json:"fqdns,omitempty"`
	BackupSettings *ServerBackupSettings `json:"backup_settings,omitempty"`
}

// ServerBackupSettings represents server backup settings.
type ServerBackupSettings struct {
	Enabled  bool                         `json:"enabled,omitempty"`
	Schedule ServerBackupSettingsSchedule `json:"schedule,omitempty"`
	Limit    UnitPlanLimit                `json:"limit,omitempty"`
}

// ServerBackupSettingsSchedule represents server backup settings schedule.
type ServerBackupSettingsSchedule struct {
	Type ServerBackupSettingsScheduleType `json:"type"`
	Time ServerBackupSettingsScheduleTime `json:"time"`
	Days []int                            `json:"days,omitempty"`
}

// ServerBackupSettingsScheduleType represents available server backup scheduling
// types.
type ServerBackupSettingsScheduleType string

//goland:noinspection GoUnusedConst
const (
	// ServerBackupSettingsScheduleTypeMonthly indicates backing up every month.
	ServerBackupSettingsScheduleTypeMonthly ServerBackupSettingsScheduleType = "monthly"

	// ServerBackupSettingsScheduleTypeWeekly indicates backing up every week.
	ServerBackupSettingsScheduleTypeWeekly ServerBackupSettingsScheduleType = "weekly"

	// ServerBackupSettingsScheduleTypeDaily indicates backing up every day.
	ServerBackupSettingsScheduleTypeDaily ServerBackupSettingsScheduleType = "daily"
)

// ServerBackupSettingsScheduleTime represents backup settings schedule time.
type ServerBackupSettingsScheduleTime struct {
	Hour    int `json:"hour"`
	Minutes int `json:"minutes"`
}

// ServersResponse represents paginated list of servers.
// This cursor can be used for iterating over all available server.
type ServersResponse struct {
	paginatedResponse

	Data []Server `json:"data"`
}

type serverResponse struct {
	Data Server `json:"data"`
}

// List lists servers.
func (s *ServersService) List(ctx context.Context, filter *FilterServers) (ServersResponse, error) {
	resp := ServersResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "servers", &resp, withFilter(filter.data))
}

// Get gets specified server.
func (s *ServersService) Get(ctx context.Context, id int) (Server, error) {
	var resp serverResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("servers/%d", id), &resp)
}

// Patch patches specified server.
func (s *ServersService) Patch(ctx context.Context, id int, data ServerUpdateRequest) (Server, error) {
	var resp serverResponse
	return resp.Data, s.client.patch(ctx, fmt.Sprintf("servers/%d", id), data, &resp)
}

// Start starts specified server.
func (s *ServersService) Start(ctx context.Context, id int) (Task, error) {
	return s.client.asyncPost(ctx, fmt.Sprintf("servers/%d/start", id))
}

// Stop stops specified server.
func (s *ServersService) Stop(ctx context.Context, id int) (Task, error) {
	return s.client.asyncPost(ctx, fmt.Sprintf("servers/%d/stop", id))
}

// Restart restarts specified server.
func (s *ServersService) Restart(ctx context.Context, id int) (Task, error) {
	return s.client.asyncPost(ctx, fmt.Sprintf("servers/%d/restart", id))
}

// Backup backing up specified server.
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

// Resize resizes specified server.
func (s *ServersService) Resize(ctx context.Context, id int, data ServerResizeRequest) (Task, error) {
	return s.client.asyncPost(ctx, fmt.Sprintf("servers/%d/resize", id), withBody(data))
}

// Delete deletes specified server.
func (s *ServersService) Delete(ctx context.Context, id int) (Task, error) {
	return s.client.asyncDelete(ctx, fmt.Sprintf("servers/%d", id))
}

// SnapshotsCreate creates a snapshot for the specified server.
func (s *ServersService) SnapshotsCreate(ctx context.Context, vmID int, data SnapshotRequest) (Snapshot, error) {
	var resp snapshotResponse
	return resp.Data, s.client.create(ctx, fmt.Sprintf("servers/%d/snapshots", vmID), data, &resp)
}
