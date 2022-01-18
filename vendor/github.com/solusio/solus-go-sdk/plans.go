package solus

import (
	"context"
	"fmt"
)

// PlansService handles all available methods with plans.
type PlansService service

// Plan represents a plan.
type Plan struct {
	ID                       int                         `json:"id"`
	Name                     string                      `json:"name"`
	VirtualizationType       VirtualizationType          `json:"virtualization_type"`
	Params                   PlanParams                  `json:"params"`
	StorageType              string                      `json:"storage_type"`
	ImageFormat              string                      `json:"image_format"`
	IsDefault                bool                        `json:"is_default"`
	IsSnapshotAvailable      bool                        `json:"is_snapshot_available"`
	IsSnapshotsEnabled       bool                        `json:"is_snapshots_enabled"`
	IsBackupAvailable        bool                        `json:"is_backup_available"`
	IsAdditionalIPsAvailable bool                        `json:"is_additional_ips_available"`
	BackupSettings           PlanBackupSettings          `json:"backup_settings"`
	BackupPrice              float64                     `json:"backup_price"`
	IsVisible                bool                        `json:"is_visible"`
	Limits                   PlanLimits                  `json:"limits"`
	TokensPerHour            int                         `json:"tokens_per_hour"`
	TokensPerMonth           int                         `json:"tokens_per_month"`
	IPTokensPerHour          int                         `json:"ip_tokens_per_hour"`
	IPTokensPerMonth         int                         `json:"ip_tokens_per_month"`
	Position                 float64                     `json:"position"`
	Price                    PlanPrice                   `json:"price"`
	ResetLimitPolicy         PlanResetLimitPolicy        `json:"reset_limit_policy"`
	NetworkTotalTrafficType  PlanNetworkTotalTrafficType `json:"network_traffic_limit_type"`
	AvailableLocations       []ShortLocation             `json:"available_locations"`
	AvailableOsImageVersions []ShortOsImageVersion       `json:"available_os_image_versions"`
	AvailableApplications    []ShortOsImageVersion       `json:"available_applications"`
}

// ShortPlan represents only ID and name of plan.
type ShortPlan struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// PlanParams represents a plan's parameters.
type PlanParams struct {
	Disk       int `json:"disk"`
	RAM        int `json:"ram"`
	VCPU       int `json:"vcpu"`
	VCPUUnits  int `json:"vcpu_units"`
	VCPULimit  int `json:"vcpu_limit"`
	IOPriority int `json:"io_priority"`
}

// PlanBackupSettings represents a plan's backup settings.
type PlanBackupSettings struct {
	IsIncrementalBackupEnabled bool `json:"is_incremental_backup_enabled"`
	IncrementalBackupsLimit    int  `json:"incremental_backups_limit"`
}

// DiskBandwidthPlanLimit represents disk bandwidth specific limit.
type DiskBandwidthPlanLimit struct {
	IsEnabled bool                       `json:"is_enabled"`
	Limit     int                        `json:"limit"`
	Unit      DiskBandwidthPlanLimitUnit `json:"unit"`
}

// DiskBandwidthPlanLimitUnit represents available units for disk bandwidth limit.
type DiskBandwidthPlanLimitUnit string

const (
	// DiskBandwidthPlanLimitUnitBps indicates bytes per second unit.
	DiskBandwidthPlanLimitUnitBps DiskBandwidthPlanLimitUnit = "Bps"
)

func (s *DiskBandwidthPlanLimit) setDefault() {
	if s.Unit == "" {
		s.Unit = DiskBandwidthPlanLimitUnitBps
	}
}

// BandwidthPlanLimit represents network bandwidth specific limit.
type BandwidthPlanLimit struct {
	IsEnabled bool                   `json:"is_enabled"`
	Limit     int                    `json:"limit"`
	Unit      BandwidthPlanLimitUnit `json:"unit"`
}

// BandwidthPlanLimitUnit represents available units for network bandwidth limit.
type BandwidthPlanLimitUnit string

const (
	// BandwidthPlanLimitUnitKbps indicates kilobits per second unit.
	BandwidthPlanLimitUnitKbps BandwidthPlanLimitUnit = "Kbps"

	// BandwidthPlanLimitUnitMbps indicates megabits per second unit.
	BandwidthPlanLimitUnitMbps BandwidthPlanLimitUnit = "Mbps"

	// BandwidthPlanLimitUnitGbps indicates gigabits per second unit.
	BandwidthPlanLimitUnitGbps BandwidthPlanLimitUnit = "Gbps"
)

func (s *BandwidthPlanLimit) setDefault() {
	if s.Unit == "" {
		s.Unit = BandwidthPlanLimitUnitKbps
	}
}

// DiskIOPSPlanLimit represents disk IOPS specific limit.
type DiskIOPSPlanLimit struct {
	IsEnabled bool                  `json:"is_enabled"`
	Limit     int                   `json:"limit"`
	Unit      DiskIOPSPlanLimitUnit `json:"unit"`
}

// DiskIOPSPlanLimitUnit represents available units for disk IOPS limit.
type DiskIOPSPlanLimitUnit string

const (
	// DiskIOPSPlanLimitUnitIOPS indicates input/output operations per second unit.
	DiskIOPSPlanLimitUnitIOPS DiskIOPSPlanLimitUnit = "iops"
)

func (s *DiskIOPSPlanLimit) setDefault() {
	if s.Unit == "" {
		s.Unit = DiskIOPSPlanLimitUnitIOPS
	}
}

// TrafficPlanLimit represents network traffic specific limit.
type TrafficPlanLimit struct {
	IsEnabled bool                 `json:"is_enabled"`
	Limit     int                  `json:"limit"`
	Unit      TrafficPlanLimitUnit `json:"unit"`
}

// TrafficPlanLimitUnit represents available units for network traffic specific
// limit.
type TrafficPlanLimitUnit string

//goland:noinspection GoUnusedConst
const (
	// TrafficPlanLimitUnitKB indicates kilobytes unit.
	TrafficPlanLimitUnitKB TrafficPlanLimitUnit = "KB"

	// TrafficPlanLimitUnitMB indicates megabytes unit.
	TrafficPlanLimitUnitMB TrafficPlanLimitUnit = "MB"

	// TrafficPlanLimitUnitGB indicates gigabytes unit.
	TrafficPlanLimitUnitGB TrafficPlanLimitUnit = "GB"

	// TrafficPlanLimitUnitTB indicates terabytes unit.
	TrafficPlanLimitUnitTB TrafficPlanLimitUnit = "TB"

	// TrafficPlanLimitUnitPB indicates petabytes unit.
	TrafficPlanLimitUnitPB TrafficPlanLimitUnit = "PB"
)

func (s *TrafficPlanLimit) setDefault() {
	if s.Unit == "" {
		s.Unit = TrafficPlanLimitUnitKB
	}
}

// UnitPlanLimit represents generic units limit.
type UnitPlanLimit struct {
	IsEnabled bool          `json:"is_enabled"`
	Limit     int           `json:"limit"`
	Unit      PlanLimitUnit `json:"unit"`
}

// PlanLimitUnit represents available generic units.
type PlanLimitUnit string

const (
	// PlanLimitUnits indicates generic unit.
	PlanLimitUnits PlanLimitUnit = "units"
)

func (s *UnitPlanLimit) setDefault() {
	if s.Unit == "" {
		s.Unit = PlanLimitUnits
	}
}

// PlanLimits represents all available limits for creating a new plan.
type PlanLimits struct {
	DiskBandwidth            DiskBandwidthPlanLimit `json:"disk_bandwidth"`
	DiskIOPS                 DiskIOPSPlanLimit      `json:"disk_iops"`
	NetworkIncomingBandwidth BandwidthPlanLimit     `json:"network_incoming_bandwidth"`
	NetworkOutgoingBandwidth BandwidthPlanLimit     `json:"network_outgoing_bandwidth"`
	NetworkIncomingTraffic   TrafficPlanLimit       `json:"network_incoming_traffic"`
	NetworkOutgoingTraffic   TrafficPlanLimit       `json:"network_outgoing_traffic"`
	NetworkTotalTraffic      TrafficPlanLimit       `json:"network_total_traffic"`
	NetworkReduceBandwidth   BandwidthPlanLimit     `json:"network_reduce_bandwidth"`
	BackupsNumber            UnitPlanLimit          `json:"backups_number"`
}

// PlanUpdateLimits represents all available limits for updating a plan.
type PlanUpdateLimits struct {
	NetworkIncomingBandwidth BandwidthPlanLimit `json:"network_incoming_bandwidth"`
	NetworkOutgoingBandwidth BandwidthPlanLimit `json:"network_outgoing_bandwidth"`
	NetworkIncomingTraffic   TrafficPlanLimit   `json:"network_incoming_traffic"`
	NetworkOutgoingTraffic   TrafficPlanLimit   `json:"network_outgoing_traffic"`
	NetworkTotalTraffic      TrafficPlanLimit   `json:"network_total_traffic"`
	NetworkReduceBandwidth   BandwidthPlanLimit `json:"network_reduce_bandwidth"`
	BackupsNumber            UnitPlanLimit      `json:"backups_number"`
}

// PlanResetLimitPolicy represent available a limit reset policies.
type PlanResetLimitPolicy string

const (
	// PlanResetLimitPolicyNever indicates we shouldn't reset a limit at all.
	PlanResetLimitPolicyNever PlanResetLimitPolicy = "never"

	// PlanResetLimitPolicyFirstDayOfMonth indicates we should reset a limit at the
	// first day of month.
	PlanResetLimitPolicyFirstDayOfMonth PlanResetLimitPolicy = "first_day_of_month"

	// PlanResetLimitPolicyVMCreatedDay indicates we should reset a limit at the VM's
	// created day.
	PlanResetLimitPolicyVMCreatedDay PlanResetLimitPolicy = "vm_created_day"
)

// PlanNetworkTotalTrafficType represent available total traffic types.
type PlanNetworkTotalTrafficType string

//goland:noinspection GoUnusedConst
const (
	// PlanNetworkTotalTrafficTypeSeparate indicates we should count incoming and
	// outgoing traffic separately.
	PlanNetworkTotalTrafficTypeSeparate PlanNetworkTotalTrafficType = "separate"

	// PlanNetworkTotalTrafficTypeTotal indicates we should count incoming and
	// outgoing traffic together.
	PlanNetworkTotalTrafficTypeTotal PlanNetworkTotalTrafficType = "total"
)

// PlanPrice represents plan price information.
type PlanPrice struct {
	AdditionalIPsPerHour     string        `json:"additional_ips_per_hour"`
	AdditionalIPsPerMonth    string        `json:"additional_ips_per_month"`
	PerHour                  string        `json:"per_hour"`
	PerMonth                 string        `json:"per_month"`
	CurrencyCode             string        `json:"currency_code"`
	TaxesInclusive           bool          `json:"taxes_inclusive"`
	Taxes                    []interface{} `json:"taxes"`
	TotalPriceWithoutBackups string        `json:"total_price_without_backups"`
	TotalPrice               string        `json:"total_price"`
	BackupPrice              string        `json:"backup_price"`
}

// PlanCreateRequest represents available properties for creating a new plan.
type PlanCreateRequest struct {
	Name                     string                      `json:"name"`
	VirtualizationType       VirtualizationType          `json:"virtualization_type"`
	Params                   PlanParams                  `json:"params"`
	StorageType              StorageTypeName             `json:"storage_type"`
	ImageFormat              ImageFormat                 `json:"image_format"`
	Limits                   PlanLimits                  `json:"limits"`
	TokensPerHour            int                         `json:"tokens_per_hour"`
	TokensPerMonth           int                         `json:"tokens_per_month"`
	IPTokensPerHour          int                         `json:"ip_tokens_per_hour"`
	IPTokensPerMonth         int                         `json:"ip_tokens_per_month"`
	IsVisible                bool                        `json:"is_visible"`
	IsDefault                bool                        `json:"is_default"`
	IsSnapshotsEnabled       bool                        `json:"is_snapshots_enabled"`
	IsBackupAvailable        bool                        `json:"is_backup_available"`
	IsAdditionalIPsAvailable bool                        `json:"is_additional_ips_available"`
	BackupSettings           PlanBackupSettings          `json:"backup_settings"`
	BackupPrice              float64                     `json:"backup_price"`
	ResetLimitPolicy         PlanResetLimitPolicy        `json:"reset_limit_policy"`
	NetworkTotalTrafficType  PlanNetworkTotalTrafficType `json:"network_traffic_limit_type"`
	AvailableLocations       []int                       `json:"available_locations,omitempty"`
	AvailableOsImageVersions []int                       `json:"available_os_image_versions,omitempty"`
	AvailableApplications    []int                       `json:"available_applications,omitempty"`
}

// PlanUpdateRequest represents available properties for updating an existing plan.
type PlanUpdateRequest struct {
	Name                     string                      `json:"name"`
	Limits                   PlanUpdateLimits            `json:"limits"`
	TokensPerHour            int                         `json:"tokens_per_hour"`
	TokensPerMonth           int                         `json:"tokens_per_month"`
	IPTokensPerHour          int                         `json:"ip_tokens_per_hour"`
	IPTokensPerMonth         int                         `json:"ip_tokens_per_month"`
	IsVisible                bool                        `json:"is_visible"`
	IsDefault                bool                        `json:"is_default"`
	IsSnapshotsEnabled       bool                        `json:"is_snapshots_enabled"`
	IsBackupAvailable        bool                        `json:"is_backup_available"`
	IsAdditionalIPsAvailable bool                        `json:"is_additional_ips_available"`
	BackupSettings           PlanBackupSettings          `json:"backup_settings"`
	BackupPrice              float64                     `json:"backup_price"`
	ResetLimitPolicy         PlanResetLimitPolicy        `json:"reset_limit_policy"`
	NetworkTotalTrafficType  PlanNetworkTotalTrafficType `json:"network_traffic_limit_type"`
	AvailableLocations       []int                       `json:"available_locations,omitempty"`
	AvailableOsImageVersions []int                       `json:"available_os_image_versions,omitempty"`
	AvailableApplications    []int                       `json:"available_applications,omitempty"`
}

// PlansResponse represents paginated list of plans.
// This cursor can be used for iterating over all available plans.
type PlansResponse struct {
	paginatedResponse

	Data []Plan `json:"data"`
}

type planResponse struct {
	Data Plan `json:"data"`
}

// List lists plans.
func (s *PlansService) List(ctx context.Context, filter *FilterPlans) (PlansResponse, error) {
	resp := PlansResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "plans", &resp, withFilter(filter.data))
}

// Get gets specified plan.
func (s *PlansService) Get(ctx context.Context, id int) (Plan, error) {
	var resp planResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("plans/%d", id), &resp)
}

// Create creates new plan.
func (s *PlansService) Create(ctx context.Context, data PlanCreateRequest) (Plan, error) {
	s.setCreateRequestDefaults(&data)
	var resp planResponse
	return resp.Data, s.client.create(ctx, "plans", data, &resp)
}

func (*PlansService) setCreateRequestDefaults(data *PlanCreateRequest) {
	if data.ResetLimitPolicy == "" {
		data.ResetLimitPolicy = PlanResetLimitPolicyNever
	}

	data.Limits.DiskBandwidth.setDefault()
	data.Limits.DiskIOPS.setDefault()
	data.Limits.NetworkIncomingBandwidth.setDefault()
	data.Limits.NetworkOutgoingBandwidth.setDefault()
	data.Limits.NetworkReduceBandwidth.setDefault()
	data.Limits.NetworkIncomingTraffic.setDefault()
	data.Limits.NetworkOutgoingTraffic.setDefault()
	data.Limits.NetworkTotalTraffic.setDefault()
	data.Limits.BackupsNumber.setDefault()
}

// Update updates specified plan.
func (s *PlansService) Update(ctx context.Context, id int, data PlanUpdateRequest) (Plan, error) {
	s.setUpdateRequestDefaults(&data)
	var resp planResponse
	return resp.Data, s.client.update(ctx, fmt.Sprintf("plans/%d", id), data, &resp)
}

func (*PlansService) setUpdateRequestDefaults(data *PlanUpdateRequest) {
	if data.ResetLimitPolicy == "" {
		data.ResetLimitPolicy = PlanResetLimitPolicyNever
	}

	data.Limits.NetworkIncomingBandwidth.setDefault()
	data.Limits.NetworkOutgoingBandwidth.setDefault()
	data.Limits.NetworkReduceBandwidth.setDefault()
	data.Limits.NetworkIncomingTraffic.setDefault()
	data.Limits.NetworkOutgoingTraffic.setDefault()
	data.Limits.NetworkTotalTraffic.setDefault()
	data.Limits.BackupsNumber.setDefault()
}

// Delete deletes specified plan.
func (s *PlansService) Delete(ctx context.Context, id int) error {
	return s.client.syncDelete(ctx, fmt.Sprintf("plans/%d", id))
}
