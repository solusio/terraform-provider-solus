package solus

import (
	"context"
	"fmt"
)

// UsersService handles all available methods with users.
type UsersService service

// UserStatus represents available user's statuses.
type UserStatus string

const (
	// UserStatusActive indicates user is active.
	UserStatusActive UserStatus = "active"

	// UserStatusLocked indicates user is locked.
	UserStatusLocked UserStatus = "locked"

	// UserStatusSuspended indicates user is suspended.
	UserStatusSuspended UserStatus = "suspended"
)

// User represents a user.
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	// CreatedAt for date in RFC3339Nano format
	CreatedAt string     `json:"created_at"`
	Status    UserStatus `json:"status"`
	Roles     []Role     `json:"roles"`
}

// UsersResponse represents paginated list of users.
// This cursor can be used for iterating over all available users.
type UsersResponse struct {
	paginatedResponse

	Data []User `json:"data"`
}

// UserCreateRequest represents available properties for creating a new user.
type UserCreateRequest struct {
	Password   string `json:"password,omitempty"`
	Email      string `json:"email,omitempty"`
	Status     string `json:"status,omitempty"`
	LanguageID int    `json:"language_id,omitempty"`
	Roles      []int  `json:"roles,omitempty"`
}

// UserUpdateRequest represents available properties for updating exists user.
type UserUpdateRequest struct {
	Password   string `json:"password,omitempty"`
	Status     string `json:"status,omitempty"`
	LanguageID int    `json:"language_id,omitempty"`
	Roles      []int  `json:"roles,omitempty"`
}

type userResponse struct {
	Data User `json:"data"`
}

// List lists users.
func (s *UsersService) List(ctx context.Context, filter *FilterUsers) (UsersResponse, error) {
	resp := UsersResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "users", &resp, withFilter(filter.data))
}

// Create creates new user.
func (s *UsersService) Create(ctx context.Context, data UserCreateRequest) (User, error) {
	var resp userResponse
	return resp.Data, s.client.create(ctx, "users", data, &resp)
}

// Update updates specified user.
func (s *UsersService) Update(ctx context.Context, id int, data UserUpdateRequest) (User, error) {
	var resp userResponse
	return resp.Data, s.client.update(ctx, fmt.Sprintf("users/%d", id), data, &resp)
}

// Delete deletes specified user.
func (s *UsersService) Delete(ctx context.Context, id int) error {
	return s.client.syncDelete(ctx, fmt.Sprintf("users/%d", id))
}
