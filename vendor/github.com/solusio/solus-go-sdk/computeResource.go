package solus

import (
	"context"
	"fmt"
	"net/http"
)

type ComputeResourcesService service

type ComputerResourceStatus string

const (
	ComputeResourceStatusActive           ComputerResourceStatus = "active"
	ComputeResourceStatusCommissioning    ComputerResourceStatus = "commissioning"
	ComputeResourceStatusConfigureNetwork ComputerResourceStatus = "configure_network"
	ComputeResourceStatusFailed           ComputerResourceStatus = "failed"
	ComputeResourceStatusUnavailable      ComputerResourceStatus = "unavailable"
)

type ComputeResourceBalanceStrategy string

const (
	ComputeResourceBalanceStrategyRoundRobin           ComputeResourceBalanceStrategy = "round-robin"
	ComputeResourceBalanceStrategyRandom               ComputeResourceBalanceStrategy = "random"
	ComputeResourceBalanceStrategyMostStorageAvailable ComputeResourceBalanceStrategy = "most-storage-available"
)

type ComputeResource struct {
	ID        int                     `json:"id"`
	Name      string                  `json:"name"`
	Host      string                  `json:"host"`
	AgentPort int                     `json:"agent_port"`
	Settings  ComputeResourceSettings `json:"settings"`
	Status    ComputerResourceStatus  `json:"status"`
	Locations []Location              `json:"locations"`
	IPBlocks  []IPBlock               `json:"ip_blocks"`
	Metrics   ComputeResourceMetrics  `json:"metrics"`
	VMsCount  int                     `json:"vms_count"`
	Version   string                  `json:"version"`
}

type ComputeResourceSettings struct {
	CachePath       string                         `json:"cache_path"`
	ISOPath         string                         `json:"iso_path"`
	BackupTmpPath   string                         `json:"backup_tmp_path"`
	VNCProxyPort    int                            `json:"vnc_proxy_port"`
	Limits          ComputeResourceSettingsLimits  `json:"limits"`
	Network         ComputeResourceSettingsNetwork `json:"network"`
	BalanceStrategy ComputeResourceBalanceStrategy `json:"balance_strategy"`
}

type ComputeResourceSettingsLimits struct {
	VM   ComputeResourceSettingsLimit `json:"vm"`
	HDD  ComputeResourceSettingsLimit `json:"hdd"`
	RAM  ComputeResourceSettingsLimit `json:"ram"`
	VCPU ComputeResourceSettingsLimit `json:"vcpu"`
}

type ComputeResourceSettingsLimit struct {
	Unlimited bool    `json:"unlimited"`
	Total     float32 `json:"total"`
	Used      float32 `json:"used,omitempty"`
}

type ComputeResourceSettingsNetworkType string

const (
	ComputeResourceSettingsNetworkTypeRouted  ComputeResourceSettingsNetworkType = "routed"
	ComputeResourceSettingsNetworkTypeBridged ComputeResourceSettingsNetworkType = "bridged"
)

type ComputeResourceSettingsNetwork struct {
	Type    ComputeResourceSettingsNetworkType     `json:"type"`
	Bridges []ComputeResourceSettingsNetworkBridge `json:"bridges"`
}

type ComputeResourceSettingsNetworkBridge struct {
	Type ComputeResourceSettingsNetworkType `json:"type"`
	Name string                             `json:"name"`
}

type ComputeResourceMetrics struct {
	Network ComputeResourceMetricsNetwork `json:"network"`
}

type ComputeResourceMetricsNetwork struct {
	IPv6Enabled bool `json:"ipv6_enabled"`
}

type ComputeResourceNetwork struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// AddrConfType for 'static' or 'dhcp'
	AddrConfType string `json:"addr_conf_type"`
	IPVersion    int    `json:"ip_version"`
	IP           string `json:"ip"`
	Mask         string `json:"mask"`
	MaskSize     int    `json:"mask_size"`
}

type ComputeResourceAuthType string

const (
	ComputeResourceAuthTypePassword ComputeResourceAuthType = "lpass"
	ComputeResourceAuthTypeKey      ComputeResourceAuthType = "lkey"
)

type ComputerResourceCreateRequest struct {
	Name  string `json:"name,omitempty"`
	Host  string `json:"host,omitempty"`
	Login string `json:"login,omitempty"`
	// SSH port number
	Port     int                     `json:"port,omitempty"`
	Type     ComputeResourceAuthType `json:"type,omitempty"`
	Password string                  `json:"password,omitempty"`
	// SSH private key
	Key       string `json:"key,omitempty"`
	AgentPort int    `json:"agent_port,omitempty"`
	IPBlocks  []int  `json:"ip_blocks,omitempty"`
	Locations []int  `json:"locations,omitempty"`
}

type SetupNetworkRequest struct {
	ID   string                             `json:"id"`
	Type ComputeResourceSettingsNetworkType `json:"type"`
}

type ComputeResourcePhysicalVolume struct {
	VGFree string `json:"vg_free"`
	VGName string `json:"vg_name"`
	VGSize string `json:"vg_size"`
	PVUsed string `json:"pv_used"`
}

type ComputeResourceThinPool struct {
	ConvertLV       string `json:"convert_lv"`
	CopyPercent     string `json:"copy_percent"`
	DataPercent     string `json:"data_percent"`
	LVAttr          string `json:"lv_attr"`
	LVLayout        string `json:"lv_layout"`
	LVMetadataSize  string `json:"lv_metadata_size"`
	LVName          string `json:"lv_name"`
	LVSize          string `json:"lv_size"`
	MetadataPrecent string `json:"metadata_percent"`
	MirrorLog       string `json:"mirror_log"`
	MovePV          string `json:"move_pv"`
	Origin          string `json:"origin"`
	PoolLV          string `json:"pool_lv"`
	VGName          string `json:"vg_name"`
}

type computeResourceResponse struct {
	Data ComputeResource `json:"data"`
}

type ComputeResourcesPaginatedResponse struct {
	paginatedResponse

	Data []ComputeResource `json:"data"`
}

func (s *ComputeResourcesService) List(ctx context.Context) (ComputeResourcesPaginatedResponse, error) {
	resp := ComputeResourcesPaginatedResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "compute_resources", &resp)
}

func (s *ComputeResourcesService) Create(
	ctx context.Context,
	data ComputerResourceCreateRequest,
) (ComputeResource, error) {
	var resp computeResourceResponse
	return resp.Data, s.client.create(ctx, "compute_resources", data, &resp)
}

func (s *ComputeResourcesService) Patch(
	ctx context.Context,
	id int,
	data ComputerResourceCreateRequest,
) (ComputeResource, error) {
	var resp computeResourceResponse
	return resp.Data, s.client.patch(ctx, fmt.Sprintf("compute_resources/%d", id), data, &resp)
}

func (s *ComputeResourcesService) Get(ctx context.Context, id int) (ComputeResource, error) {
	var resp computeResourceResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("compute_resources/%d", id), &resp)
}

type deleteRequest struct {
	Force bool `json:"force"`
}

func (s *ComputeResourcesService) Delete(ctx context.Context, id int, force bool) error {
	data := deleteRequest{
		Force: force,
	}
	path := fmt.Sprintf("compute_resources/%d", id)
	body, code, err := s.client.request(ctx, http.MethodDelete, path, withBody(data))
	if err != nil {
		return err
	}

	if code != http.StatusNoContent {
		return newHTTPError(http.MethodDelete, path, code, body)
	}
	return nil
}

func (s *ComputeResourcesService) Networks(ctx context.Context, id int) ([]ComputeResourceNetwork, error) {
	var resp struct {
		Data []ComputeResourceNetwork `json:"data"`
	}
	return resp.Data, s.client.get(ctx, fmt.Sprintf("compute_resources/%d/networks", id), &resp)
}

func (s *ComputeResourcesService) SetUpNetwork(ctx context.Context, id int, data SetupNetworkRequest) error {
	path := fmt.Sprintf("compute_resources/%d/setup_network", id)
	body, code, err := s.client.request(ctx, http.MethodPost, path, withBody(data))
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return newHTTPError(http.MethodPost, path, code, body)
	}
	return nil
}

func (s *ComputeResourcesService) PhysicalVolumes(
	ctx context.Context,
	id int,
) ([]ComputeResourcePhysicalVolume, error) {
	var resp struct {
		Data []ComputeResourcePhysicalVolume `json:"data"`
	}
	return resp.Data, s.client.get(ctx, fmt.Sprintf("compute_resources/%d/physical_volumes", id), &resp)
}

func (s *ComputeResourcesService) ThinPools(ctx context.Context, id int) ([]ComputeResourceThinPool, error) {
	var resp struct {
		Data []ComputeResourceThinPool `json:"data"`
	}
	return resp.Data, s.client.get(ctx, fmt.Sprintf("compute_resources/%d/thin_pools", id), &resp)
}
