package solus

import (
	"context"
)

// AccountService handles all available methods with a current user account.
type AccountService service

// Get retrieves current user account.
func (s *AccountService) Get(ctx context.Context) (User, error) {
	var resp struct {
		Data User `json:"data"`
	}
	return resp.Data, s.client.get(ctx, "account", &resp)
}
