package solus

// LimitGroup represent limit group.
type LimitGroup struct {
	ID                    int             `json:"id"`
	Name                  string          `json:"name"`
	CreatedAt             string          `json:"created_at"`
	UpdatedAt             string          `json:"updated_at"`
	VirtualServers        LimitGroupLimit `json:"vms"`
	RunningVirtualServers LimitGroupLimit `json:"running_vms"`
	AdditionalIPs         LimitGroupLimit `json:"additional_ips"`
}

// LimitGroupLimit represent limit's group limit.
type LimitGroupLimit struct {
	Limit     int  `json:"limit"`
	IsEnabled bool `json:"is_enabled"`
}
