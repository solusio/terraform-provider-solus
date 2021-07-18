package solus

import (
	"context"
	"fmt"
)

type SSHKeysService service

type SSHKey struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Body string `json:"body"`
}

type SSHKeyCreateRequest struct {
	Name   string `json:"name"`
	Body   string `json:"body"`
	UserID int    `json:"user_id"`
}

func (s *SSHKeysService) Create(ctx context.Context, data SSHKeyCreateRequest) (SSHKey, error) {
	var resp struct {
		Data SSHKey `json:"data"`
	}
	return resp.Data, s.client.create(ctx, "ssh_keys", data, &resp)
}

func (s *SSHKeysService) Delete(ctx context.Context, id int) error {
	return s.client.delete(ctx, fmt.Sprintf("ssh_keys/%d", id))
}
