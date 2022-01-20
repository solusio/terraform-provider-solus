package solus

import (
	"context"
	"fmt"

	"gopkg.in/guregu/null.v4"
)

// OsImagesService handles all available methods with OS image.
type OsImagesService service

// OsImage represent an OS image.
type OsImage struct {
	ID        int              `json:"id"`
	Name      string           `json:"name"`
	Icon      Icon             `json:"icon"`
	Versions  []OsImageVersion `json:"versions,omitempty"`
	IsDefault bool             `json:"is_default,omitempty"`
	IsVisible bool             `json:"is_visible,omitempty"`
	Position  float32          `json:"position"`
}

// OsImageRequest represents available properties for creating a new OS image.
type OsImageRequest struct {
	Name      string   `json:"name"`
	IconID    null.Int `json:"icon_id"`
	IsVisible bool     `json:"is_visible"`
}

// OsImagesResponse represents paginated list of OS images.
// This cursor can be used for iterating over all available OS images.
type OsImagesResponse struct {
	paginatedResponse

	Data []OsImage `json:"data"`
}

type osImageResponse struct {
	Data OsImage `json:"data"`
}

// List lists OS images.
func (s *OsImagesService) List(ctx context.Context, filter *FilterOsImages) (OsImagesResponse, error) {
	resp := OsImagesResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "os_images", &resp, withFilter(filter.data))
}

// Get gets specified OS image.
func (s *OsImagesService) Get(ctx context.Context, id int) (OsImage, error) {
	var resp osImageResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("os_images/%d", id), &resp)
}

// Create creates specified OS image.
func (s *OsImagesService) Create(ctx context.Context, data OsImageRequest) (OsImage, error) {
	var resp osImageResponse
	return resp.Data, s.client.create(ctx, "os_images", data, &resp)
}

// Update updates specified OS image.
func (s *OsImagesService) Update(ctx context.Context, id int, data OsImageRequest) (OsImage, error) {
	var resp osImageResponse
	return resp.Data, s.client.update(ctx, fmt.Sprintf("os_images/%d", id), data, &resp)
}

// Delete deletes specified OS image.
func (s *OsImagesService) Delete(ctx context.Context, id int) error {
	return s.client.syncDelete(ctx, fmt.Sprintf("os_images/%d", id))
}

// CreateVersion creates a new version for the specified OS image.
func (s *OsImagesService) CreateVersion(
	ctx context.Context,
	osImageID int,
	data OsImageVersionRequest,
) (OsImageVersion, error) {
	var resp osImageVersionResponse
	return resp.Data, s.client.create(ctx, fmt.Sprintf("os_images/%d/versions", osImageID), data, &resp)
}

// ListVersion lists specified OS image versions.
func (s *OsImagesService) ListVersion(ctx context.Context, osImageID int) ([]OsImageVersion, error) {
	var resp struct {
		Data []OsImageVersion `json:"data"`
	}
	return resp.Data, s.client.list(ctx, fmt.Sprintf("os_images/%d/versions", osImageID), &resp)
}
