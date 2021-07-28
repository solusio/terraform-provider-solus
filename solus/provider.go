package solus

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
	baseURLEnv  = "SOLUSIO_BASE_URL"
	tokenEnv    = "SOLUSIO_TOKEN" //nolint:gosec // False positive.
	insecureEnv = "SOLUSIO_INSECURE"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"base_url": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "SOLUS IO API base url like 'https://solus.io:4444'",
				ValidateFunc: validation.NoZeroValues,
				DefaultFunc:  schema.EnvDefaultFunc(baseURLEnv, nil),
			},
			"token": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				Description:  "SOLUS IO auth token",
				ValidateFunc: validation.NoZeroValues,
				DefaultFunc:  schema.EnvDefaultFunc(tokenEnv, nil),
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Skip certificate validation",
				DefaultFunc: func() (interface{}, error) {
					return strings.TrimSpace(os.Getenv(insecureEnv)) == "1", nil
				},
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"solusio_icon":             dataSourceIcon(),
			"solusio_ip_block":         dataSourceIPBlock(),
			"solusio_location":         dataSourceLocation(),
			"solusio_os_image":         dataSourceOsImage(),
			"solusio_os_image_version": dataSourceOsImageVersion(),
			"solusio_plan":             dataSourcePlan(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"solusio_ip_block":         resourceIPBlock(),
			"solusio_location":         resourceLocation(),
			"solusio_os_image":         resourceOsImage(),
			"solusio_os_image_version": resourceOsImageVersion(),
			"solusio_plan":             resourcePlan(),
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
		return nil, diag.Errorf("failed to initialize Solus IO client: %s", err)
	}

	return client, nil
}
