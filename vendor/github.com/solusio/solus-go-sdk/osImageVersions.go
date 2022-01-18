package solus

import (
	"context"
	"fmt"
)

// OsImageVersionsService handles all available methods with OS image versions.
type OsImageVersionsService service

// OsImageVersion represents an OS image version.
type OsImageVersion struct {
	ID                 int                `json:"id"`
	VirtualizationType VirtualizationType `json:"virtualization_type"`
	Position           float64            `json:"position"`
	Version            string             `json:"version"`
	URL                string             `json:"url"`
	CloudInitVersion   CloudInitVersion   `json:"cloud_init_version,omitempty"`
	OsImageID          int                `json:"os_image_id"`
	IsVisible          bool               `json:"is_visible"`
	IsSSHKeysSupported bool               `json:"is_ssh_keys_supported"`
	AvailablePlans     []ShortPlan        `json:"available_plans"`
}

// ShortOsImageVersion represents only ID and name of OS image version.
type ShortOsImageVersion struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CloudInitVersion represents available cloud-init config versions.
type CloudInitVersion string

const (
	// CloudInitVersionV0 indicates v0 cloud-init config version.
	CloudInitVersionV0 CloudInitVersion = "v0"

	// CloudInitVersionCentOS6 indicates CentOS 6 specific v0 cloud-init config
	// version.
	CloudInitVersionCentOS6 CloudInitVersion = "v0-centos6"

	// CloudInitVersionDebian9 indicates Debian 9 specific v0 cloud-init config
	// version.
	CloudInitVersionDebian9 CloudInitVersion = "v0-debian9"

	// CloudInitVersionV2 indicates v2 cloud-init config version.
	CloudInitVersionV2 CloudInitVersion = "v2"

	// CloudInitVersionV2Alpine indicates Alpine specific v2 cloud-init config
	// version.
	CloudInitVersionV2Alpine CloudInitVersion = "v2-alpine"

	// CloudInitVersionV2Centos indicates CentOS 7 & 8 specific v2 cloud-init config
	// version.
	CloudInitVersionV2Centos CloudInitVersion = "v2-centos"

	// CloudInitVersionV2Debian10 indicates Debian 10 specific v2 cloud-init config
	// version.
	CloudInitVersionV2Debian10 CloudInitVersion = "v2-debian10"

	// CloudInitVersionCloudBase indicates cloudbase config version.
	CloudInitVersionCloudBase CloudInitVersion = "cloudbase"
)

// IsValidCloudInitVersion returns true if specified cloud-init version is valid.
func IsValidCloudInitVersion(v string) bool {
	_, ok := map[CloudInitVersion]struct{}{
		CloudInitVersionV0:         {},
		CloudInitVersionCentOS6:    {},
		CloudInitVersionDebian9:    {},
		CloudInitVersionV2:         {},
		CloudInitVersionV2Alpine:   {},
		CloudInitVersionV2Centos:   {},
		CloudInitVersionV2Debian10: {},
		CloudInitVersionCloudBase:  {},
	}[CloudInitVersion(v)]
	return ok
}

// OsImageVersionRequest represents available properties for creating a new OS image
// version.
type OsImageVersionRequest struct {
	Position           float64            `json:"position,omitempty"`
	Version            string             `json:"version"`
	VirtualizationType VirtualizationType `json:"virtualization_type"`
	URL                string             `json:"url"`
	CloudInitVersion   CloudInitVersion   `json:"cloud_init_version,omitempty"`
	IsVisible          bool               `json:"is_visible"`
	AvailablePlans     []int              `json:"available_plans,omitempty"`
}

type osImageVersionResponse struct {
	Data OsImageVersion `json:"data"`
}

// Get gets specified OS image version.
func (s *OsImageVersionsService) Get(ctx context.Context, id int) (OsImageVersion, error) {
	var resp osImageVersionResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("os_image_versions/%d", id), &resp)
}

// Update updates specified OS image version.
func (s *OsImageVersionsService) Update(
	ctx context.Context,
	id int,
	data OsImageVersionRequest,
) (OsImageVersion, error) {
	var resp osImageVersionResponse
	return resp.Data, s.client.update(ctx, fmt.Sprintf("os_image_versions/%d", id), data, &resp)
}

// Delete deletes specified OS image version.
func (s *OsImageVersionsService) Delete(ctx context.Context, id int) error {
	return s.client.syncDelete(ctx, fmt.Sprintf("os_image_versions/%d", id))
}
