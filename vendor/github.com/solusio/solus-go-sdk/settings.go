package solus

import (
	"context"
)

// SettingsService handles all available methods with settings.
type SettingsService service

// Settings represents application settings.
type Settings struct {
	Hostname              string                        `json:"hostname"`
	Mail                  SettingsMail                  `json:"mail"`
	SendStatistics        bool                          `json:"send_statistics"`
	LimitGroup            LimitGroup                    `json:"limit_group"`
	Registration          SettingsRegistration          `json:"registration"`
	Theme                 SettingsTheme                 `json:"theme"`
	NetworkRules          SettingsNetworkRules          `json:"network_rules"`
	Notifications         SettingsNotifications         `json:"notifications"`
	ComputeResource       SettingsComputeResource       `json:"compute_resource"`
	NonExistentVMSRemover SettingsNonExistentVMSRemover `json:"non_existent_vms_remover"`
	BillingIntegration    SettingsBillingIntegration    `json:"billing_integration"`
	DNS                   SettingsDNS                   `json:"dns"`
	Update                SettingsUpdate                `json:"update"`
	Features              SettingsFeatures              `json:"features"`
}

// SettingsUpdateRequest represents available properties for updating settings.
type SettingsUpdateRequest struct {
	Hostname              *string                        `json:"hostname,omitempty"`
	Mail                  *SettingsMail                  `json:"mail,omitempty"`
	SendStatistics        *bool                          `json:"send_statistics,omitempty"`
	LimitGroup            *LimitGroup                    `json:"limit_group,omitempty"`
	Registration          *SettingsRegistration          `json:"registration,omitempty"`
	Theme                 *SettingsTheme                 `json:"theme,omitempty"`
	NetworkRules          *SettingsNetworkRules          `json:"network_rules,omitempty"`
	Notifications         *SettingsNotifications         `json:"notifications,omitempty"`
	ComputeResource       *SettingsComputeResource       `json:"compute_resource,omitempty"`
	NonExistentVMSRemover *SettingsNonExistentVMSRemover `json:"non_existent_vms_remover,omitempty"`
	BillingIntegration    *SettingsBillingIntegration    `json:"billing_integration,omitempty"`
	DNS                   *SettingsDNS                   `json:"dns,omitempty"`
	Update                *SettingsUpdate                `json:"update,omitempty"`
	Features              *SettingsFeatures              `json:"features,omitempty"`
}

// SettingsMail represents mail settings.
type SettingsMail struct {
	Host string `json:"host,omitempty"`
	// Port       string `json:"port,omitempty"` todo uncomment after fix SIO-3682.
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	FromEmail  string `json:"from_email,omitempty"`
	FromName   string `json:"from_name,omitempty"`
	Encryption bool   `json:"encryption,omitempty"`
}

// SettingsRegistration represents registration settings.
type SettingsRegistration struct {
	Role string `json:"role,omitempty"`
}

// SettingsTheme represents application theme setting.
type SettingsTheme struct {
	PrimaryColor          string `json:"primary_color,omitempty"`
	SecondaryColor        string `json:"secondary_color,omitempty"`
	BrandName             string `json:"brand_name,omitempty"`
	Logo                  string `json:"logo,omitempty"`
	Favicon               string `json:"favicon,omitempty"`
	TermsAndConditionsURL string `json:"terms_and_conditions_url,omitempty"`
}

// SettingsNetworkRules represents network rules settings.
type SettingsNetworkRules struct {
	ARP       bool `json:"arp,omitempty"`
	DHCP      bool `json:"dhcp,omitempty"`
	CloudInit bool `json:"cloud_init,omitempty"`
	SMTP      bool `json:"smtp,omitempty"`
	ICMP      bool `json:"icmp,omitempty"`
	ICMPReply bool `json:"icmp_reply"`
}

// SettingsNotifications represents notification settings.
type SettingsNotifications struct {
	ServerCreate                  SettingsNotificationsTemplate `json:"server_create,omitempty"`
	ServerResetPassword           SettingsNotificationsTemplate `json:"server_reset_password,omitempty"`
	UserResetPassword             SettingsNotificationsTemplate `json:"user_reset_password,omitempty"`
	UserVerifyEmail               SettingsNotificationsTemplate `json:"user_verify_email,omitempty"`
	ProjectUserInvite             SettingsNotificationsTemplate `json:"project_user_invite,omitempty"`
	ProjectUserLeft               SettingsNotificationsTemplate `json:"project_user_left,omitempty"`
	ServerIncomingTrafficExceeded SettingsNotificationsTemplate `json:"server_incoming_traffic_exceeded,omitempty"`
	ServerOutgoingTrafficExceeded SettingsNotificationsTemplate `json:"server_outgoing_traffic_exceeded,omitempty"`
}

// SettingsNotificationsTemplate represents notification template.
type SettingsNotificationsTemplate struct {
	Enabled          bool              `json:"enabled,omitempty"`
	SubjectTemplates map[string]string `json:"subject_templates,omitempty"`
	BodyTemplated    map[string]string `json:"body_templated,omitempty"`
}

// SettingsComputeResource represents compute resource settings.
type SettingsComputeResource struct {
	RescueISOURL    string `json:"rescue_iso_url,omitempty"`
	BalanceStrategy string `json:"balance_strategy,omitempty"`
}

// SettingsNonExistentVMSRemover represents settings for virtual server remover.
type SettingsNonExistentVMSRemover struct {
	Enabled  bool `json:"enabled,omitempty"`
	Interval int  `json:"interval,omitempty"`
}

// SettingsBillingIntegration represents settings for billing.
type SettingsBillingIntegration struct {
	Type    string                            `json:"type,omitempty"`
	Drivers SettingsBillingIntegrationDrivers `json:"drivers,omitempty"`
}

// SettingsBillingIntegrationDrivers represents billing drivers.
type SettingsBillingIntegrationDrivers struct {
	WHMCS SettingsBillingIntegrationDriversWHMCS `json:"whmcs,omitempty"`
}

// SettingsBillingIntegrationDriversWHMCS represents WHMCS driver.
type SettingsBillingIntegrationDriversWHMCS struct {
	URL   string `json:"url,omitempty"`
	Token string `json:"token,omitempty"`
}

// SettingsDNS represents DNS settings.
type SettingsDNS struct {
	Type                       string             `json:"type,omitempty"`
	ServerHostnameTemplate     string             `json:"server_hostname_template,omitempty"`
	RegisterFQDNOnServerCreate bool               `json:"register_fqdn_on_server_create,omitempty"`
	ReverseDNSDomainTemplate   string             `json:"reverse_dns_domain_template,omitempty"`
	Drivers                    SettingsDNSDrivers `json:"drivers,omitempty"`
}

// SettingsDNSDrivers represents DNS drivers.
type SettingsDNSDrivers struct {
	PowerDNS SettingsDNSDriversPowerDNS `json:"power_dns,omitempty"`
}

// SettingsDNSDriversPowerDNS represents PowerDNS driver.
type SettingsDNSDriversPowerDNS struct {
	Host   string `json:"host,omitempty"`
	APIKey string `json:"api_key,omitempty"`
	// Port   string `json:"port,omitempty"` todo Uncomment after fix SIO-3682
}

// SettingsUpdate represents application update settings.
type SettingsUpdate struct {
	Method        string `json:"method,omitempty"`
	Channel       string `json:"channel,omitempty"`
	ScheduledDays []int  `json:"scheduled_days,omitempty"`
	ScheduledTime string `json:"scheduled_time,omitempty"`
}

// SettingsFeatures represents application features settings.
type SettingsFeatures struct {
	HidePlanName          bool `json:"hide_plan_name"`
	HideUserData          bool `json:"hide_user_data"`
	HidePlanSection       bool `json:"hide_plan_section"`
	HideLocationSection   bool `json:"hide_location_section"`
	AllowRegistration     bool `json:"allow_registration"`
	AllowPasswordRecovery bool `json:"allow_password_recovery"`
}

type settingsResponse struct {
	Data Settings `json:"data"`
}

// Get gets settings.
func (s *SettingsService) Get(ctx context.Context) (Settings, error) {
	var resp settingsResponse
	return resp.Data, s.client.get(ctx, "settings", &resp)
}

// Patch patches settings.
func (s *SettingsService) Patch(ctx context.Context, data SettingsUpdateRequest) (Settings, error) {
	var resp settingsResponse
	return resp.Data, s.client.patch(ctx, "settings", data, &resp)
}
