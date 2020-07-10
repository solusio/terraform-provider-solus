package solus

import (
	"context"
	"fmt"
	"net/http"
)

type IPBlocksService service

type IPVersion string

const (
	IPv4 IPVersion = "IPv4"
	IPv6 IPVersion = "IPv6"
)

type IPBlockCreateRequest struct {
	ComputeResources []int     `json:"compute_resources,omitempty"`
	Name             string    `json:"name"`
	Type             IPVersion `json:"type"`
	Gateway          string    `json:"gateway"`
	Ns1              string    `json:"ns_1"`
	Ns2              string    `json:"ns_2"`

	// IPv4 related fields
	Netmask string `json:"netmask"`
	From    string `json:"from"`
	To      string `json:"to"`

	// IPv6 related fields
	Range  string `json:"range"`
	Subnet int    `json:"subnet"`
}

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
	ComputeResources []ComputeResource  `json:"compute_resources[]"`
	IPs              []IPBlockIPAddress `json:"ips[]"`
}

type IPBlockIPAddress struct {
	ID      int     `json:"id"`
	IP      string  `json:"ip"`
	IPBlock IPBlock `json:"ip_block"`
}

type IPBlocksResponse struct {
	paginatedResponse

	Data []IPBlock `json:"data"`
}

func (s *IPBlocksService) List(ctx context.Context) (IPBlocksResponse, error) {
	resp := IPBlocksResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "ip_blocks", &resp)
}

func (s *IPBlocksService) Create(ctx context.Context, data IPBlockCreateRequest) (IPBlock, error) {
	var resp struct {
		Data IPBlock `json:"data"`
	}
	return resp.Data, s.client.create(ctx, "ip_blocks", data, &resp)
}

func (s *IPBlocksService) Delete(ctx context.Context, id int) error {
	return s.client.delete(ctx, fmt.Sprintf("ip_blocks/%d", id))
}

func (s *IPBlocksService) IPAddressCreate(ctx context.Context, ipBlockID int) (IPBlockIPAddress, error) {
	path := fmt.Sprintf("ip_blocks/%d/ips", ipBlockID)
	body, code, err := s.client.request(ctx, http.MethodPost, path)
	if err != nil {
		return IPBlockIPAddress{}, err
	}

	if code != http.StatusCreated {
		return IPBlockIPAddress{}, newHTTPError(http.MethodPost, path, code, body)
	}

	var resp struct {
		Data IPBlockIPAddress `json:"data"`
	}
	return resp.Data, unmarshal(body, &resp)
}

func (s *IPBlocksService) IPAddressDelete(ctx context.Context, id int) error {
	return s.client.delete(ctx, fmt.Sprintf("ips/%d", id))
}
