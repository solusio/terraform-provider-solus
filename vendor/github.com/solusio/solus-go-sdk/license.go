package solus

import (
	"context"
	"net/http"
)

// LicenseService handles all available methods with license.
type LicenseService service

// License represent a license.
type License struct {
	CPUCores       int    `json:"cpu_cores"`
	CPUCoresInUse  int    `json:"cpu_cores_in_use"`
	IsActive       bool   `json:"is_active"`
	Key            string `json:"key"`
	KeyType        string `json:"key_type"`
	Product        string `json:"product"`
	ExpirationDate string `json:"expiration_date"`
	UpdateDate     string `json:"update_date"`
}

// LicenseActivateRequest represents available properties for activating the license.
type LicenseActivateRequest struct {
	ActivationCode string `json:"activation_code"`
}

// Activate activates the license.
func (s *LicenseService) Activate(ctx context.Context, data LicenseActivateRequest) (License, error) {
	const path = "license/activate"
	body, code, err := s.client.request(ctx, http.MethodPost, path, withBody(data))
	if err != nil {
		return License{}, err
	}

	if code != http.StatusOK {
		return License{}, newHTTPError(http.MethodPost, path, code, body)
	}

	var resp struct {
		Data License `json:"data"`
	}
	return resp.Data, unmarshal(body, &resp)
}
