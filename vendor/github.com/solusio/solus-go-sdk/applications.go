package solus

import (
	"context"
)

type ApplicationsService service

type LoginLinkType string

const (
	LoginLinkTypeNone   LoginLinkType = "none"
	LoginLinkTypeURL    LoginLinkType = "url"
	LoginLinkTypeJSCode LoginLinkType = "js_code"
	LoginLinkTypeInfo   LoginLinkType = "info"
)

type LoginLink struct {
	Type    LoginLinkType `json:"type"`
	Content string        `json:"content"`
}

type Application struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Icon             Icon      `json:"icon"`
	URL              string    `json:"url"`
	CloudInitVersion string    `json:"cloud_init_version"`
	UserData         string    `json:"user_data_template"`
	LoginLink        LoginLink `json:"login_link"`
	JSONSchema       string    `json:"json_schema"`
	IsDefault        bool      `json:"is_default"`
	IsVisible        bool      `json:"is_visible"`
	IsBuiltin        bool      `json:"is_buildin"`
}

type ApplicationCreateRequest struct {
	Name             string    `json:"name"`
	URL              string    `json:"url"`
	IconID           int       `json:"icon_id"`
	CloudInitVersion string    `json:"cloud_init_version"`
	UserDataTemplate string    `json:"user_data_template"`
	JSONSchema       string    `json:"json_schema"`
	IsVisible        bool      `json:"is_visible"`
	LoginLink        LoginLink `json:"login_link"`
}

type ApplicationsResponse struct {
	paginatedResponse

	Data []Application `json:"data"`
}

func (s *ApplicationsService) Create(ctx context.Context, data ApplicationCreateRequest) (Application, error) {
	var resp struct {
		Data Application `json:"data"`
	}
	return resp.Data, s.client.create(ctx, "applications", data, &resp)
}

func (s *ApplicationsService) List(ctx context.Context) (ApplicationsResponse, error) {
	resp := ApplicationsResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "applications", &resp)
}
