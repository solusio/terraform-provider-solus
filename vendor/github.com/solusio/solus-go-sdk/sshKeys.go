package solus

import (
	"context"
	"fmt"
)

// SSHKeysService handles all available methods with SSH keys.
type SSHKeysService service

// SSHKey represents a SSH key.
type SSHKey struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Body string `json:"body"`
}

// SSHKeyCreateRequest represents available properties for creating a new SSH key.
type SSHKeyCreateRequest struct {
	Name   string `json:"name"`
	Body   string `json:"body"`
	UserID int    `json:"user_id"`
}

// SSHKeysResponse represents paginated list of SSH keys.
// This cursor can be used for iterating over all available SSH keys.
type SSHKeysResponse struct {
	paginatedResponse

	Data []SSHKey `json:"data"`
}

type sshKeyResponse struct {
	Data SSHKey `json:"data"`
}

// List lists SSH keys.
func (s *SSHKeysService) List(ctx context.Context, filter *FilterSSHKeys) (SSHKeysResponse, error) {
	resp := SSHKeysResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "ssh_keys", &resp, withFilter(filter.data))
}

// Get gets specified SSH key.
func (s *SSHKeysService) Get(ctx context.Context, id int) (SSHKey, error) {
	var resp sshKeyResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("ssh_keys/%d", id), &resp)
}

// Create creates new SSH key.
func (s *SSHKeysService) Create(ctx context.Context, data SSHKeyCreateRequest) (SSHKey, error) {
	var resp sshKeyResponse
	return resp.Data, s.client.create(ctx, "ssh_keys", data, &resp)
}

// Delete deletes specified SSH key.
func (s *SSHKeysService) Delete(ctx context.Context, id int) error {
	return s.client.syncDelete(ctx, fmt.Sprintf("ssh_keys/%d", id))
}
