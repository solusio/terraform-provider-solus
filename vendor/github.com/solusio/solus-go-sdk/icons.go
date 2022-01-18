package solus

import (
	"context"
	"fmt"
)

// IconsService handles all available methods with icons.
type IconsService service

// Icon represents an icon.
type Icon struct {
	ID   int      `json:"id"`
	Name string   `json:"name"`
	URL  string   `json:"url"`
	Type IconType `json:"type"`
}

// IconType represents available icon types.
type IconType string

const (
	// IconTypeOS indicated OS specific icon.
	// OSes like AlmaLinux, Ubuntu, and etc.
	IconTypeOS IconType = "os"

	// IconTypeApplication indicates application specific icon.
	// Applications likes Plesk, Nginx, and etc.
	IconTypeApplication IconType = "application"

	// IconTypeFlags indicates countries flags.
	IconTypeFlags IconType = "flags"
)

// IconsResponse represents paginated list of icons.
// This cursor can be used for iterating over all available icons.
type IconsResponse struct {
	paginatedResponse

	Data []Icon `json:"data"`
}

// List lists icons.
func (s *IconsService) List(ctx context.Context, filter *FilterIcons) (IconsResponse, error) {
	resp := IconsResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "icons", &resp, withFilter(filter.data))
}

// Get gets specified icon.
func (s *IconsService) Get(ctx context.Context, id int) (Icon, error) {
	var resp struct {
		Data Icon `json:"data"`
	}
	return resp.Data, s.client.get(ctx, fmt.Sprintf("icons/%d", id), &resp)
}
