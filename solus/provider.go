package solus

import (
	"context"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"base_url": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "SOLUS IO API base url like 'https://solus.io:4444'",
				ValidateFunc: validation.NoZeroValues,
			},
			"token": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				Description:  "SOLUS IO auth token",
				ValidateFunc: validation.NoZeroValues,
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Skip certificate validation",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"solusio_icon":             dataSourceIcon(),
			"solusio_location":         dataSourceLocation(),
			"solusio_os_image":         dataSourceOsImage(),
			"solusio_os_image_version": dataSourceOsImageVersion(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"solusio_location":         resourceLocation(),
			"solusio_os_image":         resourceOsImage(),
			"solusio_os_image_version": resourceOsImageVersion(),
		},

		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	rawURL := d.Get("base_url").(string)
	token := d.Get("token").(string)
	insecure := d.Get("insecure").(bool)

	baseURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, diag.Errorf("failed to parse base URL %q: %s", rawURL, err)
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
