package solus

import (
	"context"
	"fmt"

	"gopkg.in/guregu/null.v4"
)

type OsImagesService service

type OsImage struct {
	ID        int              `json:"id"`
	Name      string           `json:"name"`
	Icon      Icon             `json:"icon"`
	Versions  []OsImageVersion `json:"versions,omitempty"`
	IsDefault bool             `json:"is_default,omitempty"`
	IsVisible bool             `json:"is_visible,omitempty"`
	Position  float32          `json:"position"`
}

type OsImageRequest struct {
	Name      string   `json:"name"`
	IconID    null.Int `json:"icon_id"`
	IsVisible bool     `json:"is_visible"`
}

type OsImagesResponse struct {
	paginatedResponse

	Data []OsImage `json:"data"`
}

type osImageResponse struct {
	Data OsImage `json:"data"`
}

func (s *OsImagesService) List(ctx context.Context, filter *FilterOsImages) (OsImagesResponse, error) {
	resp := OsImagesResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "os_images", &resp, withFilter(filter.data))
}

func (s *OsImagesService) Get(ctx context.Context, id int) (OsImage, error) {
	var resp osImageResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("os_images/%d", id), &resp)
}

func (s *OsImagesService) Create(ctx context.Context, data OsImageRequest) (OsImage, error) {
	var resp osImageResponse
	return resp.Data, s.client.create(ctx, "os_images", data, &resp)
}

func (s *OsImagesService) Update(ctx context.Context, id int, data OsImageRequest) (OsImage, error) {
	var resp osImageResponse
	return resp.Data, s.client.update(ctx, fmt.Sprintf("os_images/%d", id), data, &resp)
}

func (s *OsImagesService) Delete(ctx context.Context, id int) error {
	return s.client.delete(ctx, fmt.Sprintf("os_images/%d", id))
}

func (s *OsImagesService) CreateVersion(
	ctx context.Context,
	osImageID int,
	data OsImageVersionRequest,
) (OsImageVersion, error) {
	var resp osImageVersionResponse
	return resp.Data, s.client.create(ctx, fmt.Sprintf("os_images/%d/versions", osImageID), data, &resp)
}

func (s *OsImagesService) ListVersion(ctx context.Context, osImageID int) ([]OsImageVersion, error) {
	var resp struct {
		Data []OsImageVersion `json:"data"`
	}
	return resp.Data, s.client.list(ctx, fmt.Sprintf("os_images/%d/versions", osImageID), &resp)
}
