package solus

import (
	"context"
	"fmt"
)

// IPBlocksService handles all available methods with IP blocks.
type IPBlocksService service

// IPBlock represents an IP block.
type IPBlock struct {
	ID               int                `json:"id"`
	Name             string             `json:"name"`
	Type             IPVersion          `json:"type"`
	Gateway          string             `json:"gateway"`
	Netmask          string             `json:"netmask"`
	Ns1              string             `json:"ns_1"`
	Ns2              string             `json:"ns_2"`
	From             string             `json:"from"`
	To               string             `json:"to"`
	Subnet           int                `json:"subnet"`
	Range            string             `json:"range"`
	ComputeResources []ComputeResource  `json:"compute_resources[]"`
	IPs              []IPBlockIPAddress `json:"ips[]"`
	ReverseDNS       IPBlockReverseDNS  `json:"reverse_dns"`
}

// IPVersion represents available IP versions.
type IPVersion string

const (
	// IPv4 indicates IP v4 address.
	IPv4 IPVersion = "IPv4"

	// IPv6 indicates IP v6 address.
	IPv6 IPVersion = "IPv6"
)

// IPBlockRequest represents available properties for creating new or updating
// existing IP block.
type IPBlockRequest struct {
	ComputeResources []int             `json:"compute_resources,omitempty"`
	Name             string            `json:"name"`
	Type             IPVersion         `json:"type"`
	Gateway          string            `json:"gateway"`
	Ns1              string            `json:"ns_1"`
	Ns2              string            `json:"ns_2"`
	ReverseDNS       IPBlockReverseDNS `json:"reverse_dns"`

	// IPv4 related fields
	Netmask string `json:"netmask,omitempty"`
	From    string `json:"from,omitempty"`
	To      string `json:"to,omitempty"`

	// IPv6 related fields
	Range  string `json:"range,omitempty"`
	Subnet int    `json:"subnet,omitempty"`
}

// IPBlockReverseDNS represents an IP block's reverse DNS settings.
type IPBlockReverseDNS struct {
	Zone    string `json:"zone"`
	Enabled bool   `json:"enabled"`
}

// IPBlockIPAddress represents an IP block's IP address.
type IPBlockIPAddress struct {
	ID      int     `json:"id"`
	IP      string  `json:"ip"`
	IPBlock IPBlock `json:"ip_block"`
}

// IPBlocksResponse represents paginated list of IP blocks.
// This cursor can be used for iterating over all available IP blocks.
type IPBlocksResponse struct {
	paginatedResponse

	Data []IPBlock `json:"data"`
}

type ipBlockResponse struct {
	Data IPBlock `json:"data"`
}

// List lists IP blocks.
func (s *IPBlocksService) List(ctx context.Context, filter *FilterIPBlocks) (IPBlocksResponse, error) {
	resp := IPBlocksResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "ip_blocks", &resp, withFilter(filter.data))
}

// Get gets specified IP block.
func (s *IPBlocksService) Get(ctx context.Context, id int) (IPBlock, error) {
	var resp ipBlockResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("ip_blocks/%d", id), &resp)
}

// Create creates new IP block.
func (s *IPBlocksService) Create(ctx context.Context, data IPBlockRequest) (IPBlock, error) {
	var resp ipBlockResponse
	return resp.Data, s.client.create(ctx, "ip_blocks", data, &resp)
}

// Update updates specified IP block.
func (s *IPBlocksService) Update(ctx context.Context, id int, data IPBlockRequest) (IPBlock, error) {
	var resp ipBlockResponse
	return resp.Data, s.client.update(ctx, fmt.Sprintf("ip_blocks/%d", id), data, &resp)
}

// Delete deletes specified IP block.
func (s *IPBlocksService) Delete(ctx context.Context, id int) error {
	return s.client.syncDelete(ctx, fmt.Sprintf("ip_blocks/%d", id))
}

// IPAddressCreate creates a new IP address in the specified IP block.
func (s *IPBlocksService) IPAddressCreate(ctx context.Context, ipBlockID int) (IPBlockIPAddress, error) {
	var resp struct {
		Data IPBlockIPAddress `json:"data"`
	}
	return resp.Data, s.client.create(ctx, fmt.Sprintf("ip_blocks/%d/ips", ipBlockID), nil, &resp)
}

// IPAddressDelete deletes a provided IP address in the specified IP block.
func (s *IPBlocksService) IPAddressDelete(ctx context.Context, id int) error {
	return s.client.syncDelete(ctx, fmt.Sprintf("ips/%d", id))
}
