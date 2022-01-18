package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

func resourceOSImageVersion() *schema.Resource {
	return &schema.Resource{
		CreateContext: adoptCreate("OS Image Version", resourceOSImageVersionCreate),
		ReadContext:   adoptRead("OS Image Version", resourceOSImageVersionRead),
		UpdateContext: adoptUpdate("OS Image Version", resourceOSImageVersionUpdate),
		DeleteContext: adoptDelete("OS Image Version", resourceOSImageVersionDelete),

		Schema: map[string]*schema.Schema{
			"os_image_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"version": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"url": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"cloud_init_version": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(i interface{}, k string) ([]string, []error) {
					v, ok := i.(string)
					if !ok {
						return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
					}
					if !solus.IsValidCloudInitVersion(v) {
						return nil, []error{fmt.Errorf("unknown cloud init version %q", v)}
					}
					return nil, nil
				},
			},
			"virtualization_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validationIsVirtualizationType,
			},
		},
	}
}

func resourceOSImageVersionCreate(ctx context.Context, client *solus.Client, d *schema.ResourceData) error {
	res, err := client.OsImages.CreateVersion(ctx, d.Get("os_image_id").(int), buildOSImageVersionRequest(d))
	if err != nil {
		return normalizeAPIError(err)
	}

	d.SetId(strconv.Itoa(res.ID))
	return resourceOSImageVersionRead(ctx, client, d)
}

func resourceOSImageVersionRead(ctx context.Context, client *solus.Client, d *schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	res, err := client.OsImageVersions.Get(ctx, id)
	if err != nil {
		return normalizeAPIError(err)
	}

	return newSchemaChainSetter(d).
		SetID(res.ID).
		Set("os_image_id", res.OsImageID).
		Set("version", res.Version).
		Set("url", res.URL).
		Set("cloud_init_version", res.CloudInitVersion).
		Set("virtualization_type", res.VirtualizationType).
		Error()
}

func resourceOSImageVersionUpdate(ctx context.Context, client *solus.Client, d *schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	res, err := client.OsImageVersions.Update(ctx, id, buildOSImageVersionRequest(d))
	if err != nil {
		return normalizeAPIError(err)
	}

	d.SetId(strconv.Itoa(res.ID))
	return resourceOSImageVersionRead(ctx, client, d)
}

func resourceOSImageVersionDelete(ctx context.Context, client *solus.Client, d *schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	return normalizeAPIError(client.OsImageVersions.Delete(ctx, id))
}

func buildOSImageVersionRequest(d *schema.ResourceData) solus.OsImageVersionRequest {
	return solus.OsImageVersionRequest{
		Version:            d.Get("version").(string),
		URL:                d.Get("url").(string),
		CloudInitVersion:   solus.CloudInitVersion(d.Get("cloud_init_version").(string)),
		VirtualizationType: solus.VirtualizationType(d.Get("virtualization_type").(string)),
	}
}
