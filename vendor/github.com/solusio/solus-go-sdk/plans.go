package solus

import (
	"context"
	"fmt"
)

type PlansService service

type PlanParams struct {
	Disk int `json:"disk"`
	RAM  int `json:"ram"`
	VCPU int `json:"vcpu"`
}

type PlanBackupSettings struct {
	IsIncrementalBackupEnabled bool `json:"is_incremental_backup_enabled"`
	IncrementalBackupsLimit    int  `json:"incremental_backups_limit"`
}

type DiskBandwidthPlanLimit struct {
	IsEnabled bool                       `json:"is_enabled"`
	Limit     int                        `json:"limit"`
	Unit      DiskBandwidthPlanLimitUnit `json:"unit"`
}

type DiskBandwidthPlanLimitUnit string

const (
	DiskBandwidthPlanLimitUnitBps DiskBandwidthPlanLimitUnit = "Bps"
)

func (s *DiskBandwidthPlanLimit) setDefault() {
	if s.Unit == "" {
		s.Unit = DiskBandwidthPlanLimitUnitBps
	}
}

type BandwidthPlanLimit struct {
	IsEnabled bool                   `json:"is_enabled"`
	Limit     int                    `json:"limit"`
	Unit      BandwidthPlanLimitUnit `json:"unit"`
}

type BandwidthPlanLimitUnit string

const (
	BandwidthPlanLimitUnitKbps BandwidthPlanLimitUnit = "Kbps"
	BandwidthPlanLimitUnitMbps BandwidthPlanLimitUnit = "Mbps"
	BandwidthPlanLimitUnitGbps BandwidthPlanLimitUnit = "Gbps"
)

func (s *BandwidthPlanLimit) setDefault() {
	if s.Unit == "" {
		s.Unit = BandwidthPlanLimitUnitKbps
	}
}

type DiskIOPSPlanLimit struct {
	IsEnabled bool                  `json:"is_enabled"`
	Limit     int                   `json:"limit"`
	Unit      DiskIOPSPlanLimitUnit `json:"unit"`
}

type DiskIOPSPlanLimitUnit string

const (
	DiskIOPSPlanLimitUnitIOPS DiskIOPSPlanLimitUnit = "iops"
)

func (s *DiskIOPSPlanLimit) setDefault() {
	if s.Unit == "" {
		s.Unit = DiskIOPSPlanLimitUnitIOPS
	}
}

type TrafficPlanLimit struct {
	IsEnabled bool                 `json:"is_enabled"`
	Limit     int                  `json:"limit"`
	Unit      TrafficPlanLimitUnit `json:"unit"`
}

type TrafficPlanLimitUnit string

const (
	TrafficPlanLimitUnitKB TrafficPlanLimitUnit = "KB"
	TrafficPlanLimitUnitMB TrafficPlanLimitUnit = "MB"
	TrafficPlanLimitUnitGB TrafficPlanLimitUnit = "GB"
	TrafficPlanLimitUnitTB TrafficPlanLimitUnit = "TB"
	TrafficPlanLimitUnitPB TrafficPlanLimitUnit = "PB"
)

func (s *TrafficPlanLimit) setDefault() {
	if s.Unit == "" {
		s.Unit = TrafficPlanLimitUnitKB
	}
}

type UnitPlanLimit struct {
	IsEnabled bool          `json:"is_enabled"`
	Limit     int           `json:"limit"`
	Unit      PlanLimitUnit `json:"unit"`
}

type PlanLimitUnit string

const (
	PlanLimitUnits PlanLimitUnit = "units"
)

func (s *UnitPlanLimit) setDefault() {
	if s.Unit == "" {
		s.Unit = PlanLimitUnits
	}
}

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

type PlanUpdateLimits struct {
	NetworkIncomingBandwidth BandwidthPlanLimit `json:"network_incoming_bandwidth"`
	NetworkOutgoingBandwidth BandwidthPlanLimit `json:"network_outgoing_bandwidth"`
	NetworkIncomingTraffic   TrafficPlanLimit   `json:"network_incoming_traffic"`
	NetworkOutgoingTraffic   TrafficPlanLimit   `json:"network_outgoing_traffic"`
	NetworkTotalTraffic      TrafficPlanLimit   `json:"network_total_traffic"`
	NetworkReduceBandwidth   BandwidthPlanLimit `json:"network_reduce_bandwidth"`
	BackupsNumber            UnitPlanLimit      `json:"backups_number"`
}

type PlanResetLimitPolicy string

const (
	PlanResetLimitPolicyNever           PlanResetLimitPolicy = "never"
	PlanResetLimitPolicyFirstDayOfMonth PlanResetLimitPolicy = "first_day_of_month"
	PlanResetLimitPolicyVMCreatedDay    PlanResetLimitPolicy = "vm_created_day"
)

type PlanNetworkTotalTrafficType string

const (
	PlanNetworkTotalTrafficTypeSeparate PlanNetworkTotalTrafficType = "separate"
	PlanNetworkTotalTrafficTypeTotal    PlanNetworkTotalTrafficType = "total"
)

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

type Plan struct {
	ID                       int                         `json:"id"`
	Name                     string                      `json:"name"`
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
}

type PlanCreateRequest struct {
	Name                     string                      `json:"name"`
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
}

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
}

type PlansResponse struct {
	paginatedResponse

	Data []Plan `json:"data"`
}

type planResponse struct {
	Data Plan `json:"data"`
}

func (s *PlansService) List(ctx context.Context, filter *FilterPlans) (PlansResponse, error) {
	resp := PlansResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "plans", &resp, withFilter(filter.data))
}

func (s *PlansService) Get(ctx context.Context, id int) (Plan, error) {
	var resp planResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("plans/%d", id), &resp)
}

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

func (s *PlansService) Delete(ctx context.Context, id int) error {
	return s.client.delete(ctx, fmt.Sprintf("plans/%d", id))
}
