package solus

import (
	"context"
	"fmt"
	"net/http"
)

// VirtualizationType represents available virtualization types.
type VirtualizationType string

const (
	// VirtualizationTypeKVM indicates KVM virtualization type.
	VirtualizationTypeKVM VirtualizationType = "kvm"

	// VirtualizationTypeVZ indicates VZ virtualization type.
	VirtualizationTypeVZ VirtualizationType = "vz"
)

// IsValidVirtualizationType returns true if specified virtualization type is valid.
func IsValidVirtualizationType(v string) bool {
	_, ok := map[VirtualizationType]struct{}{
		VirtualizationTypeKVM: {},
		VirtualizationTypeVZ:  {},
	}[VirtualizationType(v)]
	return ok
}

// ComputeResourcesService handles all available methods with compute
// resources.
type ComputeResourcesService service

// ComputeResource represents a compute resource.
// The compute resource is a server where virtual machines will be running.
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

// ComputerResourceStatus represents available compute resource's statuses.
type ComputerResourceStatus string

//goland:noinspection GoUnusedConst
const (
	// ComputeResourceStatusActive indicates compute resource is active and ready
	// to handle requests.
	ComputeResourceStatusActive ComputerResourceStatus = "active"

	// ComputeResourceStatusCommissioning indicates compute resource is still
	// commissioning.
	ComputeResourceStatusCommissioning ComputerResourceStatus = "commissioning"

	// ComputeResourceStatusConfigureNetwork indicates compute resource is ready
	// for setting up a network.
	ComputeResourceStatusConfigureNetwork ComputerResourceStatus = "configure_network"

	// ComputeResourceStatusFailed indicates compute resource commissioning
	// was failed.
	ComputeResourceStatusFailed ComputerResourceStatus = "failed"

	// ComputeResourceStatusUnavailable indicates compute resource is unavailable
	// due to some reason.
	ComputeResourceStatusUnavailable ComputerResourceStatus = "unavailable"
)

// ComputeResourceBalanceStrategy represents available balancing strategies for
// compute resource.
// These strategies will be used during a server creation for choosing which
// compute resource should be used.
type ComputeResourceBalanceStrategy string

//goland:noinspection GoUnusedConst
const (
	// ComputeResourceBalanceStrategyRoundRobin indicates compute resource will
	// be choosing by round-robin algorithm.
	ComputeResourceBalanceStrategyRoundRobin ComputeResourceBalanceStrategy = "round-robin"

	// ComputeResourceBalanceStrategyRandom indicates compute resource will be
	// chosen randomly.
	ComputeResourceBalanceStrategyRandom ComputeResourceBalanceStrategy = "random"

	// ComputeResourceBalanceStrategyMostStorageAvailable indicates compute resource
	// with most available storage space will be chosen.
	ComputeResourceBalanceStrategyMostStorageAvailable ComputeResourceBalanceStrategy = "most-storage-available"
)

// ComputeResourceSettings represents available compute resource's settings.
type ComputeResourceSettings struct {
	CachePath           string                         `json:"cache_path"`
	ISOPath             string                         `json:"iso_path"`
	BackupTmpPath       string                         `json:"backup_tmp_path"`
	VNCProxyPort        int                            `json:"vnc_proxy_port"`
	Limits              ComputeResourceSettingsLimits  `json:"limits"`
	Network             ComputeResourceSettingsNetwork `json:"network"`
	BalanceStrategy     ComputeResourceBalanceStrategy `json:"balance_strategy"`
	VirtualizationTypes []VirtualizationType           `json:"virtualization_types"`
}

// ComputeResourceSettingsLimits represents available compute resources limits
// for virtual machines.
type ComputeResourceSettingsLimits struct {
	VM   ComputeResourceSettingsLimit `json:"vm"`
	HDD  ComputeResourceSettingsLimit `json:"hdd"`
	RAM  ComputeResourceSettingsLimit `json:"ram"`
	VCPU ComputeResourceSettingsLimit `json:"vcpu"`
}

// ComputeResourceSettingsLimit represents single compute resources limit.
type ComputeResourceSettingsLimit struct {
	Unlimited bool    `json:"unlimited"`
	Total     float32 `json:"total"`
	Used      float32 `json:"used,omitempty"`
}

// ComputeResourceSettingsNetworkType represents available network types on a
// compute resource.
type ComputeResourceSettingsNetworkType string

//goland:noinspection GoUnusedConst
const (
	// ComputeResourceSettingsNetworkTypeRouted indicates virtual servers don't
	// connect directly to the physical network. The compute resource's operating
	// system routes the servers' traffic to the physical network (the compute
	// resource works as the gateway).
	// The server's MAC address isn’t exposed to the physical network.
	ComputeResourceSettingsNetworkTypeRouted ComputeResourceSettingsNetworkType = "routed"

	// ComputeResourceSettingsNetworkTypeBridged indicates virtual servers get
	// direct access to the physical network. In the bridged network, the IP addresses
	// of a server and the gateway must be within the same IP network. For example,
	// if the gateway's IP address 192.168.1.1 is within the IP network 192.168.1.0/24,
	// then the server's IP address must be also within the 192.168.1.0/24 network
	// (for example, 192.168.1.2).
	// The server's MAC address is exposed to the physical network.
	// Use a bridged network if you have the large network of IP addresses with
	// the gateway.
	ComputeResourceSettingsNetworkTypeBridged ComputeResourceSettingsNetworkType = "bridged"
)

// ComputeResourceSettingsNetwork represents compute resource's network.
type ComputeResourceSettingsNetwork struct {
	Type    ComputeResourceSettingsNetworkType     `json:"type"`
	Bridges []ComputeResourceSettingsNetworkBridge `json:"bridges"`
}

// ComputeResourceSettingsNetworkBridge represents one of compute resource's network bridge.
type ComputeResourceSettingsNetworkBridge struct {
	Type ComputeResourceSettingsNetworkType `json:"type"`
	Name string                             `json:"name"`
}

// ComputeResourceMetrics represents compute resource's metrics.
type ComputeResourceMetrics struct {
	Network ComputeResourceMetricsNetwork `json:"network"`
}

// ComputeResourceMetricsNetwork represents compute resource's network metrics.
type ComputeResourceMetricsNetwork struct {
	IPv6Enabled bool `json:"ipv6_enabled"`
}

// ComputeResourceNetwork represents compute resource's network.
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

// ComputeResourceAuthType represents available authentication methods for the agent
// installation during compute resource creating.
type ComputeResourceAuthType string

//goland:noinspection GoUnusedConst
const (
	// ComputeResourceAuthTypePassword indicates connect by password.
	ComputeResourceAuthTypePassword ComputeResourceAuthType = "lpass"

	// ComputeResourceAuthTypeKey indicates connect by SSH key.
	ComputeResourceAuthTypeKey ComputeResourceAuthType = "lkey"
)

// ComputerResourceCreateRequest represents available properties for creating a
// new compute resource.
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

// SetupNetworkRequest represents available properties for setting up network for
// the compute resource.
type SetupNetworkRequest struct {
	ID   string                             `json:"id"`
	Type ComputeResourceSettingsNetworkType `json:"type"`
}

// ComputeResourcePhysicalVolume represents a compute resource's physical volume.
type ComputeResourcePhysicalVolume struct {
	VGFree string `json:"vg_free"`
	VGName string `json:"vg_name"`
	VGSize string `json:"vg_size"`
	PVUsed string `json:"pv_used"`
}

// ComputeResourceThinPool represents a compute resource's ThinLVM pool.
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

// ComputeResourcesResponse represents paginated list of compute resources.
// This cursor can be used for iterating over all available compute resources.
type ComputeResourcesResponse struct {
	paginatedResponse

	Data []ComputeResource `json:"data"`
}

type computeResourceResponse struct {
	Data ComputeResource `json:"data"`
}

// List lists compute resource.
func (s *ComputeResourcesService) List(
	ctx context.Context,
	filter *FilterComputeResources,
) (ComputeResourcesResponse, error) {
	resp := ComputeResourcesResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "compute_resources", &resp, withFilter(filter.data))
}

// Create creates new compute resource.
func (s *ComputeResourcesService) Create(
	ctx context.Context,
	data ComputerResourceCreateRequest,
) (ComputeResource, error) {
	var resp computeResourceResponse
	return resp.Data, s.client.create(ctx, "compute_resources", data, &resp)
}

// Patch patches specified compute resource.
func (s *ComputeResourcesService) Patch(
	ctx context.Context,
	id int,
	data ComputerResourceCreateRequest,
) (ComputeResource, error) {
	var resp computeResourceResponse
	return resp.Data, s.client.patch(ctx, fmt.Sprintf("compute_resources/%d", id), data, &resp)
}

// Get gets specified compute resource.
func (s *ComputeResourcesService) Get(ctx context.Context, id int) (ComputeResource, error) {
	var resp computeResourceResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("compute_resources/%d", id), &resp)
}

type deleteRequest struct {
	Force bool `json:"force"`
}

// Delete deletes specified compute resource.
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

// Networks lists specified compute resource's networks.
func (s *ComputeResourcesService) Networks(ctx context.Context, id int) ([]ComputeResourceNetwork, error) {
	var resp struct {
		Data []ComputeResourceNetwork `json:"data"`
	}
	return resp.Data, s.client.get(ctx, fmt.Sprintf("compute_resources/%d/networks", id), &resp)
}

// SetUpNetwork setups a network on the specified compute resource.
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

// PhysicalVolumes lists physical volume on the specified compute resource.
// Return available LVM volume groups.
func (s *ComputeResourcesService) PhysicalVolumes(
	ctx context.Context,
	id int,
) ([]ComputeResourcePhysicalVolume, error) {
	var resp struct {
		Data []ComputeResourcePhysicalVolume `json:"data"`
	}
	return resp.Data, s.client.get(ctx, fmt.Sprintf("compute_resources/%d/physical_volumes", id), &resp)
}

// ThinPools lists ThinLVM pools on the specified compute resource.
func (s *ComputeResourcesService) ThinPools(ctx context.Context, id int) ([]ComputeResourceThinPool, error) {
	var resp struct {
		Data []ComputeResourceThinPool `json:"data"`
	}
	return resp.Data, s.client.get(ctx, fmt.Sprintf("compute_resources/%d/thin_pools", id), &resp)
}
