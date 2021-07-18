package solus

import (
	"context"
	"fmt"
	"net/url"
	"time"

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
			"requests_timeout": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "1m",
				Description: "Timeout for each API request",
				ValidateFunc: func(i interface{}, k string) ([]string, []error) {
					v, ok := i.(string)
					if !ok {
						return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
					}

					if _, err := time.ParseDuration(v); err != nil {
						return nil, []error{fmt.Errorf("can't parse duration from %q: %w", v, err)}
					}
					return nil, nil
				},
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"solusio_location": dataSourceLocation(),
			"solusio_icon":     dataSourceIcon(),
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
	rawRequestTimeout := d.Get("requests_timeout").(string)

	baseURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, diag.Errorf("failed to parse base URL %q: %s", rawURL, err)
	}

	requestTimeout, err := time.ParseDuration(rawRequestTimeout)
	if err != nil {
		return nil, diag.Errorf("failed to parse request timeout %q: %s", rawRequestTimeout, err)
	}

	opts := make([]solus.ClientOption, 0)

	if insecure {
		opts = append(opts, solus.AllowInsecure())
	}

	client, err := solus.NewClient(baseURL, solus.APITokenAuthenticator{Token: token}, opts...)
	if err != nil {
		return nil, diag.Errorf("failed to initialize Solus IO client: %s", err)
	}

	return metadata{
		Client:         client,
		RequestTimeout: requestTimeout,
	}, nil
}
