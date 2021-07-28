package solus

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

func resourceOsImageVersion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOsImageVersionCreate,
		ReadContext:   resourceOsImageVersionRead,
		UpdateContext: resourceOsImageVersionUpdate,
		DeleteContext: resourceOsImageVersionDelete,

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
		},
	}
}

func resourceOsImageVersionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*solus.Client)

	osImageID := d.Get("os_image_id").(int)
	version := d.Get("version").(string)
	url := d.Get("url").(string)
	cloudInitVersion := d.Get("cloud_init_version").(string)

	v, err := client.OsImages.CreateVersion(ctx, osImageID, solus.OsImageVersionRequest{
		Version:          version,
		URL:              url,
		CloudInitVersion: solus.CloudInitVersion(cloudInitVersion),
	})
	if err != nil {
		return diag.Errorf("failed to create new os image version: %s", err)
	}

	d.SetId(strconv.Itoa(v.ID))
	return resourceOsImageVersionRead(ctx, d, m)
}

func resourceOsImageVersionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*solus.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	v, err := client.OsImageVersions.Get(ctx, id)
	if err != nil {
		return diag.Errorf("failed to get os image version by id %d: %s", id, err)
	}

	err = (&schemaChainSetter{d: d}).
		SetID(v.ID).
		Set("os_image_id", v.OsImageID).
		Set("version", v.Version).
		Set("url", v.URL).
		Set("cloud_init_version", v.CloudInitVersion).
		Error()
	if err != nil {
		return diag.Errorf("failed to map os image version response to resource: %s", err)
	}

	return nil
}

func resourceOsImageVersionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*solus.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	version := d.Get("version").(string)
	url := d.Get("url").(string)
	cloudInitVersion := d.Get("cloud_init_version").(string)

	v, err := client.OsImageVersions.Update(ctx, id, solus.OsImageVersionRequest{
		Version:          version,
		URL:              url,
		CloudInitVersion: solus.CloudInitVersion(cloudInitVersion),
	})
	if err != nil {
		return diag.Errorf("failed to update os image version with id %d: %s", id, err)
	}

	d.SetId(strconv.Itoa(v.ID))
	return resourceOsImageVersionRead(ctx, d, m)
}

func resourceOsImageVersionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*solus.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.OsImageVersions.Delete(ctx, id)
	if err != nil {
		return diag.Errorf("failed to delete os image version by id %d: %s", id, err)
	}
	return nil
}
