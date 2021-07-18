package solus

import (
	"context"
	"fmt"
)

type IconsService service

type IconType string

const (
	IconTypeOS          IconType = "os"
	IconTypeApplication IconType = "application"
	IconTypeFlags       IconType = "flags"
)

type Icon struct {
	ID   int      `json:"id"`
	Name string   `json:"name"`
	URL  string   `json:"url"`
	Type IconType `json:"type"`
}

type IconsResponse struct {
	paginatedResponse

	Data []Icon `json:"data"`
}

func (s *IconsService) List(ctx context.Context, filter *FilterIcons) (IconsResponse, error) {
	resp := IconsResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "icons", &resp, withFilter(filter.data))
}

func (s *IconsService) Get(ctx context.Context, id int) (Icon, error) {
	var resp struct {
		Data Icon `json:"data"`
	}
	return resp.Data, s.client.get(ctx, fmt.Sprintf("icons/%d", id), &resp)
}
