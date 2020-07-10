package solus

import (
	"context"
	"fmt"
)

type PlansService service

type PlanParams struct {
	Disk int `json:"disk"`
	RAM  int `json:"ram"`
	VCPU int `json:"vcpu"`
}

type PlanLimit struct {
	IsEnabled bool `json:"is_enabled"`
	Limit     int  `json:"limit"`
}

type PlanLimits struct {
	TotalBytes PlanLimit `json:"total_bytes"`
	TotalIops  PlanLimit `json:"total_iops"`
}

type PlanPrice struct {
	PerHour        string        `json:"per_hour"`
	PerMonth       string        `json:"per_month"`
	CurrencyCode   string        `json:"currency_code"`
	TaxesInclusive bool          `json:"taxes_inclusive"`
	Taxes          []interface{} `json:"taxes"`
	TotalPrice     string        `json:"total_price"`
	BackupPrice    string        `json:"backup_price"`
}

type Plan struct {
	ID                  int        `json:"id"`
	Name                string     `json:"name"`
	Params              PlanParams `json:"params"`
	StorageType         string     `json:"storage_type"`
	ImageFormat         string     `json:"image_format"`
	IsDefault           bool       `json:"is_default"`
	IsSnapshotAvailable bool       `json:"is_snapshot_available"`
	IsSnapshotsEnabled  bool       `json:"is_snapshots_enabled"`
	IsBackupAvailable   bool       `json:"is_backup_available"`
	BackupPrice         float32    `json:"backup_price"`
	IsVisible           bool       `json:"is_visible"`
	Limits              PlanLimits `json:"limits"`
	TokensPerHour       float64    `json:"tokens_per_hour"`
	TokensPerMonth      float64    `json:"tokens_per_month"`
	Position            float64    `json:"position"`
	Price               PlanPrice  `json:"price"`
}

type PlanCreateRequest struct {
	Name               string          `json:"name"`
	Params             PlanParams      `json:"params"`
	StorageType        StorageTypeName `json:"storage_type"`
	ImageFormat        ImageFormat     `json:"image_format"`
	Limits             PlanLimits      `json:"limits"`
	TokensPerHour      float64         `json:"tokens_per_hour"`
	TokensPerMonth     float64         `json:"tokens_per_month"`
	Position           float64         `json:"position"`
	IsVisible          bool            `json:"is_visible"`
	IsDefault          bool            `json:"is_default"`
	IsSnapshotsEnabled bool            `json:"is_snapshots_enabled"`
	IsBackupAvailable  bool            `json:"is_backup_available"`
	BackupPrice        float32         `json:"backup_price"`
}

type PlansResponse struct {
	paginatedResponse

	Data []Plan `json:"data"`
}

func (s *PlansService) List(ctx context.Context) (PlansResponse, error) {
	resp := PlansResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "plans", &resp)
}

func (s *PlansService) Create(ctx context.Context, data PlanCreateRequest) (Plan, error) {
	var resp struct {
		Data Plan `json:"data"`
	}
	return resp.Data, s.client.create(ctx, "plans", data, &resp)
}

func (s *PlansService) Delete(ctx context.Context, id int) error {
	return s.client.delete(ctx, fmt.Sprintf("plans/%d", id))
}
