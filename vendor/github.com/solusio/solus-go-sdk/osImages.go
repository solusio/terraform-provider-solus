package solus

import (
	"context"
	"fmt"
)

type OsImagesService service

type OsImage struct {
	ID        int              `json:"id"`
	Name      string           `json:"name"`
	Icon      Icon             `json:"icon"`
	Versions  []OsImageVersion `json:"versions,omitempty"`
	IsDefault bool             `json:"is_default,omitempty"`
}

type OsImageVersion struct {
	ID               int     `json:"id"`
	Position         float64 `json:"position"`
	Version          string  `json:"version"`
	URL              string  `json:"url"`
	CloudInitVersion string  `json:"cloud_init_version"`
}

type OsImageCreateRequest struct {
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	IsVisible bool   `json:"is_visible"`
}

type OsImageVersionRequest struct {
	Position         float64 `json:"position"`
	Version          string  `json:"version"`
	URL              string  `json:"url"`
	CloudInitVersion string  `json:"cloud_init_version"`
}

type OsImagesResponse struct {
	paginatedResponse

	Data []OsImage `json:"data"`
}

func (s *OsImagesService) List(ctx context.Context) (OsImagesResponse, error) {
	resp := OsImagesResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "os_images", &resp)
}

func (s *OsImagesService) Create(ctx context.Context, data OsImageCreateRequest) (OsImage, error) {
	var resp struct {
		Data OsImage `json:"data"`
	}
	return resp.Data, s.client.create(ctx, "os_images", data, &resp)
}

func (s *OsImagesService) Delete(ctx context.Context, id int) error {
	return s.client.delete(ctx, fmt.Sprintf("os_images/%d", id))
}

func (s *OsImagesService) OsImageVersionCreate(
	ctx context.Context,
	osImageID int,
	data OsImageVersionRequest,
) (OsImageVersion, error) {
	var resp struct {
		Data OsImageVersion `json:"data"`
	}
	return resp.Data, s.client.create(ctx, fmt.Sprintf("os_images/%d/versions", osImageID), data, &resp)
}
