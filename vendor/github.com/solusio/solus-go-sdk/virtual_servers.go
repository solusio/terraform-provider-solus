package solus

import (
	"context"
	"fmt"
	"net/http"
)

// VirtualServersService handles all available methods with virtual servers.
type VirtualServersService service

// VirtualServer represents a virtual server.
type VirtualServer struct {
	ID                    int                         `json:"id"`
	Name                  string                      `json:"name"`
	Description           string                      `json:"description"`
	VirtualizationType    VirtualizationType          `json:"virtualization_type"`
	UUID                  string                      `json:"uuid"`
	Specifications        VirtualServerSpecifications `json:"specifications"`
	Status                VirtualServerStatus         `json:"status"`
	IPs                   []IPBlockIPAddress          `json:"ips"`
	Location              Location                    `json:"location"`
	Plan                  Plan                        `json:"plan"`
	FQDNs                 []string                    `json:"fqdns"`
	BootMode              BootMode                    `json:"boot_mode"`
	IsSuspended           bool                        `json:"is_suspended"`
	IsProcessing          bool                        `json:"is_processing"`
	User                  User                        `json:"user"`
	Project               Project                     `json:"project"`
	Usage                 VirtualServerUsage          `json:"usage"`
	BackupSettings        VirtualServerBackupSettings `json:"backup_settings"`
	NextScheduledBackupAt string                      `json:"next_scheduled_backup_at"`
	SSHKeys               []SSHKey                    `json:"ssh_keys"`
	CreatedAt             string                      `json:"created_at"`
}

// VirtualServerStatus represents available virtual server statuses.
type VirtualServerStatus string

//goland:noinspection GoUnusedConst
const (
	// VirtualServerStatusNotExists indicates virtual server didn't exists.
	VirtualServerStatusNotExists VirtualServerStatus = "not exists"

	// VirtualServerStatusProcessing indicates some action is performed right now
	// on a virtual server.
	VirtualServerStatusProcessing VirtualServerStatus = "processing"

	// VirtualServerStatusStarted indicates virtual server is started.
	VirtualServerStatusStarted VirtualServerStatus = "started"

	// VirtualServerStatusStopped indicates virtual server is stopped.
	VirtualServerStatusStopped VirtualServerStatus = "stopped"

	// VirtualServerStatusPaused indicates virtual server is paused.
	VirtualServerStatusPaused VirtualServerStatus = "paused"

	// VirtualServerStatusUnavailable indicates virtual server is unavailable.
	// This may happen due to some problem with a compute resource or a server.
	VirtualServerStatusUnavailable VirtualServerStatus = "unavailable"
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

// VirtualServerSpecifications represent virtual server specification.
type VirtualServerSpecifications struct {
	Disk int `json:"disk"`
	RAM  int `json:"ram"`
	VCPU int `json:"vcpu"`
}

// VirtualServerUsage represent virtual server usage.
type VirtualServerUsage struct {
	CPU float64 `json:"cpu"`
}

type VirtualServerCreateRequest struct {
	Name             string                 `json:"name"`
	BootMode         BootMode               `json:"boot_mode"`
	Description      string                 `json:"description,omitempty"`
	UserData         string                 `json:"user_data,omitempty"`
	FQDNs            []string               `json:"fqdns,omitempty"`
	Password         string                 `json:"password,omitempty"`
	SSHKeys          []int                  `json:"ssh_keys"`
	PlanID           int                    `json:"plan"`
	ProjectID        int                    `json:"project"`
	LocationID       int                    `json:"location"`
	OSImageVersionID int                    `json:"os,omitempty"`
	ApplicationID    int                    `json:"application,omitempty"`
	ApplicationData  map[string]interface{} `json:"applicationData,omitempty"`
}

// VirtualServerUpdateRequest represents available properties for updating an existing
// virtual server.
type VirtualServerUpdateRequest struct {
	Name           string                       `json:"name,omitempty"`
	BootMode       BootMode                     `json:"boot_mode,omitempty"`
	Description    string                       `json:"description,omitempty"`
	UserData       string                       `json:"user_data,omitempty"`
	FQDNs          []string                     `json:"fqdns,omitempty"`
	BackupSettings *VirtualServerBackupSettings `json:"backup_settings,omitempty"`
}

// VirtualServerBackupSettings represents virtual server backup settings.
type VirtualServerBackupSettings struct {
	Enabled  bool                                `json:"enabled,omitempty"`
	Schedule VirtualServerBackupSettingsSchedule `json:"schedule,omitempty"`
	Limit    UnitPlanLimit                       `json:"limit,omitempty"`
}

// VirtualServerBackupSettingsSchedule represents virtual server backup settings
// schedule.
type VirtualServerBackupSettingsSchedule struct {
	Type VirtualServerBackupSettingsScheduleType `json:"type"`
	Time VirtualServerBackupSettingsScheduleTime `json:"time"`
	Days []int                                   `json:"days,omitempty"`
}

// VirtualServerBackupSettingsScheduleType represents available server backup scheduling
// types.
type VirtualServerBackupSettingsScheduleType string

//goland:noinspection GoUnusedConst
const (
	// ServerBackupSettingsScheduleTypeMonthly indicates backing up every month.
	ServerBackupSettingsScheduleTypeMonthly VirtualServerBackupSettingsScheduleType = "monthly"

	// ServerBackupSettingsScheduleTypeWeekly indicates backing up every week.
	ServerBackupSettingsScheduleTypeWeekly VirtualServerBackupSettingsScheduleType = "weekly"

	// ServerBackupSettingsScheduleTypeDaily indicates backing up every day.
	ServerBackupSettingsScheduleTypeDaily VirtualServerBackupSettingsScheduleType = "daily"
)

// VirtualServerBackupSettingsScheduleTime represents backup settings schedule time.
type VirtualServerBackupSettingsScheduleTime struct {
	Hour    int `json:"hour"`
	Minutes int `json:"minutes"`
}

// VirtualServersResponse represents paginated list of servers.
// This cursor can be used for iterating over all available server.
type VirtualServersResponse struct {
	paginatedResponse

	Data []VirtualServer `json:"data"`
}

type virtualServerResponse struct {
	Data VirtualServer `json:"data"`
}

// Create creates virtual server.
func (s *VirtualServersService) Create(ctx context.Context, data VirtualServerCreateRequest) (VirtualServer, error) {
	var resp virtualServerResponse
	s.setDefaultsForCreateRequest(&data)
	return resp.Data, s.client.create(ctx, "servers", data, &resp)
}

func (*VirtualServersService) setDefaultsForCreateRequest(r *VirtualServerCreateRequest) {
	if r.BootMode == "" {
		r.BootMode = BootModeDisk
	}
}

// List lists virtual servers.
func (s *VirtualServersService) List(
	ctx context.Context,
	filter *FilterVirtualServers,
) (VirtualServersResponse, error) {
	resp := VirtualServersResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "servers", &resp, withFilter(filter.data))
}

// Get gets specified virtual server.
func (s *VirtualServersService) Get(ctx context.Context, id int) (VirtualServer, error) {
	var resp virtualServerResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("servers/%d", id), &resp)
}

// Patch patches specified virtual server.
func (s *VirtualServersService) Patch(
	ctx context.Context,
	id int,
	data VirtualServerUpdateRequest,
) (VirtualServer, error) {
	var resp virtualServerResponse
	return resp.Data, s.client.patch(ctx, fmt.Sprintf("servers/%d", id), data, &resp)
}

// Start starts specified virtual server.
func (s *VirtualServersService) Start(ctx context.Context, id int) (Task, error) {
	return s.client.asyncPost(ctx, fmt.Sprintf("servers/%d/start", id))
}

// Stop stops specified virtual server.
func (s *VirtualServersService) Stop(ctx context.Context, id int) (Task, error) {
	return s.client.asyncPost(ctx, fmt.Sprintf("servers/%d/stop", id))
}

// Restart restarts specified virtual server.
func (s *VirtualServersService) Restart(ctx context.Context, id int) (Task, error) {
	return s.client.asyncPost(ctx, fmt.Sprintf("servers/%d/restart", id))
}

// Backup backing up specified virtual server.
func (s *VirtualServersService) Backup(ctx context.Context, id int) (Backup, error) {
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

type ViretualServerResizeRequest struct {
	PreserveDisk   bool                         `json:"preserve_disk"`
	PlanID         int                          `json:"plan_id"`
	BackupSettings *VirtualServerBackupSettings `json:"backup_settings,omitempty"`
}

// Resize resizes specified virtual server.
func (s *VirtualServersService) Resize(ctx context.Context, id int, data ViretualServerResizeRequest) (Task, error) {
	return s.client.asyncPost(ctx, fmt.Sprintf("servers/%d/resize", id), withBody(data))
}

// Delete deletes specified virtual server.
func (s *VirtualServersService) Delete(ctx context.Context, id int) (Task, error) {
	return s.client.asyncDelete(ctx, fmt.Sprintf("servers/%d", id))
}

// SnapshotsCreate creates a snapshot for the specified virtual server.
func (s *VirtualServersService) SnapshotsCreate(ctx context.Context, vmID int, data SnapshotRequest) (Snapshot, error) {
	var resp snapshotResponse
	return resp.Data, s.client.create(ctx, fmt.Sprintf("servers/%d/snapshots", vmID), data, &resp)
}
