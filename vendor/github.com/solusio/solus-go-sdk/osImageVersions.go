package solus

import (
	"context"
	"fmt"
)

type OsImageVersionsService service

type CloudInitVersion string

const (
	CloudInitVersionV0         CloudInitVersion = "v0"
	CloudInitVersionCentOS6    CloudInitVersion = "v0-centos6"
	CloudInitVersionDebian9    CloudInitVersion = "v0-debian9"
	CloudInitVersionV2         CloudInitVersion = "v2"
	CloudInitVersionV2Alpine   CloudInitVersion = "v2-alpine"
	CloudInitVersionV2Centos   CloudInitVersion = "v2-centos"
	CloudInitVersionV2Debian10 CloudInitVersion = "v2-debian10"
	CloudInitVersionCloudBase  CloudInitVersion = "cloudbase"
)

func IsValidCloudInitVersion(v string) bool {
	m := map[CloudInitVersion]struct{}{
		CloudInitVersionV0:         {},
		CloudInitVersionCentOS6:    {},
		CloudInitVersionDebian9:    {},
		CloudInitVersionV2:         {},
		CloudInitVersionV2Alpine:   {},
		CloudInitVersionV2Centos:   {},
		CloudInitVersionV2Debian10: {},
		CloudInitVersionCloudBase:  {},
	}

	_, ok := m[CloudInitVersion(v)]
	return ok
}

type OsImageVersion struct {
	ID                 int              `json:"id"`
	Position           float64          `json:"position"`
	Version            string           `json:"version"`
	URL                string           `json:"url"`
	CloudInitVersion   CloudInitVersion `json:"cloud_init_version"`
	OsImageID          int              `json:"os_image_id"`
	IsSSHKeysSupported bool             `json:"is_ssh_keys_supported"`
}

type OsImageVersionRequest struct {
	Position         float64          `json:"position,omitempty"`
	Version          string           `json:"version"`
	URL              string           `json:"url"`
	CloudInitVersion CloudInitVersion `json:"cloud_init_version"`
}

type osImageVersionResponse struct {
	Data OsImageVersion `json:"data"`
}

func (s *OsImageVersionsService) Get(ctx context.Context, id int) (OsImageVersion, error) {
	var resp osImageVersionResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("os_image_versions/%d", id), &resp)
}

func (s *OsImageVersionsService) Update(
	ctx context.Context,
	id int,
	data OsImageVersionRequest,
) (OsImageVersion, error) {
	var resp osImageVersionResponse
	return resp.Data, s.client.update(ctx, fmt.Sprintf("os_image_versions/%d", id), data, &resp)
}

func (s *OsImageVersionsService) Delete(ctx context.Context, id int) error {
	return s.client.delete(ctx, fmt.Sprintf("os_image_versions/%d", id))
}
