package solus

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/solusio/solus-go-sdk"
	"net/url"
	"time"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"base_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SOLUS IO API base url like 'https://solus.io:4444'",
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SOLUS IO auth token",
				Sensitive:   true,
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Skip certificate validation",
				Default:     false,
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"solusio_location": dataSourceLocation(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"solusio_location": resourceLocation(),
		},

		ConfigureContextFunc: configureProvider,
	}
}

type metadata struct {
	Client         *solus.Client
	RequestTimeout time.Duration
}

func configureProvider(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	rawURL := d.Get("base_url").(string)
	token := d.Get("token").(string)
	insecure := d.Get("insecure").(bool)

	baseURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, diagErr("Failed to parse base URL", "Invalid URL %q: %s", rawURL, err)
	}

	opts := make([]solus.ClientOption, 0)

	if insecure {
		opts = append(opts, solus.AllowInsecure())
	}

	client, err := solus.NewClient(baseURL, solus.APITokenAuthenticator{Token: token}, opts...)
	if err != nil {
		return nil, diagErr("Failed to initialize Solus IO client", err.Error())
	}

	return metadata{
		Client: client,
	}, nil
}
