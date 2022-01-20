package solus

import (
	"context"
)

// ApplicationsService handles all available methods with applications.
type ApplicationsService service

// Application represents an application.
type Application struct {
	ID               int         `json:"id"`
	Name             string      `json:"name"`
	Icon             Icon        `json:"icon"`
	URL              string      `json:"url"`
	CloudInitVersion string      `json:"cloud_init_version"`
	UserData         string      `json:"user_data_template"`
	LoginLink        LoginLink   `json:"login_link"`
	JSONSchema       string      `json:"json_schema"`
	IsDefault        bool        `json:"is_default"`
	IsVisible        bool        `json:"is_visible"`
	IsBuiltin        bool        `json:"is_buildin"`
	AvailablePlans   []ShortPlan `json:"available_plans"`
}

// ShortApplication represents only ID and name of application.
type ShortApplication struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// LoginLinkType a type of login link to the application.
type LoginLinkType string

const (
	// LoginLinkTypeNone indicates application without login link.
	LoginLinkTypeNone LoginLinkType = "none"

	// LoginLinkTypeURL indicates application with URL pattern login link.
	LoginLinkTypeURL LoginLinkType = "url"

	// LoginLinkTypeJSCode indicates application with custom JS code auth code.
	LoginLinkTypeJSCode LoginLinkType = "js_code"

	// LoginLinkTypeInfo indicates application with custom information in popup.
	LoginLinkTypeInfo LoginLinkType = "info"
)

// LoginLink represents an application login link.
type LoginLink struct {
	Type    LoginLinkType `json:"type"`
	Content string        `json:"content"`
}

// ApplicationCreateRequest represents available properties for creating a new
// application.
type ApplicationCreateRequest struct {
	Name             string    `json:"name"`
	URL              string    `json:"url"`
	IconID           int       `json:"icon_id"`
	CloudInitVersion string    `json:"cloud_init_version"`
	UserDataTemplate string    `json:"user_data_template"`
	JSONSchema       string    `json:"json_schema"`
	IsVisible        bool      `json:"is_visible"`
	LoginLink        LoginLink `json:"login_link"`
	AvailablePlans   []int     `json:"available_plans,omitempty"`
}

// ApplicationsResponse represents paginated list of applications.
// This cursor can be used for iterating over all available applications.
type ApplicationsResponse struct {
	paginatedResponse

	Data []Application `json:"data"`
}

// Create creates new application.
func (s *ApplicationsService) Create(ctx context.Context, data ApplicationCreateRequest) (Application, error) {
	var resp struct {
		Data Application `json:"data"`
	}
	return resp.Data, s.client.create(ctx, "applications", data, &resp)
}

// List lists all applications.
func (s *ApplicationsService) List(ctx context.Context) (ApplicationsResponse, error) {
	resp := ApplicationsResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "applications", &resp)
}
