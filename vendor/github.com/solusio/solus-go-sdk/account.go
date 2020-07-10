package solus

import (
	"context"
)

type AccountService service

func (s *AccountService) Get(ctx context.Context) (User, error) {
	var resp struct {
		Data User `json:"data"`
	}
	return resp.Data, s.client.get(ctx, "account", &resp)
}
