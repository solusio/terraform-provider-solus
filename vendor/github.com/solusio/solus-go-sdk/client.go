//go:generate go run generators/paginatorgen.go

package solus

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
)

// Client a Solus API client.
type Client struct {
	BaseURL     *url.URL
	UserAgent   string
	Credentials Credentials
	Headers     http.Header
	HTTPClient  *http.Client
	Logger      Logger
	Retries     int
	RetryAfter  time.Duration

	s service

	Account           *AccountService
	ActivityLogs      *ActivityLogsService
	Applications      *ApplicationsService
	BackupNodes       *BackupNodesService
	Backups           *BackupsService
	ComputeResources  *ComputeResourcesService
	IPBlocks          *IPBlocksService
	Icons             *IconsService
	License           *LicenseService
	Locations         *LocationsService
	OsImageVersions   *OsImageVersionsService
	OsImages          *OsImagesService
	Permission        *PermissionsService
	Plans             *PlansService
	Projects          *ProjectsService
	Roles             *RolesService
	SSHKeys           *SSHKeysService
	ServersMigrations *ServersMigrationsService
	Settings          *SettingsService
	Snapshots         *SnapshotsService
	Storage           *StorageService
	StorageTypes      *StorageTypesService
	Tasks             *TasksService
	Users             *UsersService
	VirtualServers    *VirtualServersService
}

type service struct {
	client *Client
}

// Authenticator interface for client authentication.
type Authenticator interface {
	// Authenticate authenticates client and return credentials
	// which should be used for making further API calls.
	// The Client is fully initialized. Any endpoints which is not requires
	// authentication may be called.
	Authenticate(c *Client) (Credentials, error)
}

// EmailAndPasswordAuthenticator authenticate with specified email
// and password.
type EmailAndPasswordAuthenticator struct {
	Email    string
	Password string
}

var _ Authenticator = EmailAndPasswordAuthenticator{}

// Authenticate authenticates by email and password.
func (a EmailAndPasswordAuthenticator) Authenticate(c *Client) (Credentials, error) {
	resp, err := c.authLogin(context.Background(), AuthLoginRequest(a))
	if err != nil {
		return Credentials{}, err
	}

	return resp.Credentials, nil
}

// APITokenAuthenticator authenticate by provided API token.
type APITokenAuthenticator struct {
	Token string
}

var _ Authenticator = APITokenAuthenticator{}

// Authenticate authenticates by API token.
func (a APITokenAuthenticator) Authenticate(*Client) (Credentials, error) {
	return Credentials{
		AccessToken: a.Token,
		TokenType:   "Bearer",
		ExpiresAt:   "",
	}, nil
}

// ClientOption represent client initialization options.
type ClientOption func(c *Client)

// AllowInsecure allows skipping certificate verify.
func AllowInsecure() ClientOption {
	return func(c *Client) {
		c.HTTPClient.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec // We should give an ability to disable cert check.
	}
}

// SetRetryPolicy sets number of retries and timeout between them.
func SetRetryPolicy(retries int, retryAfter time.Duration) ClientOption {
	return func(c *Client) {
		c.Retries = retries
		c.RetryAfter = retryAfter
	}
}

// WithLogger inject specific logger into client.
func WithLogger(logger Logger) ClientOption {
	return func(c *Client) {
		c.Logger = logger
	}
}

// NewClient create and initialize Client instance.
func NewClient(
	baseURL *url.URL,
	a Authenticator,
	opts ...ClientOption,
) (*Client, error) {
	client := &Client{
		BaseURL:   baseURL,
		UserAgent: "Go SDK client",
		Headers: map[string][]string{
			"Accept":       {"application/json"},
			"Content-Type": {"application/json"},
		},
		HTTPClient: &http.Client{
			Timeout:   time.Second * 35,
			Transport: http.DefaultTransport.(*http.Transport).Clone(),
		},
		Logger:     NullLogger{},
		Retries:    5,
		RetryAfter: 1 * time.Second,
	}

	for _, o := range opts {
		o(client)
	}

	c, err := a.Authenticate(client)
	if err != nil {
		return nil, err
	}

	client.Credentials = c
	client.Headers["Authorization"] = []string{client.Credentials.TokenType + " " + client.Credentials.AccessToken}

	client.s.client = client

	client.Account = (*AccountService)(&client.s)
	client.ActivityLogs = (*ActivityLogsService)(&client.s)
	client.Applications = (*ApplicationsService)(&client.s)
	client.BackupNodes = (*BackupNodesService)(&client.s)
	client.Backups = (*BackupsService)(&client.s)
	client.ComputeResources = (*ComputeResourcesService)(&client.s)
	client.IPBlocks = (*IPBlocksService)(&client.s)
	client.Icons = (*IconsService)(&client.s)
	client.License = (*LicenseService)(&client.s)
	client.Locations = (*LocationsService)(&client.s)
	client.OsImageVersions = (*OsImageVersionsService)(&client.s)
	client.OsImages = (*OsImagesService)(&client.s)
	client.Permission = (*PermissionsService)(&client.s)
	client.Plans = (*PlansService)(&client.s)
	client.Projects = (*ProjectsService)(&client.s)
	client.Roles = (*RolesService)(&client.s)
	client.SSHKeys = (*SSHKeysService)(&client.s)
	client.ServersMigrations = (*ServersMigrationsService)(&client.s)
	client.Settings = (*SettingsService)(&client.s)
	client.Snapshots = (*SnapshotsService)(&client.s)
	client.Storage = (*StorageService)(&client.s)
	client.StorageTypes = (*StorageTypesService)(&client.s)
	client.Tasks = (*TasksService)(&client.s)
	client.Users = (*UsersService)(&client.s)
	client.VirtualServers = (*VirtualServersService)(&client.s)

	return client, nil
}

func (c *Client) authLogin(ctx context.Context, data AuthLoginRequest) (AuthLoginResponse, error) {
	const path = "auth/login"
	body, code, err := c.request(ctx, http.MethodPost, path, withBody(data))
	if err != nil {
		return AuthLoginResponse{}, err
	}

	if code != http.StatusOK {
		return AuthLoginResponse{}, newHTTPError(http.MethodPost, path, code, body)
	}

	var resp struct {
		Data AuthLoginResponse `json:"data"`
	}
	return resp.Data, unmarshal(body, &resp)
}
