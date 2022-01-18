package provider

import (
	"context"
	"net/url"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

const (
	baseURLEnv  = "SOLUS_BASE_URL"
	tokenEnv    = "SOLUS_TOKEN"
	insecureEnv = "SOLUS_INSECURE"
)

func New() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"base_url": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Solus API base url like 'https://solus.io:4444'",
				ValidateFunc: validation.NoZeroValues,
				DefaultFunc:  schema.EnvDefaultFunc(baseURLEnv, ""),
			},
			"token": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				Description:  "Solus auth token",
				ValidateFunc: validation.NoZeroValues,
				DefaultFunc:  schema.EnvDefaultFunc(tokenEnv, ""),
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Skip SSL/TLS certificate validation",
				DefaultFunc: func() (interface{}, error) {
					return strings.TrimSpace(os.Getenv(insecureEnv)) == "1", nil
				},
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"solus_icon":             dataSourceIcon(),
			"solus_ip_block":         dataSourceIPBlock(),
			"solus_location":         dataSourceLocation(),
			"solus_os_image":         dataSourceOsImage(),
			"solus_os_image_version": dataSourceOsImageVersion(),
			"solus_plan":             dataSourcePlan(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"solus_ip_block":         resourceIPBlock(),
			"solus_location":         resourceLocation(),
			"solus_os_image":         resourceOSImage(),
			"solus_os_image_version": resourceOSImageVersion(),
			"solus_plan":             resourcePlan(),
		},

		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	rawBaseURL := d.Get("base_url").(string)
	token := d.Get("token").(string)
	insecure := d.Get("insecure").(bool)

	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, diag.Errorf("failed to parse base URL %q: %s", rawBaseURL, err)
	}

	opts := make([]solus.ClientOption, 0)

	if insecure {
		opts = append(opts, solus.AllowInsecure())
	}

	client, err := solus.NewClient(baseURL, solus.APITokenAuthenticator{Token: token}, opts...)
	if err != nil {
		return nil, diag.Errorf("failed to initialize API client: %s", err)
	}

	return client, nil
}
